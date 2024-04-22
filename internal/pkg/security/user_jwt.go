package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
)

type UserJWTInterface interface {
	GenerateToken(email string) (string, error)
	GenerateRefreshToken(email string) (string, error)
	ValidateToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, error)
	LogOff(token string)
}

var (
	secretKey      = os.Getenv("SECRET_KEY")
	refreshKey     = os.Getenv("REFRESH_KEY")
	secretTime, _  = strconv.Atoi(os.Getenv("SECRET_TIME"))
	refreshTime, _ = strconv.Atoi(os.Getenv("REFRESH_TIME"))
	logOffTokens   = cache.New(12*time.Hour, 10*time.Minute)
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(email string) (string, error) {
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(secretTime) * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(email string) (string, error) {
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(refreshTime) * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshKey))
}

func ValidateToken(token string) (string, error) {

	if IsLoggedOff(token) {
		return "", errors.New("invalid token")
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", err
	}

	return claims.Email, nil
}

func ValidateRefreshToken(token string) (string, error) {

	if IsLoggedOff(token) {
		return "", errors.New("invalid token")
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshKey), nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", err
	}

	return claims.Email, nil
}

func LogOff(token string) {
	logOffTokens.Set(token, true, cache.DefaultExpiration)
}

func IsLoggedOff(token string) bool {
	_, found := logOffTokens.Get(token)
	return found
}
