package app

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmadi02/mini-project/app/config"
	"github.com/nrmadi02/mini-project/app/router"
	docs "github.com/nrmadi02/mini-project/docs"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"time"
)

// @title UMKM applications Documentation
// @description This is a UMKM management application
// @version 2.0
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization

func Run() {
	docs.SwaggerInfo.Host = os.Getenv("APP_HOST")

	db := config.InitDB()

	e := echo.New()
	e.Use(loggingMiddleware())
	router.SetupRouter(e, db)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	address := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}

func loggingMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		var res map[string]interface{}
		_ = json.Unmarshal(resBody, &res)
		start := time.Now()
		log.WithFields(log.Fields{
			"method":     c.Request().Method,
			"path":       c.Path(),
			"status":     c.Response().Status,
			"latency_ms": time.Since(start).Milliseconds(),
		}).Info(res["message"])
	})
}
