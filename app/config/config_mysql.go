package config

import (
	"fmt"
	"go_cicd/domain"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func InitDB() *gorm.DB {

	config := Config{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Name:     os.Getenv("DB_NAME"),
	}

	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",

		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	var err error
	DB, err = gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	InitialMigration()

	return DB
}

func InitialMigration() {
	DB.AutoMigrate(&domain.User{})
}
