package controller

import (
	"go-micro-blog/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLikes(c *gin.Context) {
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
	count, err := service.GetLikesByArticleID(articleID)
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}
	Success(c, count)
}
func CreateLike(c *gin.Context) {
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
	err = service.CreateLike(articleID)
	if err != nil {
		Fail(c, CodeServerBusy)
		return
	}
	Success(c, nil)
}
