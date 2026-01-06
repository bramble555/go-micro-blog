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

	// 📡 公开 API 分组
	apiPublic := r.Group("/api")
	{
		// 登录动作 (POST)：接收 JSON 账号密码，签发 Token
		apiPublic.POST("/login", controller.Login)

		apiPublic.GET("/articles", controller.GetArticleList)
		apiPublic.GET("/articles/:id", controller.GetArticleDetail)
		// 获取当前用户信息（需要 Authorization header）
		apiPublic.GET("/me", controller.Me)
		apiPublic.GET("/comments", controller.GetComments)
		// 提交评论（可选认证，让后端能识别管理员身份）
		apiPublic.POST("/comments", middleware.JWTAuth(), controller.CreateComment)

	}

	// 2. 管理员(可选)
	r.GET("/admin/create", middleware.JWTAuth(), controller.RenderCreateArticle)

	// 3. 管理员操作
	apiAdminRoutes := r.Group("/api/admin")
	apiAdminRoutes.Use(middleware.RequireAdmin())
	{
		InitArticleRoutes(apiAdminRoutes)
		apiAdminRoutes.POST("/comments/:id/delete", controller.DeleteComment)
	}
	return r
}
