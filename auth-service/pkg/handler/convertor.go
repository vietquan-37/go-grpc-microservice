package handler

import (
	"github.com/vietquan-37/auth-service/pkg/model"
	"github.com/vietquan-37/auth-service/pkg/model/enum"
	"github.com/vietquan-37/auth-service/pkg/pb"
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
		UserName:    user.PhoneNumber,
		PhoneNumber: user.PhoneNumber,
		Role:        string(user.Role),
	}

}
