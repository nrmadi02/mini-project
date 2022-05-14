package main

import (
	"github.com/nrmadi02/mini-project/app"
	"github.com/nrmadi02/mini-project/app/config"
	"github.com/nrmadi02/mini-project/app/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	
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
