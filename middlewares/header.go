package middlewares

import (
	"github.com/gin-gonic/gin"
)

func HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			c.Header("Cache-Control", "no-cache")
			c.Header("X-Frame-Options", "deny")
			c.Header("X-Content-Type-Options", "nosniff")
			c.Header("Content-Security-Policy", "default-src 'self' data:; script-src 'unsafe-inline' 'unsafe-eval' 'self' https:; style-src 'unsafe-inline' https:; connect-src 'self'; img-src https: data:; child-src https:; media-src 'none'; object-src 'none';")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next()
	}
}
