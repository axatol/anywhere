package models

import (
	"context"
	"fmt"
	"time"

	"github.com/tunes-anywhere/anywhere/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func artistCol() *mongo.Collection {
	return database.Database.Collection("artists")
}

type PartialArtist struct {
	Name string `json:"name"`
}

func (pt *PartialArtist) Valid() error {
	if pt.Name == "" {
		return fmt.Errorf("PartialArtist: must specify name")
	}

	return nil
}

type Artist struct {
	ID        primitive.ObjectID `json:"id"         bson:"_id"`
	Name      string             `json:"name"       bson:"name"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
}

func ListArtists(ctx context.Context) ([]Artist, error) {
	var artists []Artist
	cur, err := artistCol().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func CreateArtist(ctx context.Context, pt *PartialArtist) (*Artist, error) {
	if err := pt.Valid(); err != nil {
		return nil, err
	}

	artist := Artist{
		ID:        primitive.NewObjectID(),
		Name:      pt.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if _, err := artistCol().InsertOne(ctx, artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

func ReadArtist(ctx context.Context, id string) (*Artist, error) {
	var artist Artist
	if err := artistCol().FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

func UpdateArtist(ctx context.Context, id string, pt *PartialArtist) (*Artist, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{
		{Key: "name", Value: pt.Name},
		{Key: "updated_at", Value: time.Now().Unix()},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var artist Artist
	if err := artistCol().FindOneAndUpdate(ctx, filter, update, opts).Decode(&artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

func DeleteArtist(ctx context.Context, id string) error {
	if _, err := artistCol().DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return err
	}

	return nil
}
