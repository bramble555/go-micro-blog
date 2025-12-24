package router

import (
	"go-micro-blog/internal/controller"
	"go-micro-blog/internal/middleware"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, mode string, wg *sync.WaitGroup) *gin.Engine {
	// 如果是发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}

	// ==========================================
	// 1. 公开区域：所有人可见 (无需任何中间件)
	// ==========================================

	// 基础测试
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// 首页渲染
	r.GET("/", func(c *gin.Context) {
		// 临时数据，后续从 Service 获取
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

		c.HTML(http.StatusOK, "base.html", gin.H{
			"Title":   "首页",
			"Posts":   posts,
			"IsAdmin": false, // 后续可以通过 cookie 或 session 判断
		})
	})

	// 🔑 登录页面 (GET)：显示 HTML 界面
	// 🚨 注意：这里一定要用 RenderLogin
	r.GET("/admin/login", controller.RenderLogin)

	// 📡 公开 API 分组
	apiPublic := r.Group("/api")
	{
		// 登录动作 (POST)：接收 JSON 账号密码，签发 Token
		apiPublic.POST("/login", controller.Login)

		// 获取文章列表/详情 (所有人可见)
		// apiPublic.GET("/articles", controller.GetArticleList)
		// apiPublic.GET("/articles/:id", controller.GetArticleDetail)
	}

	// ==========================================
	// 2. 私密区域：仅管理员可见 (受 JWT 保护)
	// ==========================================

	// 🔴 管理员页面渲染分组
	adminPage := r.Group("/admin")
	adminPage.Use(middleware.JWTAuth()) // 挂载严格校验中间件
	{
		// 只有带 Token 的管理员才能看写文章页面
		adminPage.GET("/create", controller.RenderCreateArticle)
	}

	// 🔴 管理员操作 API 分组
	apiAdmin := r.Group("/api")
	apiAdmin.Use(middleware.JWTAuth())
	{
		// 只有带 Token 的管理员才能通过接口发文章
		InitArticleRoutes(apiAdmin)

		// 以后可以加更多
		// InitMessageRoutes(apiAdmin)
	}

	return r
}
