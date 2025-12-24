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
func GetArticleList(c *gin.Context) {
	// 调用 Service 获取真实数据
	articles, err := service.GetArticleList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文章列表失败",
		})
		return
	}

	// 返回标准 JSON 格式
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": articles,
	})
}

// RenderArticleDetail 渲染文章详情 HTML 页面
func RenderArticleDetail(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Title": "文章详情",
		// 注意：这里不需要传数据，让前端页面加载后自己 fetch API
	})
}

// GetArticleDetail 处理获取单篇文章数据的 API 请求
func GetArticleDetail(c *gin.Context) {
	// 1. 获取路径参数 :id
	id := c.Param("id")

	// 2. 调用我们之前写好的 service 方法获取数据
	article, err := service.GetArticleByID(id)
	if err != nil {
		// 如果查不到数据，返回 404
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "很抱歉，找不到该文章",
		})
		return
	}

	// 3. 增加阅读量逻辑 (可选)
	// global.DB.Model(article).Update("view_count", article.ViewCount+1)

	// 4. 返回成功 JSON
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": article, // 这里包含 title, summary, content 等所有加了 json 标签的字段
	})
}
