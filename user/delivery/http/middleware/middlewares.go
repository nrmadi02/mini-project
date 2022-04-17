package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

type GoMiddleware struct {
}

func NewGoMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}

func (m *GoMiddleware) LogMiddleware(e *echo.Echo) {
	e.Use(mid.LoggerWithConfig(mid.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency_human=${latency_human}\n",
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
