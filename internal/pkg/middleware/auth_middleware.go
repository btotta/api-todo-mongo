package middleware

import (
	"strings"
	"todo-app-mongo/internal/pkg/security"
	"todo-app-mongo/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			utils.DefaultErrorResponse(c, 401, "Unauthorized")
			c.Abort()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		email, err := security.ValidateToken(tokenString)

		if err != nil {
			utils.DefaultErrorResponse(c, 401, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("email", email)

		c.Next()
	}
}
