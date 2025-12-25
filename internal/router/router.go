package router

import (
	"go-micro-blog/internal/controller"
	"go-micro-blog/internal/controller/front"
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
	// 基础测试
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// 1. 首页渲染：只返回 HTML 结构
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"Title": "首页",
		})
	})
	// 🔑 登录页面 (GET)：显示 HTML 界面
	r.GET("/admin/login", controller.RenderLogin)
	// 🚀 在加载静态文件和模板的代码附近增加

	// 📡 公开 API 分组
	// --- 页面路由 (用于返回 HTML 壳子) ---
	apiPublic := r.Group("/api")
	{
		// 登录动作 (POST)：接收 JSON 账号密码，签发 Token
		apiPublic.POST("/login", controller.Login)

		// 获取文章列表/详情 (所有人可见)
		apiPublic.GET("/articles", controller.GetArticleList)
		apiPublic.GET("/articles/:id", controller.GetArticleDetail)
		// 🚀 修改这里：去掉多余层级，保持简单
		apiPublic.GET("/comments", front.GetComments)    // 获取评论列表
		apiPublic.POST("/comments", front.CreateComment) // 提交评论
	}

	// ==========================================
	// 2. 私密区域：仅管理员可见 (受 JWT 保护)
	// ==========================================

	// 🔴 管理员页面渲染分组
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth()) // 挂载严格校验中间件
	{
		// 只有带 Token 的管理员才能看写文章页面
		admin.GET("/create", controller.RenderCreateArticle)
		admin.POST("/comments/:id/delete", controller.DeleteComment)
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
