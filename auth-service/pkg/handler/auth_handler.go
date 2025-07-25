package handler

import (
	"common/cache"
	"common/kafka/producer"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/message"
	"github.com/vietquan-37/auth-service/pkg/model"
	"github.com/vietquan-37/auth-service/pkg/model/enum"
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
	redis    cache.Client
}

func NewAuthHandler(jwt config.JwtWrapper, repo repository.IAuthRepo, config config.Config, wait *sync.WaitGroup, producer *producer.Producer, redis cache.Client) *Handler {
	return &Handler{
		Jwt:      jwt,
		Repo:     repo,
		Config:   config,
		wg:       wait,
		redis:    redis,
		Producer: producer,
	}
}
func (handler *Handler) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

	hashPassword, err := config.HashedPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while hashing password: %s", err)
	}
	req.Password = hashPassword
	m := convertUser(req)

	user, err := handler.Repo.CreateUser(ctx, m)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Errorf(codes.AlreadyExists, "email %s already register before", req.UserName)
		}
		return nil, status.Errorf(codes.Internal, "error while creating user: %s", err)
	}

	token, err := handler.Jwt.GenerateJWT(user, time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating token: %v", err)
	}
	expiresAt := time.Now().Add(time.Hour).Unix()
	redisPayload := model.TokenPayload{
		Token:     token,
		UserID:    int32(user.ID),
		TokenType: enum.Verification,
		ExpiredAt: expiresAt,
	}
	err = handler.SaveVerificationToken(ctx, redisPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while saving token to redis: %s", err)
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
		return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
	}
	if err = config.CheckPassword(req.GetPassword(), user.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if !user.Active {
		return nil, status.Error(codes.Unauthenticated, "user is not active")
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
func (handler *Handler) SaveVerificationToken(ctx context.Context, payload model.TokenPayload) error {
	key := fmt.Sprintf("verify:%s", payload.Token)
	return handler.redis.Set(ctx, key, payload, time.Hour*2)
}
func (handler *Handler) GetVerificationToken(ctx context.Context, token string) (*model.TokenPayload, error) {
	key := fmt.Sprintf("verify:%s", token)
	val, err := handler.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var bytesVal []byte
	switch v := val.(type) {
	case string:
		bytesVal = []byte(v)
	case []byte:
		bytesVal = v
	default:
		return nil, fmt.Errorf("unexpected redis value type: %T", val)
	}

	var payload model.TokenPayload
	if err := json.Unmarshal(bytesVal, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (handler *Handler) VerifyAccount(ctx context.Context, req *pb.ValidateRequest) (*pb.CommonResponse, error) {
	payload, err := handler.GetVerificationToken(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, cache.ErrorCacheMiss) {
			return nil, status.Error(codes.InvalidArgument, "Invalid token")
		}
		return nil, status.Errorf(codes.Internal, "error while verifying token: %v", err)
	}
	if payload.TokenType != enum.Verification {
		return nil, status.Error(codes.InvalidArgument, "Invalid token")
	}
	user, err := handler.Repo.FindOneUser(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", payload.UserID)
		}
		return nil, status.Errorf(codes.Internal, "error while retrieving user: %v", err)
	}
	if time.Now().Unix() > payload.ExpiredAt {
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
		return nil, status.Error(codes.DeadlineExceeded, "Token has expired")
	}
	if user.Active {
		return nil, status.Error(codes.InvalidArgument, "User is already active")
	}
	user.Active = true
	err = handler.Repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while updating user: %v", err)
	}
	return &pb.CommonResponse{
		Message: "User verified successfully",
	}, nil

}
