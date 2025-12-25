package front

import (
	"go-micro-blog/internal/service"
	"net/http" // 确保引入了模型
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 提交评论（适配 Fetch API）
func CreateComment(c *gin.Context) {
	// 1. 定义接收前端 JSON 的结构体
	var req struct {
		ArticleID int64  `json:"article_id,string" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Nickname  string `json:"nickname"`
	}

	// 2. 解析 JSON Body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误: " + err.Error()})
		return
	}

	if req.Nickname == "" {
		req.Nickname = "匿名用户"
	}

	// 3. 调用 Service 保存到数据库
	// 注意：确保你的 service.CreateComment 接收的是 (int64, string, string)
	err := service.CreateComment(req.ArticleID, req.Nickname, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "发布失败"})
		return
	}

	// 4. 返回 JSON 而不是重定向
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发布成功",
	})
}

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	// 1. 获取 URL 里的字符串参数: /api/comments?article_id=xxxx
	articleIDStr := c.Query("article_id")
	if articleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "缺少 article_id"})
		return
	}

	// 2. 🚀 关键步骤：将 string 转为 int64
	// 10 表示十进制，64 表示 int64 类型
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID 格式错误"})
		return
	}

	// 3. 现在可以传给接收 int64 的 Service 函数了
	comments, err := service.GetCommentsByArticleID(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加载失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": comments,
	})
}
