package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vietquan-37/auth-service/pkg/model"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("expired token")

type JwtWrapper struct {
	JwtSecret []byte
}
type JwtClaims struct {
	jwt.StandardClaims
	Id    int32
	Role  string
	Email string
}

func NewJwtWrapper(secret string) (*JwtWrapper, error) {
	if secret == "" {
		return nil, errors.New("no secret provide")
	}
	return &JwtWrapper{
		JwtSecret: []byte(secret),
	}, nil
}
func (wrapper *JwtWrapper) GenerateJWT(user *model.User, expiration time.Duration) (tokenString string, err error) {
	expirationTime := time.Now().Add(expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		Id:    int32(user.ID),
		Role:  string(user.Role),
		Email: user.Username,
	})
	signed, err := token.SignedString(wrapper.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return signed, nil
}
func (wrapper *JwtWrapper) ValidateToken(token string) (claims *JwtClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return wrapper.JwtSecret, nil
	})
	if err != nil {
		return nil, errors.Join(ErrInvalidToken, err)
	}
	claims, ok := t.Claims.(*JwtClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpiredToken
	}
	return claims, nil
}
