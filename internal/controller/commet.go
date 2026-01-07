package controller

import (
	"fmt"
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
		Fail(c, CodeInvalidParam)
		return
	}
	// 根据是否管理员决定保存的昵称：
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
		Fail(c, CodeServerBusy)
		return
	}

	// 4. 返回统一 JSON 格式
	Success(c, nil)
}

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	// 1. 获取 URL 里的字符串参数: /api/comments?article_id=xxxx
	articleIDStr := c.Query("article_id")
	if articleIDStr == "" {
		Fail(c, CodeInvalidParam)
		return
	}

	// 2. 将 string 转为 int64
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		Fail(c, CodeInvalidParam)
		return
	}

	// 3. 调用 Service
	comments, err := service.GetCommentsByArticleID(articleID)
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}

	Success(c, comments)
}

// DeleteComment 删除评论（admin）
func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, CodeInvalidParam)
		return
	}

	if err := service.DeleteComment(id); err != nil {
		Fail(c, CodeServerBusy)
		return
	}

	Success(c, nil)
}
