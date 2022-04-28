package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/nrmadi02/mini-project/app/config"
	"github.com/nrmadi02/mini-project/app/router"
	_ "github.com/nrmadi02/mini-project/docs"
	mid "github.com/nrmadi02/mini-project/user/delivery/http/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

// @title UMKM applications Documentation
// @description This is a UMKM management application
// @version 2.0
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization

func Run() {

	db := config.InitDB()

	e := echo.New()
	mid.NewGoMiddleware().LogMiddleware(e)

	router.SetupRouter(e, db)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	address := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
