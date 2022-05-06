package config

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
	"time"
)

func ConnectMongo() (*mongo.Database, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_URL")).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	log.Info("success connect mongo")
	_, err = client.Database("logger").Collection("log").DeleteMany(context.Background(), bson.M{"level": "info"})
	if err != nil {
		log.Fatal(err.Error())
	}
	return client.Database("logger"), nil
}
