package controller

import (
	"go-micro-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	isAdmin := false
	if rolesIf, exists := c.Get("roles"); exists {
		if roles, ok := rolesIf.([]string); ok {
			for _, r := range roles {
				if r == "admin" {
					isAdmin = true
					break
				}
			}
		}
	} else if roleIf, exists := c.Get("role"); exists {
		if roleStr, ok := roleIf.(string); ok && roleStr == "admin" {
			isAdmin = true
		}
	}

	articles, _ := service.GetArticleList()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Articles": articles,
		"IsAdmin":  isAdmin,
	})
}
