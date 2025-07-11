package handler

import (
	"common/kafka/producer"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/message"
	"github.com/vietquan-37/auth-service/pkg/oauth2"
	"sync"

	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	wg       *sync.WaitGroup
	Jwt      config.JwtWrapper
	Repo     repository.IAuthRepo
	Config   config.Config
	Producer *producer.Producer
}

func NewAuthHandler(jwt config.JwtWrapper, repo repository.IAuthRepo, config config.Config, wait *sync.WaitGroup, producer *producer.Producer) *Handler {
	return &Handler{
		Jwt:      jwt,
		Repo:     repo,
		Config:   config,
		wg:       wait,
		Producer: producer,
	}
}
func (handler *Handler) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

	hashPassword, err := config.HashedPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while hashing password: %s", err)
	}
	req.Password = hashPassword
	model := convertUser(req)

	user, err := handler.Repo.CreateUser(ctx, model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Errorf(codes.AlreadyExists, "email %s already register before", req.UserName)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Errorf(codes.DeadlineExceeded, "Request timeout: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error while creating user: %s", err)
	}

	token, err := handler.Jwt.GenerateJWT(user, time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}
	payload, err := message.NewUserCreatedEnvelope("auth-service", "1", message.UserCreateMessage{
		ID:       int32(user.ID),
		Email:    user.Username,
		FullName: user.FullName,
		Token:    token})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating user message: %s", err)
	}
	handler.wg.Add(1)
	go func() {
		defer handler.wg.Done()
		if err := handler.Producer.SendMessage(context.Background(), handler.Config.Topic, nil, payload); err != nil {
			log.Error().
				Err(err).
				Uint("user_id", user.ID).
				Str("email", user.Username).
				Str("topic", handler.Config.Topic).
				Msg("CRITICAL: Failed to send user created message to Kafka")
		}
	}()

	return convertUserResponse(user), nil
}
func (handler *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := handler.Repo.GetUserByUserName(ctx, req.GetUserName())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with email %s not found", req.GetUserName())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Errorf(codes.DeadlineExceeded, "Request timeout: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
	}
	if err = config.CheckPassword(req.GetPassword(), user.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, err := handler.Jwt.GenerateJWT(user, time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}
	refreshToken, err := handler.Jwt.GenerateJWT(user, time.Hour*5)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}

	rsp := &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       int32(user.ID),
	}
	return rsp, nil
}
func (handler *Handler) GetOneUser(ctx context.Context, req *pb.GetOneUserRequest) (*pb.UserResponse, error) {
	user, err := handler.Repo.FindOneUser(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Errorf(codes.DeadlineExceeded, "Request timeout: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
	}
	rsp := convertUserResponse(user)

	return rsp, nil
}
func (handler *Handler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := handler.Jwt.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	user, err := handler.Repo.FindOneUser(ctx, claims.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", claims.Id)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Errorf(codes.DeadlineExceeded, "Request timeout: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
	}
	return convertValidate(user), nil
}
func (handler *Handler) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.LoginResponse, error) {
	rsp, err := oauth2.ExchangeToken(oauth2.ExchangeTokenRequest{
		Code:         req.GetCode(),
		ClientId:     handler.Config.ClientId,
		ClientSecret: handler.Config.ClientSecret,
		GrantType:    handler.Config.GrantType,
		RedirectUri:  handler.Config.RedirectUri,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while exchanging token: %v", err)
	}
	info, err := oauth2.GetUserInfo("json", rsp.AccessToken)
	log.Info().Msgf("info : %v", info)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while retrieving user info: %v", err)
	}

	user, err := handler.Repo.GetUserByUserName(ctx, info.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = handler.Repo.CreateUser(ctx, convertUserInfo(info))
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error while creating user: %v", err)
			}

		} else {
			return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
		}
	}
	accessToken, err := handler.Jwt.GenerateJWT(user, time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}
	refreshToken, err := handler.Jwt.GenerateJWT(user, time.Hour*5)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}
	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
