package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/nrmadi02/mini-project/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type GoJWT struct {
}

func NewGoJWT() *GoJWT {
	return &GoJWT{}
}

func (j *GoJWT) CreateTokenJWT(user *domain.User) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID":    user.ID,
		"Roles":     user.Roles,
		"ExpiresAt": time.Now().Add(time.Hour * 48).Unix(),
	})

	fixToken, err := token.SignedString([]byte("220220"))
	if err != nil {
		panic(err.Error())
		log.Info("error create token jwt")
	}

	return fixToken
}
