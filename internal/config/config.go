package config

// Config 聚合所有配置
type Config struct {
	Server Server `mapstructure:"server"`
	MySQL  MySQL  `mapstructure:"mysql"`
	Logger Logger `mapstructure:"logger"`
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
