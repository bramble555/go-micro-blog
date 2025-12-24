package config

// Config 聚合所有配置
type Config struct {
	Server Server `mapstructure:"server"`
	MySQL  MySQL  `mapstructure:"mysql"`
	Logger Logger `mapstructure:"logger"`
	Admin  Admin  `mapstructure:"admin"`
	JWT    JWT    `mapstructure:"jwt"`
}

type Server struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug 或 release
}

type MySQL struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Logger struct {
	Level      string `mapstructure:"level"`       // debug, info, error
	Filename   string `mapstructure:"filename"`    // 日志保存路径
	MaxSize    int    `mapstructure:"max_size"`    // 每个日志文件最大 MB
	MaxAge     int    `mapstructure:"max_age"`     // 保留多少天
	MaxBackups int    `mapstructure:"max_backups"` // 保留多少个文件
}

// Admin 初始管理员配置
type Admin struct {
	Username string `mapstructure:"username"` // 管理员用户名
	Password string `mapstructure:"password"` // 管理员密码
}

type JWT struct {
	Secret string `mapstructure:"secret"` // JWT 密钥
}
