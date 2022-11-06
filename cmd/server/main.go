package main

import (
	"context"

	"github.com/tunes-anywhere/anywhere/config"
	"github.com/tunes-anywhere/anywhere/database"
	"github.com/tunes-anywhere/anywhere/datastore"
	"github.com/tunes-anywhere/anywhere/models"
	"github.com/tunes-anywhere/anywhere/server"
	"github.com/tunes-anywhere/anywhere/services"
)

func main() {
	ctx := context.Background()

	config.Log.Debug("initialising datastore client")
	datastore.Init(ctx)

	config.Log.Debug("initialising database client")
	database.Init(ctx)
	defer database.Close(ctx)

	config.Log.Debug("initialising database models")
	models.Init(ctx)

	config.Log.Debug("configuring server")
	server.Init()

	server.R.GET("/api/health", services.Health)

	server.R.GET("/api/artists", services.ListArtists)
	server.R.POST("/api/artists", services.CreateArtist)
	server.R.GET("/api/artists/:id", services.ReadArtist)
	server.R.PUT("/api/artists/:id", services.UpdateArtist)
	server.R.DELETE("/api/artists/:id", services.DeleteArtist)

	server.R.GET("/api/tracks", services.ListTracks)
	server.R.POST("/api/tracks", services.CreateTrack)
	server.R.GET("/api/tracks/:id", services.ReadTrack)
	server.R.PUT("/api/tracks/:id", services.UpdateTrack)
	server.R.DELETE("/api/tracks/:id", services.DeleteTrack)

	server.Start(config.Values.Server.Port)
}
