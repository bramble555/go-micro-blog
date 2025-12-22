package main

import (
	"html/template"
	"io/fs"
	"log"
	"path/filepath"

	"go-micro-blog/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin 默认引擎（包含 Logger + Recovery）
	r := gin.Default()

	// 加载 HTML 模板（递归 templates 目录下所有 .html 文件）
	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatal("failed to load templates:", err)
	}
	r.SetHTMLTemplate(tmpl)

	// 静态资源
	r.Static("/assets", "./assets")

	// 注册路由
	router.Register(r)

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatal("server start failed:", err)
	}
}

// loadTemplates 会递归遍历模板目录，加载所有 .html 文件
func loadTemplates(root string) (*template.Template, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fs.ErrNotExist
	}
	// template.ParseFiles 会使用文件名作为模板名
	tmpl := template.Must(template.ParseFiles(files...))
	log.Println("[INFO] Loaded HTML templates:", files)
	return tmpl, nil
}
