package handler

import (
	"github.com/gin-gonic/gin"
)

type errorAcc struct {
	Message string `json: error-message`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	//log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorAcc{message})
}
