package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	docs "todo-app-mongo/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/", s.healthHandler.HelloWorldHandler)
	r.GET("/health", s.healthHandler.HealthHandler)

	// Todo routes	
	r.GET("/todos", s.todoHandler.GetAll)
	r.GET("/todo/:id", s.todoHandler.Get)
	r.POST("/todo", s.todoHandler.Create)
	r.PUT("/todo/:id", s.todoHandler.Update)
	r.DELETE("/todo/:id", s.todoHandler.Delete)

	return r
}
