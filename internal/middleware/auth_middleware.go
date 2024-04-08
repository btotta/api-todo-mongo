package middleware

import (
	"net/http"
	"os"
	"time"
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/models"
	"todo-app-mongo/internal/models/dtos"
	"todo-app-mongo/internal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var (
	secretKey    = os.Getenv("SECRET_KEY")
	refreshKey   = os.Getenv("REFRESH_KEY")
	logOffTokens = cache.New(60*time.Minute, 10*time.Minute)
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type AuthMiddlewareInterface interface {
	AuthMiddleware() gin.HandlerFunc
	GenerateTokenAndRefresh(user *models.User) (*dtos.UserLoginResponseDTO, error)
	RefreshToken(c *gin.Context) *dtos.UserLoginResponseDTO
	LogOff(c *gin.Context) bool
}

type authMiddleware struct {
	userRepository database.UserDAOInterface
}

func NewAuthMiddleware(userRepository database.UserDAOInterface) AuthMiddlewareInterface {
	return &authMiddleware{userRepository: userRepository}
}

// AuthMiddleware é um middleware para autenticar o usuário recebendo um token JWT da solicitação
func (a *authMiddleware) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Extrair o token JWT do cabeçalho de autorização
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Authorization header missing")
			c.Abort()
			return
		}
		// Parse e validar o token JWT
		tokenString := authHeader[len("Bearer "):]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		if !token.Valid {
			utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Verificar se o token está na lista de tokens inválidos
		if _, found := logOffTokens.Get(tokenString); found {
			utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Next()
	}
}

// GenerateTokenAndRefresh é um middleware para gerar um token JWT e um token de atualização para o usuário
func (a *authMiddleware) GenerateTokenAndRefresh(user *models.User) (*dtos.UserLoginResponseDTO, error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenExpirationTime := time.Now().Add(1 * time.Hour)

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	refreshToken.Claims = &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshKey))

	if err != nil {
		return nil, err
	}

	return &dtos.UserLoginResponseDTO{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// RefreshToken é um middleware para atualizar o token JWT
func (a *authMiddleware) RefreshToken(c *gin.Context) *dtos.UserLoginResponseDTO {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Authorization header missing")
		c.Abort()
	}

	refreshHeader := c.GetHeader("Refresh")
	if refreshHeader == "" {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Refresh header missing")
		c.Abort()
	}

	refreshTokenString := refreshHeader[len("Bearer "):]
	claims := &Claims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshKey), nil //
	})
	if err != nil {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		c.Abort()
	}

	if !refreshToken.Valid {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		c.Abort()
	}

	user, err := a.userRepository.GetByEmail(c, claims.Email)
	if err != nil {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid user")
		c.Abort()
	}

	tokenResponse, err := a.GenerateTokenAndRefresh(user)
	if err != nil {
		utils.DefaultErrorResponse(c, http.StatusInternalServerError, "Error generating tokens")
		c.Abort()
	}

	return tokenResponse
}

func (a *authMiddleware) LogOff(c *gin.Context) bool {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Authorization header missing")
		c.Abort()
		return false
	}

	tokenString := authHeader[len("Bearer "):]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		c.Abort()
		return false
	}

	if !token.Valid {
		utils.DefaultErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		c.Abort()
		return false
	}

	err = logOffTokens.Add(tokenString, true, cache.DefaultExpiration)
	if err != nil {
		utils.DefaultErrorResponse(c, http.StatusInternalServerError, "Error logging off")
		c.Abort()
		return false
	}

	return true
}
