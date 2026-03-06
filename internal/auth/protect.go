package auth

import (
	"slices"
	"fmt"
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Protect(signature []byte, allowedRoles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		if s == "" {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Authorization header required", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		tokenString := s[len("Bearer "):]

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
		if !ok {
			appErr := response.NewAppResponse(http.StatusUnauthorized, "Missing role in token", nil)
			c.AbortWithStatusJSON(appErr.Status, appErr.Response())
			return
		}

		// --- Role Filtering Logic ---
		isAuthorized := slices.Contains(allowedRoles, role)

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Set("role", role)
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
