package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tunes-anywhere/anywhere/config"
	"github.com/tunes-anywhere/anywhere/contrib/mongoutil"
	"github.com/tunes-anywhere/anywhere/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var trackLogger = config.Log.Named("track")

func trackCol() *mongo.Collection {
	return database.Database.Collection("tracks")
}

type TrackStatus string

const (
	TrackStatusUnknown TrackStatus = "UNKNOWN"
	TrackStatusPending TrackStatus = "PENDING"
	TrackStatusInvalid TrackStatus = "INVALID"
)

func (ts TrackStatus) Valid() error {
	if ts != TrackStatusUnknown &&
		ts != TrackStatusPending &&
		ts != TrackStatusInvalid {
		return fmt.Errorf("TrackStatus: invalid track status, must be one of: [%s]", strings.Join([]string{
			string(TrackStatusUnknown),
			string(TrackStatusPending),
			string(TrackStatusInvalid),
		}, ", "))
	}

	return nil
}

type PartialCreateTrack struct {
	Title     string `json:"title"`
	SourceURL string `json:"source_url"`
}

func (pct *PartialCreateTrack) Valid() error {
	if pct.Title == "" {
		return fmt.Errorf("PartialCreateTrack: must specify title")
	}

	if pct.SourceURL == "" {
		return fmt.Errorf("PartialCreateTrack: must specify source url")
	}

	return nil
}

type Track struct {
	ID          primitive.ObjectID    `json:"id"          bson:"_id"`
	Title       string                `json:"title"       bson:"title"`
	SourceURL   string                `json:"source_url"  bson:"source_url"`
	Duration    *int64                `json:"duration"    bson:"duration"`
	DataKey     *string               `json:"data_key"    bson:"data_key"`
	ArtistIDs   *[]primitive.ObjectID `json:"artist_ids"  bson:"artists"`
	TrackStatus TrackStatus           `json:"status"      bson:"status"`
	AccessedAt  int64                 `json:"accessed_at" bson:"accessed_at"`
	CreatedAt   int64                 `json:"created_at"  bson:"created_at"`
	UpdatedAt   int64                 `json:"updated_at"  bson:"updated_at"`
}

func (t *Track) Valid() error {
	if t.Title == "" {
		return fmt.Errorf("Track: Title: must specify a valid title")
	}
	if t.SourceURL == "" {
		return fmt.Errorf("Track: SourceURL: must specify a valid source url")
	}

	if t.Duration != nil && *t.Duration < 1 {
		return fmt.Errorf("Track: Duration: must specify a positive duration")
	}

	if t.DataKey != nil && len(*t.DataKey) < 1 {
		return fmt.Errorf("Track: DataKey: must specify a valid datakey")
	}

	if err := t.TrackStatus.Valid(); err != nil {
		return fmt.Errorf("Track: %s", err.Error())
	}

	return nil
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

func CreateTrack(ctx context.Context, pct *PartialCreateTrack) (*Track, error) {
	if err := pct.Valid(); err != nil {
		return nil, err
	}

	track := Track{
		ID:          primitive.NewObjectID(),
		Title:       pct.Title,
		SourceURL:   pct.SourceURL,
		Duration:    nil,
		DataKey:     nil,
		ArtistIDs:   nil,
		TrackStatus: TrackStatusUnknown,
		AccessedAt:  time.Now().Unix(),
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

	trackLogger.Debugw("read track", "id", id)

	return &track, nil
}

func UpdateTrack(ctx context.Context, id string, t *Track) (*Track, error) {
	if err := t.Valid(); err != nil {
		return nil, err
	}

	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "source_url", Value: t.SourceURL},
			{Key: "data_key", Value: t.DataKey},
			{Key: "duration", Value: t.Duration},
			{Key: "artist_ids", Value: t.ArtistIDs},
			{Key: "track_status", Value: t.TrackStatus},
			{Key: "updated_at", Value: time.Now().Unix()},
		},
	}}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var track Track
	if err := trackCol().FindOneAndUpdate(ctx, oid, update, opts).Decode(&track); err != nil {
		return nil, err
	}

	trackLogger.Debugw("updated track", "id", id)

	return &track, nil
}

func DeleteTrack(ctx context.Context, id string) error {
	if _, err := trackCol().DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return err
	}

	trackLogger.Debugw("deleted track", "id", id)

	return nil
}
