package handler

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/oauth2"

	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	Jwt    config.JwtWrapper
	Repo   repository.IAuthRepo
	Config config.Config
}

func NewAuthHandler(jwt config.JwtWrapper, repo repository.IAuthRepo, config config.Config) *Handler {
	return &Handler{
		Jwt:    jwt,
		Repo:   repo,
		Config: config,
	}
}
func (handler *Handler) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

	hashPassword, err := config.HashedPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while hashing password: %s", err)
	}
	req.Password = hashPassword
	model := convertUser(req)

	user, err := handler.Repo.CreateUser(model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Errorf(codes.AlreadyExists, "email %s already register before", req.UserName)
		}
		return nil, status.Errorf(codes.Internal, "error while creating user: %s", err)
	}

	return convertUserResponse(user), nil
}
func (handler *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.Repo.GetUserByUserName(req.GetUserName())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, status.Errorf(codes.NotFound, "user with email %s not found", req.GetUserName())

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
	}
	return rsp, nil
}
func (handler *Handler) GetOneUser(ctx context.Context, req *pb.GetOneUserRequest) (*pb.UserResponse, error) {

	user, err := handler.Repo.FindOneUser(req.GetId())
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
	user, err := handler.Repo.FindOneUser(claims.Id)
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

	user, err := handler.Repo.GetUserByUserName(info.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = handler.Repo.CreateUser(convertUserInfo(info))
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
