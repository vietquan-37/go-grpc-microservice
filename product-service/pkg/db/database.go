package db

import (
	"github.com/vietquan-37/product-service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func DbConn(DbSource string) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(DbSource), &gorm.Config{TranslateError: true},
	)
	err = db.AutoMigrate(model.Product{})
	if err != nil {
		log.Fatalf("err while migrating model %v", err)
	}
	if err != nil {
		log.Fatalf("database connection error:%s", err)
		return nil
	}
	return db
}
