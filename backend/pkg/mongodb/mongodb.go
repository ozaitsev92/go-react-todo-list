package mongodb

import (
	"context"
	"sync"

	"github.com/ozaitsev92/tododdd/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var m *mongo.Database
var hdlOnce sync.Once

// NewOrGetSingleton -.
func NewOrGetSingleton(config config.Config) *mongo.Database {
	hdlOnce.Do(func() {
		db, err := initMongo(config)
		if err != nil {
			panic(err)
		}

		m = db
	})

	return m
}

func initMongo(config config.Config) (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoUrl))
	if err != nil {
		return nil, err
	}

	return client.Database(config.MongoDBName), nil
}
