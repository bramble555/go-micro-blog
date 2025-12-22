package mysql

import (
	"fmt"
	"go-micro-blog/global"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init 初始化 MySQL 连接 (使用 GORM)
func Init() (*gorm.DB, error) {
	// 1. 从全局配置拼接 DSN
	m := global.Config.MySQL
	// 建议 DSN 包含 charset 和 loc，保证时间格式和中文字符正确
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
	// 2. 配置 GORM 的日志 (将 GORM 的日志对接到你的 Zap 上)
	// 如果是 Debug 模式，打印所有 SQL；如果是 Release，只打印错误
	newLogger := logger.New(
		writer{log: global.Log}, // 自定义一个 writer，见下方代码
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略未找到记录错误
			Colorful:                  true,                   // 彩色打印
		},
	)

	// 3. 打开连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		global.Log.Error("MySQL 连接失败", zap.Error(err))
		return nil, err
	}

	// 4. 配置连接池 (底层依然是 sql.DB)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数 (这些可以以后写进 config.yaml)
	sqlDB.SetMaxIdleConns(10)           // 空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	global.Log.Info("MySQL 数据库连接成功!")
	return db, nil
}

// writer 定义一个简单的结构体来实现 gorm logger 的 Writer 接口
type writer struct {
	log *zap.Logger
}

func (w writer) Printf(format string, args ...interface{}) {
	w.log.Info(fmt.Sprintf(format, args...))
}
