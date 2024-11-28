package model

import (
	"github.com/vietquan-37/auth-service/pkg/model/enum"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	PhoneNumber string
	Role        enum.Role
}
