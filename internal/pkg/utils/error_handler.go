package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorHandler struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`
}

func NewErrorHandler(message string, status int, timestamp string) *ErrorHandler {
	return &ErrorHandler{
		Message:   message,
		Status:    status,
		Timestamp: timestamp,
	}
}

func DefaultErrorResponse(c *gin.Context, status int, message string) {
	errorHandler := NewErrorHandler(message, status, time.Now().Format(time.RFC3339))
	c.JSON(status, errorHandler)
}
