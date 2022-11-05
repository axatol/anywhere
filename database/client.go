package database

import (
	"context"

	"github.com/tunes-anywhere/anywhere/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var Database *mongo.Database
var logger = config.Log.Named("database")

func Init(ctx context.Context) {
	opts := options.Client()
	opts.ApplyURI(config.Config.Database.Host)
	opts.Auth = &options.Credential{
		Username: config.Config.Database.User,
		Password: config.Config.Database.Pass,
	}

	newClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Fatalln(err)
	}

	client = newClient
	Database = client.Database(config.Config.Database.Name)
}

func Close(ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		logger.Errorln(err)
	}
}
