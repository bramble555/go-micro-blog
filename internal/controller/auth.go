package controller

import (
	"fmt"
	"go-micro-blog/global"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid params"})
		return
	}
	fmt.Printf("Input: u:%s, p:%s\n", req.Username, req.Password)
	fmt.Printf("Config: u:%s, p:%s\n", global.Config.Admin.Username, global.Config.Admin.Password)
	// ✅ 从配置中校验管理员账号
	if req.Username != global.Config.Admin.Username ||
		req.Password != global.Config.Admin.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid credentials"})
		return
	}

	// ✅ 生成 JWT（roles = ["admin"]）
	claims := jwt.MapClaims{
		"username": req.Username,
		"roles":    []string{"admin"},
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(global.Config.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200, // 🚨 必须加上这个，前端才能判断成功
		"msg":  "登录成功",
		"data": gin.H{
			"token": tokenStr,
		},
	})
}

// RenderLogin 渲染登录页面
func RenderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{
		"Title": "管理员登录",
	})
}

// Me 返回当前用户信息（需要 Authorization header）
func Me(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing auth"})
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid auth header"})
		return
	}
	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid claims"})
		return
	}

	// 提取 username
	username := ""
	if u, ok := claims["username"]; ok {
		username = fmt.Sprintf("%v", u)
	}

	// 提取 roles
	roles := []string{}
	if r, ok := claims["roles"]; ok {
		switch v := r.(type) {
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					roles = append(roles, s)
				}
			}
		case []string:
			roles = v
		case string:
			roles = []string{v}
		}
	} else if r, ok := claims["role"]; ok {
		if s, ok := r.(string); ok {
			roles = []string{s}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"username": username, "roles": roles}})
}

// RenderCreateArticle 页面渲染
func RenderCreateArticle(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/create_article.html", nil)
}
