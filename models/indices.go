package models

import (
	"context"

	"github.com/aws/jsii-runtime-go"
	"github.com/axatol/anywhere/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(ctx context.Context) {
	uniqueOpt := &options.IndexOptions{Unique: jsii.Bool(true)}

	artistCol().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: 1}}, Options: uniqueOpt},
	})

	trackCol().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "title", Value: 1}}, Options: uniqueOpt},
		{Keys: bson.D{{Key: "source_url", Value: 1}}, Options: uniqueOpt},
		{Keys: bson.D{{Key: "artists", Value: 1}}},
	})

	config.Log.Info("initialised database models")
}
