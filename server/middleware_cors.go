package server

import (
	"github.com/axatol/anywhere/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func middlewareCORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowHeaders = []string{"authorization", "x-timestamp"}
	corsConfig.AllowOrigins = config.Values.Server.AllowOrigins

	return cors.New(corsConfig)
}
