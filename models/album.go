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

func albumCol() *mongo.Collection {
	return database.Database.Collection("albums")
}

type Album struct {
	ID        primitive.ObjectID   `json:"id"                   bson:"_id"                  validate:"required"`
	MBID      *string              `json:"mbid"                 bson:"mbid"                 validate:"required,min=1"`
	Name      string               `json:"name"                 bson:"name"                 validate:"required,min=1"`
	ArtistIDs []primitive.ObjectID `json:"artist_ids,omitempty" bson:"artist_ids,omitempty" validate:"omitempty"`
	TrackIDs  []primitive.ObjectID `json:"track_ids,omitempty"  bson:"track_ids,omitempty"  validate:"omitempty"`
	Metadata  map[string]string    `json:"metadata,omitempty"   bson:"metadata,omitempty"   validate:"omitempty"`
	CreatedAt int64                `json:"created_at"           bson:"created_at"           validate:"required"`
	UpdatedAt int64                `json:"updated_at"           bson:"updated_at"           validate:"required,gtefield=CreatedAt"`

	// TODO
	// CoverArt  *string              `json:"cover_art,omitempty"  bson:"cover_art,omitempty"  validate:"omitempty"`
}

func ListAlbums(ctx context.Context) ([]Album, error) {
	var albums []Album
	cur, err := albumCol().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &albums); err != nil {
		return nil, err
	}

	config.Log.Debugw("listed albums", "count", len(albums))

	return albums, nil
}

func CreateAlbum(ctx context.Context, a *Album) (*Album, error) {
	if err := validator.Validate.Struct(a); err != nil {
		return nil, err
	}

	album := Album{
		ID:        primitive.NewObjectID(),
		Name:      a.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if _, err := albumCol().InsertOne(ctx, album); err != nil {
		return nil, err
	}

	config.Log.Debugw("created album", "id", album.ID.Hex())

	return &album, nil
}

func ReadAlbum(ctx context.Context, id string) (*Album, error) {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return nil, err
	}

	var album Album
	if err := albumCol().FindOne(ctx, oid).Decode(&album); err != nil {
		return nil, err
	}

	config.Log.Debugw("read album", "id", album.ID.Hex())

	return &album, nil
}

func UpdateAlbum(ctx context.Context, id string, a *Album) (*Album, error) {
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

	if a.ArtistIDs != nil {
		updates = append(updates, bson.E{Key: "artist_ids", Value: a.ArtistIDs})
	}

	if a.TrackIDs != nil {
		updates = append(updates, bson.E{Key: "track_ids", Value: a.TrackIDs})
	}

	if a.Metadata != nil {
		updates = append(updates, bson.E{Key: "metadata", Value: a.Metadata})
	}

	update := bson.D{{Key: "$set", Value: updates}}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(false)

	var album Album
	if err := albumCol().FindOneAndUpdate(ctx, oid, update, opts).Decode(&album); err != nil {
		return nil, err
	}

	config.Log.Debugw("updated album", "id", album.ID.Hex())

	return &album, nil
}

func DeleteAlbum(ctx context.Context, id string) error {
	oid, err := mongoutil.BsonID(id)
	if err != nil {
		return err
	}

	var album Album
	if err = albumCol().FindOneAndDelete(ctx, oid).Decode(&album); err != nil {
		return err
	}

	config.Log.Debugw("deleted album", "id", album.ID.Hex())

	return nil
}
