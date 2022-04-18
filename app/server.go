package app

import (
	"fmt"
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
