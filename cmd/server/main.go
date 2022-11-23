package main

import (
	"context"

	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/database"
	"github.com/axatol/anywhere/datastore"
	"github.com/axatol/anywhere/models"
	"github.com/axatol/anywhere/server"
	"github.com/axatol/anywhere/services"
)

func main() {
	ctx := context.Background()

	datastore.Init(ctx)

	database.Init(ctx)
	defer database.Close(ctx)

	models.Init(ctx)

	server.Init()

	server.R.GET("/api/health", services.Health)

	server.R.GET("/api/artists", services.ListArtists)
	server.R.GET("/api/artists/:id", services.ReadArtist)

	server.R.GET("/api/albums", services.ListAlbums)
	server.R.GET("/api/albums/:id", services.ReadAlbum)
	server.R.PUT("/api/albums/:id", services.UpdateAlbum)

	server.R.GET("/api/tracks", services.ListTracks)
	server.R.POST("/api/tracks", services.CreateTrack)
	server.R.GET("/api/tracks/:id", services.ReadTrack)
	server.R.PUT("/api/tracks/:id", services.UpdateTrack)
	server.R.DELETE("/api/tracks/:id", services.DeleteTrack)

	server.R.GET("/api/tracks/search", services.SearchTracks)
	server.R.GET("/api/tracks/metadata", services.SearchTrackMetadata)

	server.Start(config.Values.Server.Port)
}
