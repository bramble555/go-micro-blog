package global

import (
	"go-micro-blog/internal/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config // 全局配置对象
	DB     *gorm.DB       // 全局数据库对象
	Log    *zap.Logger    // 全局日志对象
)
