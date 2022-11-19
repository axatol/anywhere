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
	server.R.POST("/api/artists", services.CreateArtist)
	server.R.GET("/api/artists/metadata", services.SearchArtistMetadata)
	server.R.GET("/api/artists/:id", services.ReadArtist)
	server.R.PUT("/api/artists/:id", services.UpdateArtist)
	server.R.DELETE("/api/artists/:id", services.DeleteArtist)

	server.R.GET("/api/tracks", services.ListTracks)
	server.R.POST("/api/tracks", services.CreateTrack)
	server.R.GET("/api/tracks/metadata", services.SearchTrackMetadata)
	server.R.GET("/api/tracks/:id", services.ReadTrack)
	server.R.PUT("/api/tracks/:id", services.UpdateTrack)
	server.R.DELETE("/api/tracks/:id", services.DeleteTrack)

	server.Start(config.Values.Server.Port)
}
