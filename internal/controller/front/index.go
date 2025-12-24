package front

import (
	"go-micro-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	role, _ := c.Get("role")

	articles, _ := service.GetArticleList()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Articles": articles,
		"IsAdmin":  role == "admin",
	})
}
