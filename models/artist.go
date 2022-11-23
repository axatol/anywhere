package models

import (
	"context"
	"time"

	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/contrib/mongoutil"
	"github.com/axatol/anywhere/contrib/validator"
	"github.com/axatol/anywhere/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func artistCol() *mongo.Collection {
	return database.Database.Collection("artists")
}

type Artist struct {
	ID        primitive.ObjectID `json:"id"                 bson:"_id"                validate:"required"`
	MBID      *string            `json:"mbid,omitempty"     bson:"mbid,omitempty"     validate:"omitempty"`
	Name      string             `json:"name"               bson:"name"               validate:"required,min=1"`
	Metadata  map[string]string  `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`
	CreatedAt int64              `json:"created_at"         bson:"created_at"         validate:"required"`
	UpdatedAt int64              `json:"updated_at"         bson:"updated_at"         validate:"required,gtefield=CreatedAt"`
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

	config.Log.Debugw("listed artists", "count", len(artists))

	return artists, nil
}

func CreateArtist(ctx context.Context, a *Artist) (*Artist, error) {
	if err := validator.Validate.Struct(a); err != nil {
		return nil, err
	}

	artist := Artist{
		ID:        primitive.NewObjectID(),
		Name:      a.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if _, err := artistCol().InsertOne(ctx, artist); err != nil {
		return nil, err
	}

	config.Log.Debugw("created artist", "id", artist.ID.Hex())

	return &artist, nil
}

func ReadArtist(ctx context.Context, id string) (*Artist, error) {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	var artist Artist
	if err := artistCol().FindOne(ctx, oid).Decode(&artist); err != nil {
		return nil, err
	}

	config.Log.Debugw("read artist", "id", artist.ID.Hex())

	return &artist, nil
}

func UpdateArtist(ctx context.Context, id string, a *Artist) (*Artist, error) {
	if err := validator.Validate.Struct(a); err != nil {
		return nil, err
	}

	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	updates := bson.D{
		{Key: "name", Value: a.Name},
		{Key: "updated_at", Value: time.Now().Unix()},
	}

	if a.MBID != nil {
		updates = append(updates, bson.E{Key: "mbid", Value: a.MBID})
	}

	if a.Metadata != nil {
		updates = append(updates, bson.E{Key: "metadata", Value: a.Metadata})
	}

	update := bson.D{{Key: "$set", Value: updates}}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var artist Artist
	if err := artistCol().FindOneAndUpdate(ctx, oid, update, opts).Decode(&artist); err != nil {
		return nil, err
	}

	config.Log.Debugw("updated artist", "id", artist.ID.Hex())

	return &artist, nil
}

func DeleteArtist(ctx context.Context, id string) error {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return err
	}

	var artist Artist
	if err = artistCol().FindOneAndDelete(ctx, oid).Decode(&artist); err != nil {
		return err
	}

	config.Log.Debugw("deleted artist", "id", artist.ID.Hex())

	return nil
}
