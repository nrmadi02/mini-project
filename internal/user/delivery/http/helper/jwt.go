package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/nrmadi02/mini-project/domain"
	"time"
)

type GoJWT struct {
}

func NewGoJWT() *GoJWT {
	return &GoJWT{}
}

func (j *GoJWT) CreateTokenJWT(user *domain.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID":    user.ID,
		"Roles":     user.Roles,
		"ExpiresAt": time.Now().Add(time.Hour * 48).Unix(),
	})

	return token.SignedString([]byte("220220"))
}
