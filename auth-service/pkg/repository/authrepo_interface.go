package repository

import "github.com/vietquan-37/auth-service/pkg/model"

type IAuthRepo interface {
	GetUserByUserName(username string) (*model.User, error)
	CreateUser(*model.User) (*model.User, error)
	FindOneUser(id int32) (*model.User, error)
}
