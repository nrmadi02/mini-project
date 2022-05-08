package main

import (
	"github.com/joho/godotenv"
	"github.com/nrmadi02/mini-project/app"
	"github.com/nrmadi02/mini-project/app/config"
	"github.com/nrmadi02/mini-project/app/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	connectMongo, err := config.ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}
	mw := utils.NewMongoWriter(connectMongo)
	log.SetOutput(mw)
}

func main() {
	app.Run()
}
