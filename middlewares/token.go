package middlewares

import (
	"gin-api/models"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Api-Token")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error":   1,
				"message": "API token is required",
			})
			return
		}

		user := models.UserValidateToken(token)
		if user.ID > 0 {
			c.Set("user", user)
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"error":   1,
				"message": "Invalid API token",
			})
			return
		}

		c.Next()
	}
}
