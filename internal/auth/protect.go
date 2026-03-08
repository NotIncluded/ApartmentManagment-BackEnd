package auth

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Protect(signature []byte, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Authorization header required", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		// Split "Bearer token"
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := response.NewAppResponse(
				http.StatusUnauthorized,
				"Invalid authorization format. Expected: Bearer <token>",
				nil,
			)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Token is empty", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return signature, nil
		})

		if err != nil || !token.Valid {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Invalid token", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Invalid token claims", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Role missing in token", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		// If roles specified -> check
		if len(allowedRoles) > 0 && !slices.Contains(allowedRoles, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Permission denied",
			})
			return
		}

		c.Set("role", role)
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}