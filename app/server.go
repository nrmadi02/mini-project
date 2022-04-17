package app

import (
	"fmt"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go_cicd/app/config"
	_userController "go_cicd/user/delivery/http"
	mid "go_cicd/user/delivery/http/middleware"
	"go_cicd/user/repository"
	"go_cicd/user/usecase"
	"os"
)

func Run() {
	client := appinsights.NewTelemetryClient(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))
	request := appinsights.NewRequestTelemetry("GET", "https://myapp.azurewebsites.net/", 1, "Success")
	client.Track(request)

	db := config.InitDB()

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)

	e := echo.New()
	mid.NewGoMiddleware().LogMiddleware(e)
	_userController.NewUserController(e, userUsecase)
	address := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
