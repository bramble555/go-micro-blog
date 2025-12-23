package controller

import (
	"net/http"

	"go-micro-blog/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
	// 1. 定义接收参数的结构体 (DTO)
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Summary string `json:"summary"`
	}

	// 2. 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 3. 调用 Logic 层
	article, err := service.CreateArticle(req.Title, req.Content, req.Summary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": article, "msg": "创建成功"})
}
