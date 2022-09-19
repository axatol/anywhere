package services

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getParamUint(c *gin.Context, idKey ...string) (uint, error) {
	key := "id"
	if len(idKey) > 0 {
		key = idKey[0]
	}

	raw := c.Param(key)
	if raw == "" {
		return 0, fmt.Errorf("param '%s' was not provided", key)
	}

	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
