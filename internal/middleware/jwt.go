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
		if u, ok := claims["username"]; ok {
			c.Set("username", u)
		}

		// 支持新格式 roles: []string 以及老格式 role: "admin"
		if r, ok := claims["roles"]; ok {
			switch v := r.(type) {
			case []interface{}:
				roles := make([]string, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						roles = append(roles, s)
					}
				}
				c.Set("roles", roles)
				for _, role := range roles {
					if role == "admin" {
						c.Set("role", "admin")
						break
					}
				}
			case []string:
				c.Set("roles", v)
				for _, role := range v {
					if role == "admin" {
						c.Set("role", "admin")
						break
					}
				}
			case string:
				// 兼容老字段
				c.Set("roles", []string{v})
				c.Set("role", v)
			}
		} else if r, ok := claims["role"]; ok {
			if s, ok := r.(string); ok {
				c.Set("role", s)
				c.Set("roles", []string{s})
			}
		}

		c.Next()
	}
}

// RequireAdmin 强制要求 Authorization 且必须包含 admin 角色
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "missing auth"})
			c.Abort()
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

		// 提取角色并检查 admin
		isAdmin := false
		if r, ok := claims["roles"]; ok {
			switch v := r.(type) {
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok && s == "admin" {
						isAdmin = true
						break
					}
				}
			case []string:
				for _, s := range v {
					if s == "admin" {
						isAdmin = true
						break
					}
				}
			case string:
				if v == "admin" {
					isAdmin = true
				}
			}
		} else if r, ok := claims["role"]; ok {
			if s, ok := r.(string); ok && s == "admin" {
				isAdmin = true
			}
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"msg": "need admin role"})
			c.Abort()
			return
		}

		// 放入 context 供后续处理器使用
		if u, ok := claims["username"]; ok {
			c.Set("username", u)
		}
		c.Set("roles", []string{"admin"})

		c.Next()
	}
}
