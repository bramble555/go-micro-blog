package logger

import (
	"go-micro-blog/global"
	"go-micro-blog/internal/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化 Zap 日志
func InitLogger(cfg config.Logger) {
	// 1. 配置日志切割 (Lumberjack)
	// 这是一个 WriteSyncer，用于文件输出
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,   // 日志文件路径
		MaxSize:    cfg.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: cfg.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     cfg.MaxAge,     // 文件最多保存多少天
		Compress:   true,           // 是否压缩
	})

	// 2. 设置日志级别
	// 从配置中读取 level 字符串 (如 "debug", "info") 并解析为 zapcore.Level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel // 解析失败默认使用 Info
	}

	// 3. 配置编码器 (Encoder)
	// 我们使用自定义的配置，让输出更美观
	encoderConfig := zap.NewProductionEncoderConfig()

	// time 格式：2006-01-02 15:04:05
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// level 格式：大写 + 彩色 (DEBUG, INFO...)
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// caller 格式：简短路径 (main.go:15)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// duration 格式：秒数
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	// 4. 创建核心 (Core)
	// 这是一个“多路复用”的核心：
	// - 控制台输出：使用 Console 编码（方便人类阅读）
	// - 文件输出：使用 JSON 编码
	// 这里为了简单和统一，我们全部使用 ConsoleEncoder (文本格式)，
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder, // 编码器
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileWriteSyncer), // 输出位置：控制台 + 文件
		level, // 日志级别
	)

	// 5. 构造 Logger
	// AddCaller: 添加调用者信息 (文件名:行号)
	logger := zap.New(core, zap.AddCaller())

	// 6. 赋值给全局变量
	global.Log = logger
}
