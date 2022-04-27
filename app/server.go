package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/nrmadi02/mini-project/app/config"
	"github.com/nrmadi02/mini-project/app/router"
	mid "github.com/nrmadi02/mini-project/user/delivery/http/middleware"
	"os"
)

func Run() {

	db := config.InitDB()

	e := echo.New()
	mid.NewGoMiddleware().LogMiddleware(e)

	router.SetupRouter(e, db)

	address := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
