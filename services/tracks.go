package services

import (
	"fmt"
	"net/http"

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

	c.JSON(http.StatusOK, server.OKResponse(result))
}

func CreateTrack(c *gin.Context) {
	var input models.PartialCreateTrack
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

	var input models.PartialUpdateTrack
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

	c.JSON(http.StatusNoContent, server.EmptyResponse())
}
