package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/axatol/anywhere/datasource/musicbrainz"
	"github.com/axatol/anywhere/datasource/youtube"
	"github.com/axatol/anywhere/models"
	"github.com/axatol/anywhere/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ListTracks(c *gin.Context) {
	result, err := models.ListTracks(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	if result == nil {
		result = []models.Track{}
	}

	c.JSON(http.StatusOK, server.OKResponse(result))
}

func CreateTrack(c *gin.Context) {
	var input models.Track
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	track, err := models.CreateTrack(c, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, server.OKResponse(track))
}

func ReadTrack(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	track, err := models.ReadTrack(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(track))
}

func UpdateTrack(c *gin.Context) {
	id := c.Param("id")

	var input models.Track
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	track, err := models.UpdateTrack(c, id, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(track))
}

func DeleteTrack(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	if err := models.DeleteTrack(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	// TODO cleanup artists with no tracks
	// TODO cleanup albums with no tracks

	c.JSON(http.StatusNoContent, server.EmptyResponse())
}

func SearchTrackMetadata(c *gin.Context) {
	query := c.Query("query")
	if len(query) < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide query")))
		return
	}

	client, err := musicbrainz.NewClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	results, err := client.LookupRecording(query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(results))
}

func SearchTracks(c *gin.Context) {
	query := c.Query("query")
	if len(query) < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide query")))
		return
	}

	limit := []int{}
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.ParseInt(rawLimit, 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("invalid limit")))
			return
		}

		limit = append(limit, int(parsedLimit))
	}

	client, err := youtube.GetClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	results, err := client.Query(query, limit...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(results))
}
