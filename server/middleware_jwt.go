package server

import (
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/axatol/anywhere/config"
	"github.com/gin-gonic/gin"
)

var jwtLogger = config.Log.Named("jwtmiddleware")

func middlewareJWT() gin.HandlerFunc {
	if config.Values.Server.Auth.Issuer == "" ||
		config.Values.Server.Auth.Audience == "" {
		return nil
	}

	issuer, err := url.Parse(config.Values.Server.Auth.Issuer)
	if err != nil {
		panic(err)
	}

	provider := jwks.NewCachingProvider(issuer, time.Minute*5)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		config.Values.Server.Auth.Issuer,
		[]string{config.Values.Server.Auth.Audience},
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			jwtLogger.Warnw("request denied", "path", c.Request.URL.Path, "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler))

		middleware.CheckJWT(httpHandlerAdapter{c}).ServeHTTP(c.Writer, c.Request)
	}
}

type httpHandlerAdapter struct {
	c *gin.Context
}

func (h httpHandlerAdapter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.c.Next()
}
