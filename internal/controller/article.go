package controller

import (
	"go-micro-blog/internal/service"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

// CreateArticle 创建文章
func CreateArticle(c *gin.Context) {
	// 1. 定义接收参数的结构体 (DTO)
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	// 2. 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, CodeInvalidParam)
		return
	}

	// 3. 调用 Logic 层
	article, err := service.CreateArticle(req.Title, req.Content)
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}

	Success(c, article)
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	// 调用 Service 获取真实数据
	articles, err := service.GetArticleList()
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}
	Success(c, articles)
}

// GetArticleDetailAPI 处理获取单篇文章数据的 JSON API 请求
func GetArticleDetailAPI(c *gin.Context) {
	id := c.Param("id")
	article, err := service.GetArticleByID(id)
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}

	// 获取评论列表
	comments, _ := service.GetCommentsByArticleID(article.ID)

	// Markdown 转 HTML
	htmlContent := blackfriday.Run([]byte(article.Content))

	Success(c, gin.H{
		"article":      article,
		"html_content": string(htmlContent),
		"comments":     comments,
	})
}

// GetArticleDetail 处理渲染文章详情页面的请求
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	article, err := service.GetArticleByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "文章不存在")
		return
	}

	// 获取评论列表
	comments, _ := service.GetCommentsByArticleID(article.ID)

	// Markdown 转 HTML
	htmlContent := blackfriday.Run([]byte(article.Content))
	article.Content = string(htmlContent)

	// 渲染 post.html 模板
	c.HTML(http.StatusOK, "post.html", gin.H{
		"Article":     article,
		"HTMLContent": template.HTML(article.Content),
		"Comments":    comments,
	})
}
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteArticle(id); err != nil {
		Fail(c, CodeServerBusy)
		return
	}
	Success(c, nil)
}
