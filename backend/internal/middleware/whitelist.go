package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/config"
)

func WhitelistEnabledMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Cfg.Features.EmailWhitelistEnabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "whitelist not enabled"})
			return
		}

		c.Next()
	}
}
