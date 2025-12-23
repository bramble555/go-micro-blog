package router

import (
	"go-micro-blog/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitArticleRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.POST("articles", controller.CreateArticle) // 提交文章数据
	// r.GET("/articles", controller.CreateArticle)
	// r.GET("/articles/:id", controller.GetArticlesDetailHandler)
	// r.PUT("articles/:id", controller.UpdateArticlesHandler)
	// r.POST("articles/digg", controller.PostArticleDigHandler)
	return r
}
