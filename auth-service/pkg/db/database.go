package db

import (
	"log"

	"github.com/vietquan-37/auth-service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConn(DbSource string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(DbSource), &gorm.Config{TranslateError: true},
	)
	if err != nil {
		log.Fatalf("fail to connect to db: %v", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("error while migrating user:%v", err)

		return nil
	}
	return db
}
