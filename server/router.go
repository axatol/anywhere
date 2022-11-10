package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/axatol/anywhere/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

var (
	logger = config.Log.Named("server")
	R      *gin.Engine
)

func Init() {
	if config.Values.Debug {
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
	corsConfig.AllowOrigins = config.Values.Server.AllowOrigins
	R.Use(cors.New(corsConfig))
	R.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, false))
	R.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	R.Use(gzip.Gzip(gzip.DefaultCompression))
	R.Use(middlewareJWT())
}

func Start(port string) {
	logger.Debugw("starting server",
		"port", config.Values.Server.Port,
	)

	srv := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: R}
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalln(err)
	}
}
