package router

import (
	"go-micro-blog/internal/controller"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, mode string, wg *sync.WaitGroup) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
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

		// 注意：这里渲染的是定义了整个 HTML 结构的那个模板名
		c.HTML(http.StatusOK, "base.html", gin.H{
			"Title": "首页",
			"Posts": posts,
		})
	})
	r.GET("/admin/create", controller.RenderCreateArticle) // 渲染创建文章页面
	apiGroup := r.Group("/api")
	InitArticleRoutes(apiGroup)
	// InitMessageRoutes(apiGroup)
	// InitCommentRoutes(apiGroup)
	return r
}

