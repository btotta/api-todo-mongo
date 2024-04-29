package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func CorsMiddleware() gin.HandlerFunc {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(allowedOrigins, ";")
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Authorization",
        "Refresh",
		"Content-Type",
		"Origin",
		"Referer",
	}
	config.ExposeHeaders = []string{"Access-Control-Allow-Origin"}
	config.AllowCredentials = true
	config.MaxAge = 300 * time.Second

	return cors.New(config)
}
