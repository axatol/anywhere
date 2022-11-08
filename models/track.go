package models

import (
	"context"
	"time"

	"github.com/aws/jsii-runtime-go"
	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/contrib/mongoutil"
	"github.com/axatol/anywhere/contrib/validator"
	"github.com/axatol/anywhere/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var trackLogger = config.Log.Named("track")

func trackCol() *mongo.Collection {
	return database.Database.Collection("tracks")
}

const (
	TrackStatusUnknown = "UNKNOWN"
	TrackStatusPending = "PENDING"
	TrackStatusInvalid = "INVALID"
)

type PartialCreateTrack struct {
	SourceURL string `json:"source_url"   bson:"source_url" validate:"required,url"`
	Title     string `json:"title"        bson:"title"      validate:"required,min=3"`
}

type PartialUpdateTrack struct {
	Duration    *int64                `json:"duration,omitempty"     bson:"duration"     validate:"omitempty,min=1"`
	DataKey     *string               `json:"data_key,omitempty"     bson:"data_key"     validate:"omitempty,min=1"`
	ArtistIDs   *[]primitive.ObjectID `json:"artist_ids,omitempty"   bson:"artist_ids"   validate:"omitempty,min=1"`
	TrackStatus *string               `json:"track_status,omitempty" bson:"track_status" validate:"omitempty,oneof=UNKNOWN PENDING INVALID"`
}

type Track struct {
	ID          primitive.ObjectID    `json:"id"                     bson:"_id"          validate:"required"`
	SourceURL   string                `json:"source_url"             bson:"source_url"   validate:"required,url"`
	Title       string                `json:"title"                  bson:"title"        validate:"required,min=3"`
	Duration    *int64                `json:"duration,omitempty"     bson:"duration"     validate:"omitempty,min=1"`
	DataKey     *string               `json:"data_key,omitempty"     bson:"data_key"     validate:"omitempty,min=1"`
	ArtistIDs   *[]primitive.ObjectID `json:"artist_ids,omitempty"   bson:"artist_ids"   validate:"omitempty,min=1"`
	TrackStatus *string               `json:"track_status,omitempty" bson:"track_status" validate:"omitempty,oneof=UNKNOWN PENDING INVALID"`
	AccessedAt  int64                 `json:"accessed_at"            bson:"accessed_at"  validate:"required"`
	CreatedAt   int64                 `json:"created_at"             bson:"created_at"   validate:"required"`
	UpdatedAt   int64                 `json:"updated_at"             bson:"updated_at"   validate:"required,gtefield=CreatedAt"`
}

func ListTracks(ctx context.Context) ([]Track, error) {
	var tracks []Track
	cur, err := trackCol().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &tracks); err != nil {
		return nil, err
	}

	trackLogger.Debugw("listed tracks", "count", len(tracks))

	return tracks, nil
}

func CreateTrack(ctx context.Context, t *PartialCreateTrack) (*Track, error) {
	if err := validator.Validate.Struct(t); err != nil {
		return nil, err
	}

	track := Track{
		ID:          primitive.NewObjectID(),
		Title:       t.Title,
		SourceURL:   t.SourceURL,
		Duration:    nil,
		DataKey:     nil,
		ArtistIDs:   nil,
		TrackStatus: jsii.String(TrackStatusUnknown),
		AccessedAt:  0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if _, err := trackCol().InsertOne(ctx, track); err != nil {
		return nil, err
	}

	trackLogger.Debugw("created track", "id", track.ID.Hex())

	return &track, nil
}

func ReadTrack(ctx context.Context, id string) (*Track, error) {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "accessed_at", Value: time.Now().Unix()},
		},
	}}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var track Track
	if err := trackCol().FindOneAndUpdate(ctx, oid, update, opts).Decode(&track); err != nil {
		return nil, err
	}

	trackLogger.Debugw("read track", "id", track.ID.Hex())

	return &track, nil
}

func UpdateTrack(ctx context.Context, id string, t *PartialUpdateTrack) (*Track, error) {
	if err := validator.Validate.Struct(t); err != nil {
		return nil, err
	}

	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	updates := bson.D{{Key: "updated_at", Value: time.Now().Unix()}}
	if t.Duration != nil {
		updates = append(updates, bson.E{Key: "duration", Value: t.Duration})
	}

	if t.DataKey != nil {
		updates = append(updates, bson.E{Key: "data_key", Value: t.DataKey})
	}

	if t.ArtistIDs != nil {
		updates = append(updates, bson.E{Key: "artist_ids", Value: t.ArtistIDs})
	}

	if t.TrackStatus != nil {
		updates = append(updates, bson.E{Key: "track_status", Value: t.TrackStatus})
	}

	update := bson.D{{Key: "$set", Value: updates}}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var track Track
	if err := trackCol().FindOneAndUpdate(ctx, oid, update, opts).Decode(&track); err != nil {
		return nil, err
	}

	trackLogger.Debugw("updated track", "id", track.ID.Hex())

	return &track, nil
}

func DeleteTrack(ctx context.Context, id string) error {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return err
	}

	var track Track
	if err := trackCol().FindOneAndDelete(ctx, oid).Decode(&track); err != nil {
		return err
	}

	trackLogger.Debugw("deleted track", "id", track.ID.Hex())

	return nil
}
