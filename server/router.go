package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/axatol/anywhere/config"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

var (
	R *gin.Engine
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

	R.Use(ginzap.Ginzap(config.Log.Desugar(), time.RFC3339, false))
	R.Use(ginzap.RecoveryWithZap(config.Log.Desugar(), true))
	R.Use(gzip.Gzip(gzip.DefaultCompression))

	R.Use(middlewareCORS())
	R.Use(middlewareJWT())

	config.Log.Info("configured server")
}

func Start(port string) {
	config.Log.Infow("starting server",
		"port", config.Values.Server.Port,
	)

	srv := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: R}
	if err := srv.ListenAndServe(); err != nil {
		config.Log.Fatalln(err)
	}
}
