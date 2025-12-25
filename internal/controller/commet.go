package controller

import (
	"net/http"
	"strconv"

	"go-micro-blog/internal/service"

	"github.com/gin-gonic/gin"
)

// 删除评论（admin）
func DeleteComment(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := service.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "delete failed"})
		return
	}

	c.Status(http.StatusOK)
}
