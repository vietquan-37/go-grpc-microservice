package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/model"
	"github.com/vietquan-37/auth-service/pkg/model/enum"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB, adminUsername, adminPassword string) IAuthRepo {
	repo := &AuthRepo{DB: db}
	ctx := context.Background()
	user, err := repo.GetUserByUserName(ctx, adminUsername)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal().Err(err).Msg("error while checking user:")
	}
	if user == nil {
		hashPassword, err := config.HashedPassword(adminPassword)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to hash password:")
		}
		m := &model.User{

			Username:    adminUsername,
			Password:    hashPassword,
			PhoneNumber: "0912021638",
			Role:        enum.AdminRole,
		}
		_, err = repo.CreateUser(ctx, m)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to create admin user: ")
		}
	}
	return repo
}

func (repo *AuthRepo) GetUserByUserName(ctx context.Context, username string) (*model.User, error) {
	//time.Sleep(10 * time.Second)
	var dbUser model.User
	err := repo.DB.WithContext(ctx).Where("username = ?", username).First(&dbUser).Error
	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}
func (repo *AuthRepo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := repo.DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *AuthRepo) FindOneUser(ctx context.Context, id int32) (user *model.User, err error) {
	//time.Sleep(time.Second * 7)
	err = repo.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
