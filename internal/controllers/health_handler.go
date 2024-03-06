package controllers

import (
	"net/http"
	"todo-app-mongo/internal/database"

	"github.com/gin-gonic/gin"
)

type HealthHandlerInterface interface {
	HealthHandler(c *gin.Context)
	HelloWorldHandler(c *gin.Context)
}

type healthHandler struct {
	db database.Service
}

func NewHealthController(db database.Service) HealthHandlerInterface {
	return &healthHandler{db: db}
}

// @Summary Health check
// @Description Check if the server is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /health [get]
func (h *healthHandler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.db.Health())
}

// @Summary HelloWorld
// @Description Main entry point
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router / [get]
func (h *healthHandler) HelloWorldHandler(c *gin.Context) {
	body := map[string]string{
		"message": "Hello World",
	}

	c.JSON(http.StatusOK, body)
}
