package main

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"go-micro-blog/global"
	"go-micro-blog/internal/dao/mysql"
	"go-micro-blog/internal/router"
	"go-micro-blog/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	var err error
	// 1. 初始化配置 (Viper)
	viper.SetConfigFile("configs/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置读取失败: %s", err))
	}
	// 将配置映射到全局变量 global.Config
	if err = viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("配置解析失败: %s", err))
	}

	// 2. 初始化日志 (Zap)
	logger.InitLogger(global.Config.Logger)
	global.Log.Info("项目启动初始化中...")

	// 3. 初始化 Gin
	gin.SetMode(global.Config.Server.Mode)

	// 初始化数据库
	global.DB, err = mysql.Init()
	if err != nil {
		global.Log.Fatal("数据库初始化失败", zap.Error(err))
	}
	// 初始化雪花算法
	global.InitSnowflake(1)
	startServer()
}

// 启动 Gin 服务器
func startServer() {
	r := gin.Default()
	// 加载 HTML 模板（递归 templates 目录下所有 .html 文件）
	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatal("failed to load templates:", err)
	}
	r.SetHTMLTemplate(tmpl)

	// 静态资源
	r.Static("/assets", "./assets")
	wg := sync.WaitGroup{}
	router.InitRouter(r, global.Config.Server.Mode, &wg)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", global.Config.Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutdown Server ...")

	// 等待所有后台任务完成
	wg.Wait()

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 优雅关闭服务
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("服务器监听失败", zap.Error(err))
	}
	global.Log.Info("Server exiting")
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
