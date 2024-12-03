package handler

import (
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/model"
	"github.com/vietquan-37/auth-service/pkg/model/enum"
	"github.com/vietquan-37/auth-service/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(req *pb.CreateUserRequest) *model.User {
	return &model.User{
		Username:    req.UserName,
		Password:    req.GetPassword(),
		PhoneNumber: req.GetPhoneNumber(),
		Role:        enum.UserRole,
	}
}

func convertUserResponse(user model.User) *pb.UserResponse {
	return &pb.UserResponse{
		UserId:      int32(user.ID),
		UserName:    user.PhoneNumber,
		PhoneNumber: user.PhoneNumber,
		Role:        string(user.Role),
		CreateAt:    timestamppb.New(user.CreatedAt),
	}

}
func convertValidate(claims *config.JwtClaims) *pb.ValidateResponse {
	return &pb.ValidateResponse{
		UserId: int32(claims.Id),
		Role:   string(claims.Role),
	}
}
