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
	// 1. 初始化配置 (Viper)
	viper.SetConfigFile("configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置读取失败: %s", err))
	}
	// 将配置映射到全局变量 global.Config
	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("配置解析失败: %s", err))
	}

	// 2. 初始化日志 (Zap)
	logger.InitLogger(global.Config.Logger)
	global.Log.Info("项目启动初始化中...")

	// 3. 初始化 Gin
	gin.SetMode(global.Config.Server.Mode)
	r := gin.Default()

	// 加载 HTML 模板（递归 templates 目录下所有 .html 文件）
	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatal("failed to load templates:", err)
	}
	r.SetHTMLTemplate(tmpl)

	// 静态资源
	r.Static("/assets", "./assets")

	// 初始化数据库
	global.DB, err = mysql.Init()
	if err != nil {
		global.Log.Fatal("数据库初始化失败", zap.Error(err))
	}
	// 初始化雪花算法
	global.InitSnowflake(1)

	// 注册路由
	router.Register(r)

	// 4.想要优雅关机
	srv := &http.Server{
		Addr:    ":" + global.Config.Server.Port,
		Handler: r,
	}

	// 5. 在另一个协程里启动服务，否则它会阻塞主线程，导致后面监听信号的代码跑不到
	go func() {
		global.Log.Info("服务器启动成功", zap.String("port", global.Config.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Fatal("监听失败", zap.Error(err))
		}
	}()

	// 6. 优雅关机逻辑：等待中断信号
	quit := make(chan os.Signal, 1)
	// 监听 Ctrl+C (SIGINT) 和 杀进程 (SIGTERM)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在这里，收到信号才往下走
	global.Log.Info("正在关闭服务器...")

	// 7. 给 5 秒钟时间处理剩下的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("服务器强制关闭", zap.Error(err))
	}

	global.Log.Info("服务器已优雅关机")
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
