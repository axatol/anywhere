package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tunes-anywhere/anywhere/config"
)

var logger = config.Logger.Named("server")
var R *gin.Engine

func Init() {
	if config.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	R = gin.New()
	R.RedirectTrailingSlash = true
	R.RedirectFixedPath = false
	R.HandleMethodNotAllowed = false
	R.ForwardedByClientIP = true
	R.UseRawPath = false
	R.UnescapePathValues = true
	R.RemoveExtraSlash = true

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowOrigins = config.Config.Server.AllowOrigins
	R.Use(cors.New(corsConfig))
	R.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, false))
	R.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	R.Use(gzip.Gzip(gzip.DefaultCompression))
}

func Start(port int) {
	logger.Debugw("starting server",
		"port", config.Config.Server.Port,
	)

	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: R}
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalln(err)
	}
}
