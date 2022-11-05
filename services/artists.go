package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tunes-anywhere/anywhere/config"
	"github.com/tunes-anywhere/anywhere/models"
	"github.com/tunes-anywhere/anywhere/server"
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
	raw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	var partialArtist models.PartialArtist
	if err := json.Unmarshal(raw, &partialArtist); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	artist, err := models.CreateArtist(c, &partialArtist)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	config.Log.Debugw("created artist",
		"partial_artist", partialArtist,
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
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	raw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	var partialArtist models.PartialArtist
	if err := json.Unmarshal(raw, &partialArtist); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	artist, err := models.UpdateArtist(c, id, &partialArtist)
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
