package db

import (
	"github.com/rs/zerolog/log"

	"github.com/vietquan-37/auth-service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConn(DbSource string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(DbSource), &gorm.Config{TranslateError: true},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file: ")
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file: ")

		return nil
	}
	return db
}
