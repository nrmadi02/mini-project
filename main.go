package main

import (
	"github.com/joho/godotenv"
	"github.com/nrmadi02/mini-project/app"
	"log"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app.Run()
}
