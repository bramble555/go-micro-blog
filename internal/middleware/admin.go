package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly 仅允许 admin 访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "admin only",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
