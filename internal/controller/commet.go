package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"go-micro-blog/internal/service"

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
	// 根据是否管理员决定保存的昵称：
	// - 管理员：如果提交了 nickname 则使用它，否则使用 token 的 username（或回退为 "管理员"）
	// - 非管理员：强制使用默认昵称 "游客"
	nickname := "游客"
	isAdmin := false
	if rolesIf, exists := c.Get("roles"); exists {
		if roles, ok := rolesIf.([]string); ok {
			for _, r := range roles {
				if r == "admin" {
					isAdmin = true
					break
				}
			}
		}
	} else if roleIf, exists := c.Get("role"); exists {
		if roleStr, ok := roleIf.(string); ok && roleStr == "admin" {
			isAdmin = true
		}
	}

	if isAdmin {
		if req.Nickname != "" {
			nickname = req.Nickname
		} else if u, ok := c.Get("username"); ok {
			nickname = fmt.Sprintf("%v", u)
		} else {
			nickname = "管理员"
		}
	}

	// 3. 调用 Service 保存到数据库
	if err := service.CreateComment(req.ArticleID, nickname, req.Content); err != nil {
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

// 删除评论（admin）
func DeleteComment(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := service.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "delete failed"})
		return
	}

	c.Status(http.StatusOK)
}
