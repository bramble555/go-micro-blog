package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderCreateArticle 页面渲染
func RenderCreateArticle(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/create_article.html", nil)
}
