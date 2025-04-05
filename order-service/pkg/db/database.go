package db

import (
	"github.com/vietquan-37/order-service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func DbConn(DbSource string) *gorm.DB {

	db, err := gorm.Open(
		postgres.Open(DbSource), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("timeout connection error: %v", err)
	}
	err = db.AutoMigrate(model.Order{}, model.OrderDetail{})
	if err != nil {
		log.Fatalf("timeout migration error: %v", err)
	}
	return db
}
