package middleware

import (
	"net/http"
	"strings"

	"go-micro-blog/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth 解析 JWT，把 role 放入 gin.Context
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			// visitor，直接放行（不 abort）
			c.Next()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid auth header"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.Config.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid claims"})
			c.Abort()
			return
		}

		// ✅ 放入 Context
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
