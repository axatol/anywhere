package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "OK",
		"time":   time.Now().Format(time.RFC3339),
	})
}
