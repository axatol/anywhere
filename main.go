package main

import (
	"github.com/tunes-anywhere/anywhere/config"
	"github.com/tunes-anywhere/anywhere/database"
	"github.com/tunes-anywhere/anywhere/models"
	"github.com/tunes-anywhere/anywhere/server"
	"github.com/tunes-anywhere/anywhere/services"
)

func main() {
	config.Logger.Debug(config.Config)

	server.Init()
	database.Init()
	models.Migrate(database.DB)

	api := server.R.Group("/api")
	api.GET("/health", services.Health)

	api.GET("/artists", services.ListArtists)
	api.POST("/artists", services.CreateArtist)
	api.GET("/artists/:id", services.ReadArtist)
	api.PUT("/artists/:id", services.UpdateArtist)
	api.DELETE("/artists/:id", services.DeleteArtist)

	api.GET("/tracks", services.ListTracks)
	api.POST("/tracks", services.CreateTrack)
	api.GET("/tracks/:id", services.ReadTrack)
	api.PUT("/tracks/:id", services.UpdateTrack)
	api.DELETE("/tracks/:id", services.DeleteTrack)

	server.Start(config.Config.Server.Port)
}
