package config

import (
	"fmt"
	"github.com/nrmadi02/mini-project/app/utils"
	"github.com/nrmadi02/mini-project/db/seeds"
	"github.com/nrmadi02/mini-project/domain"
	log "github.com/sirupsen/logrus"
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

	cm, err := ConnectMongo()
	if err != nil {
		log.Fatal(err.Error())
	}

	DB, err = gorm.Open(sqlserver.Open(connectionString), &gorm.Config{
		Logger: utils.SlowLoggerGorm(cm),
	})
	if err != nil {
		panic(err)
	}

	InitialMigration()

	return DB
}

func InitialMigration() {
	err := DB.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.Tag{}, &domain.Enterprise{}, &domain.RatingEnterprise{}, &domain.Favorite{}, &domain.Review{})

	if err != nil {
		panic("could not connect to db " + err.Error())
		return
	}
	seeds.Execute(DB)
}
