package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"todo-app-mongo/internal/controllers"
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/middleware"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port           int
	db             database.Service
	healthHandler  controllers.HealthHandlerInterface
	todoHandler    controllers.TodoHandlerInterface
	userHandler    controllers.UserHandlerInterface
	authMiddleware middleware.AuthMiddlewareInterface
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()

	// Initialize DAOs
	todoDao := database.NewTodoDAO(*db.GetDB())
	userDao := database.NewUserDAO(*db.GetDB())

	// Initialize Handlers
	healthHandler := controllers.NewHealthController(db)
	todoHandler := controllers.NewTodoHandler(todoDao, userDao)
	authMiddleware := middleware.NewAuthMiddleware(userDao)
	userHandler := controllers.NewUserHandler(userDao, authMiddleware)

	NewServer := &Server{
		port:           port,
		db:             database.New(),
		healthHandler:  healthHandler,
		todoHandler:    todoHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
