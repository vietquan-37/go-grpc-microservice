package model

import "github.com/vietquan-37/auth-service/pkg/model/enum"

type TokenPayload struct {
	Token     string    `json:"token"`
	UserID    int32     `json:"user_id"`
	TokenType enum.Type `json:"token_type"`
	ExpiredAt int64     `json:"expired_at"`
}
