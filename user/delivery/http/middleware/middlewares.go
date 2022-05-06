package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

type GoMiddleware struct {
}

func NewGoMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}

func (m *GoMiddleware) LogMiddleware(e *echo.Echo) {
	fp, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.New(fp, "\r\n", log.LstdFlags)
	e.Use(mid.LoggerWithConfig(mid.LoggerConfig{
		Format:           "[${time_custom}] ${status} ${method} ${uri} ${latency_human}\n",
		CustomTimeFormat: "15:04:05, 02-01-2006",
		Output:           fp,
	}))
}

func (m *GoMiddleware) AuthMiddleware() echo.MiddlewareFunc {
	signingKey := []byte("220220")

	config := mid.JWTConfig{
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingKey, nil
			}

			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
	}

	return mid.JWTWithConfig(config)
}
