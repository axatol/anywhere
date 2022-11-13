package database

import (
	"context"

	"github.com/axatol/anywhere/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	Database *mongo.Database
)

func Init(ctx context.Context) {
	opts := options.Client()
	opts.ApplyURI(config.Values.Database.Host)
	opts.SetAuth(options.Credential{
		Username: config.Values.Database.User,
		Password: config.Values.Database.Pass,
	})

	newClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		config.Log.Fatalln(err)
	}

	client = newClient
	Database = client.Database(config.Values.Database.Name)

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	config.Log.Infow("initialised database client",
		"name", config.Values.Database.Name,
	)
}

func Close(ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		config.Log.Errorln(err)
	}
}
