package models

import (
	"context"
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/axatol/anywhere/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createIndex(col string, name string, unique bool) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.D{{Key: name, Value: 1}},
		Options: &options.IndexOptions{
			Name:   jsii.String(fmt.Sprintf("%s.%s", col, name)),
			Unique: jsii.Bool(unique),
		},
	}

}

func Init(ctx context.Context) {
	albumCol().Indexes().CreateMany(ctx, []mongo.IndexModel{
		createIndex(albumCol().Name(), "mbid", true),
		createIndex(albumCol().Name(), "name", true),
	})

	artistCol().Indexes().CreateMany(ctx, []mongo.IndexModel{
		createIndex(artistCol().Name(), "mbid", true),
		createIndex(artistCol().Name(), "name", true),
	})

	trackCol().Indexes().CreateMany(ctx, []mongo.IndexModel{
		createIndex(trackCol().Name(), "mbid", true),
		createIndex(trackCol().Name(), "source_url", true),
		createIndex(trackCol().Name(), "name", true),
		createIndex(trackCol().Name(), "data_key", true),
	})

	config.Log.Info("initialised database models")
}
