package utils

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type MongoWriter interface {
	Write(p []byte) (n int, err error)
}

type mongoWriter struct {
	db *mongo.Database
}

func NewMongoWriter(db *mongo.Database) MongoWriter {
	return &mongoWriter{
		db: db,
	}
}

func (mw *mongoWriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	ctx := context.Background()
	if err != nil {
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}
	_, err = mw.db.Collection("log").InsertOne(ctx, data)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
