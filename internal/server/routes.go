package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	docs "todo-app-mongo/docs"
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/handlers"
	"todo-app-mongo/internal/pkg/middleware"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// Initialize DAOs
	todoDao := database.NewTodoDAO(*s.db.GetDB())
	userDao := database.NewUserDAO(*s.db.GetDB())

	// Initialize Handlers
	healthHandler := handlers.NewHealthController(s.db)
	todoHandler := handlers.NewTodoHandler(todoDao, userDao)
	userHandler := handlers.NewUserHandler(userDao)

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Health routes
	r.GET("/", healthHandler.HelloWorldHandler)
	r.GET("/health", healthHandler.HealthHandler)

	//User routes
	user := r.Group("/user")
	{
		user.POST("", userHandler.Create)
		user.GET("/:id", middleware.AuthMiddleware(), userHandler.GetUser)
		user.PUT("/:id", middleware.AuthMiddleware(), userHandler.Update)
		user.DELETE("/:id", middleware.AuthMiddleware(), userHandler.Delete)

		//Auth routes
		user.POST("/login", userHandler.Login)
		user.POST("/refresh", userHandler.Refresh)
		user.POST("/logout", middleware.AuthMiddleware(), userHandler.Logout)
	}

	//Todo routes
	todo := r.Group("/todo")
	{
		todo.GET("/pagination", middleware.AuthMiddleware(), todoHandler.GetAll)
		todo.GET("/:id", middleware.AuthMiddleware(), todoHandler.Get)
		todo.POST("", middleware.AuthMiddleware(), todoHandler.Create)
		todo.PUT("/:id", middleware.AuthMiddleware(), todoHandler.Update)
		todo.DELETE("/:id", middleware.AuthMiddleware(), todoHandler.Delete)
	}

	return r
}
