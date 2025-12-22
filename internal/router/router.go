package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 首页
	r.GET("/", func(c *gin.Context) {
		posts := []gin.H{
			{
				"Title":   "Go 微型博客系统设计",
				"Summary": "基于 Gin + Redis + Snowflake 的高性能博客系统。",
				"Date":    "2025-03-01",
			},
			{
				"Title":   "为什么选择 Go 作为后端语言",
				"Summary": "从并发模型到工程实践，聊聊 Go 的优势。",
				"Date":    "2025-03-02",
			},
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home",
			"Posts": posts,
		})
	})
}
