package repository

import (
	"context"
	"github.com/vietquan-37/auth-service/pkg/model"
)

type IAuthRepo interface {
	GetUserByUserName(ctx context.Context, username string) (*model.User, error)
	CreateUser(context.Context, *model.User) (*model.User, error)
	FindOneUser(ctx context.Context, id int32) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
}
