package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, err error, statusCode int) {
	log.Printf("ERROR: %s", err.Error())
	c.AbortWithStatusJSON(statusCode, map[string]string{
		"error": err.Error(),
	})
}
