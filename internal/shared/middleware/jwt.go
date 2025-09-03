package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"bankmore/internal/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AccountID     string `json:"accountId"`
	AccountNumber string `json:"accountNumber"`
	CPF           string `json:"cpf"`
	jwt.RegisteredClaims
}

func GenerateJWT(accountID, accountNumber, cpf string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here-change-in-production"
	}

	claims := Claims{
		AccountID:     accountID,
		AccountNumber: accountNumber,
		CPF:           cpf,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Type:    models.ErrorUserUnauthorized,
				Message: "Token de autorização não fornecido",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Type:    models.ErrorUserUnauthorized,
				Message: "Formato de token inválido",
			})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key-here-change-in-production"
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Type:    models.ErrorUserUnauthorized,
				Message: "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("accountId", claims.AccountID)
			c.Set("accountNumber", claims.AccountNumber)
			c.Set("cpf", claims.CPF)
		}

		c.Next()
	}
}
