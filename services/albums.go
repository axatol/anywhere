package services

import (
	"fmt"
	"net/http"

	"github.com/axatol/anywhere/models"
	"github.com/axatol/anywhere/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ListAlbums(c *gin.Context) {
	result, err := models.ListAlbums(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	if result == nil {
		result = []models.Album{}
	}

	c.JSON(http.StatusOK, server.OKResponse(result))
}

func CreateAlbum(c *gin.Context) {
	var input models.Album
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	album, err := models.CreateAlbum(c, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, server.OKResponse(album))
}

func ReadAlbum(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	album, err := models.ReadAlbum(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(album))
}

func UpdateAlbum(c *gin.Context) {
	id := c.Param("id")

	var input models.Album
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	album, err := models.UpdateAlbum(c, id, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, server.OKResponse(album))
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrorResponse(fmt.Errorf("must provide id")))
		return
	}

	if err := models.DeleteAlbum(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, server.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, server.EmptyResponse())
}
