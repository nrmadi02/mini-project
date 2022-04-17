package helper

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type MyClaims struct {
	UserID int    `json:"userID"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type GoJWT struct {
}

func NewGoJWT() *GoJWT {
	return &GoJWT{}
}

func (j *GoJWT) CreateTokenJWT(userID int, email string) (string, error) {
	claims := MyClaims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("220220"))
}
