package db

import (
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/product-service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConn(DbSource string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(DbSource), &gorm.Config{TranslateError: true},
	)
	err = db.AutoMigrate(model.Product{})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to migrate model:")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("fail to open timeout connection:")
		return nil
	}
	return db
}
