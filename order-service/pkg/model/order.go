package model

import (
	"github.com/vietquan-37/order-service/pkg/enum"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	Amount      float64
	OrderDate   time.Time
	Status      enum.Status
	UserId      uint
	OrderDetail []OrderDetail `gorm:"foreignKey:OrderId;constraint:OnDelete:CASCADE;"`
}
