package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/models"
	"github.com/torbenconto/spooler/util"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			log.Printf("error getting token from cookie: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := util.ParseJWT(token)
		if err != nil {
			log.Printf("error parsing jwt: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func RoleAuthMiddleware(requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := userData.(*util.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user claims type"})
			return
		}

		if models.Role(claims.Role).Permissions() < requiredRole.Permissions() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you don't have the required permissions to access this resource"})
			return
		}

		c.Next()
	}
}
