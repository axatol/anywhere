package services

import (
	"fmt"
	"net/http"

	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/models"
	"github.com/axatol/anywhere/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ListArtists(c *gin.Context) {
	result, err := models.ListArtists(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(result))
}

func CreateArtist(c *gin.Context) {
	var input models.PartialArtist
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	artist, err := models.CreateArtist(c, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	config.Log.Debugw("created artist",
		"partial_artist", input,
		"artist", artist,
	)

	c.JSON(http.StatusCreated, server.OKResponse(artist))
}

func ReadArtist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	artist, err := models.ReadArtist(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(artist))
}

func UpdateArtist(c *gin.Context) {
	id := c.Param("id")

	var input models.PartialArtist
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	artist, err := models.UpdateArtist(c, id, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(artist))
}

func DeleteArtist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	if err := models.DeleteArtist(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, server.EmptyResponse())
}
