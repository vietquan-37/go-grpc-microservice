package repository

import (
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

func NewAuthRepo(db *gorm.DB) IAuthRepo {
	repo := &AuthRepo{DB: db}
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file: ")
	}
	user, err := repo.GetUserByUserName(c.AdminUserName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal().Err(err).Msg("error while checking user:")
	}
	if user == nil {
		hashPassword, err := config.HashedPassword(c.AdminPassword)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to hash password:")
		}
		m := &model.User{

			Username:    c.AdminUserName,
			Password:    hashPassword,
			PhoneNumber: "0912021638",
			Role:        enum.AdminRole,
		}
		_, err = repo.CreateUser(m)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to create admin user: ")
		}
	}
	return repo
}

func (repo *AuthRepo) GetUserByUserName(username string) (*model.User, error) {
	var dbUser model.User
	err := repo.DB.Where("username = ?", username).First(&dbUser).Error
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}
func (repo *AuthRepo) CreateUser(user *model.User) (*model.User, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *AuthRepo) FindOneUser(id int32) (user *model.User, err error) {
	err = repo.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
