package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "todo-app-mongo/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {

	r := gin.Default()
	r.Use(cors.Default())

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/", s.healthHandler.HelloWorldHandler)
	r.GET("/health", s.healthHandler.HealthHandler)

	//User routes
	r.POST("/user", s.userHandler.Create)
	r.POST("/login", s.userHandler.Login)

	// Middleware
	r.Use(s.authMiddleware.AuthMiddleware())

	//Private routes
	r.PUT("/user", s.userHandler.Update)
	r.DELETE("/user", s.userHandler.Delete)
	r.GET("/user", s.userHandler.GetUser)
	r.POST("/refresh", s.userHandler.Refresh)
	r.POST("/logout", s.userHandler.Logout)

	//Todo routes
	s.todoRoutes(r)

	return r
}

func (s *Server) todoRoutes(r *gin.Engine) {
	r.GET("/todos", s.todoHandler.GetAll)
	r.GET("/todo/:id", s.todoHandler.Get)
	r.POST("/todo", s.todoHandler.Create)
	r.PUT("/todo/:id", s.todoHandler.Update)
	r.DELETE("/todo/:id", s.todoHandler.Delete)
}
