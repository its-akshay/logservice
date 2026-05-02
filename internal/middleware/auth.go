package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/logservice/internal/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Read and normalize header
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "missing authorization header",
				Code:  "AUTH_HEADER_MISSING",
			})
			c.Abort()
			return
		}

		// 2. Validate Bearer prefix (case-insensitive)
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "invalid authorization format",
				Code:  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		// 3. Extract token safely
		tokenStr := strings.TrimSpace(authHeader[7:]) // remove "Bearer "
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "empty token",
				Code:  "EMPTY_TOKEN",
			})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")

		// 4. Parse and validate token
		parsedToken, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
			// enforce HS256
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "invalid token",
				Code:  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 5. Extract claims
		claims, ok := parsedToken.Claims.(*model.Claims)
		if !ok || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "invalid token claims",
				Code:  "INVALID_TOKEN_CLAIMS",
			})
			c.Abort()
			return
		}

		// 6. Store in context
		c.Set("claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}