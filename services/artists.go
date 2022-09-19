package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tunes-anywhere/anywhere/models"
	"github.com/tunes-anywhere/anywhere/server"
)

func ListArtists(c *gin.Context) {
	result, err := models.ListArtists()
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

	if partialArtist.Name == "" {
		err := fmt.Errorf("must provide artist name")
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	artist, err := models.CreateArtist(partialArtist.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, server.OKResponse(artist))
}

func ReadArtist(c *gin.Context) {
	id, err := getParamUint(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	artist, err := models.ReadArtist(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(artist))
}

func UpdateArtist(c *gin.Context) {
	id, err := getParamUint(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
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

	if partialArtist.Name == "" {
		err := fmt.Errorf("must provide artist name")
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	artist, err := models.UpdateArtist(id, &partialArtist)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(artist))
}

func DeleteArtist(c *gin.Context) {
	id, err := getParamUint(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	if err := models.DeleteArtist(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, server.EmptyResponse())
}
