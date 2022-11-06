package database

import (
	"context"

	"github.com/axatol/anywhere/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var Database *mongo.Database
var logger = config.Log.Named("database")

func Init(ctx context.Context) {
	opts := options.Client()
	opts.ApplyURI(config.Values.Database.Host)
	opts.Auth = &options.Credential{
		Username: config.Values.Database.User,
		Password: config.Values.Database.Pass,
	}

	newClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Fatalln(err)
	}

	client = newClient
	Database = client.Database(config.Values.Database.Name)
}

func Close(ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		logger.Errorln(err)
	}
}
