package handler

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error" example:"error description"`
}

func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: message})
}
