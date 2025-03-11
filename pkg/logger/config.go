package logger

type LogOutput string

const (
	LOG_OUTPUT_CONSOLE LogOutput = "console"
	LOG_OUTPUT_FILE    LogOutput = "file"
)

type LogLevel string

const (
	LOG_TRACE   LogLevel = "trace"
	LOG_DEBUG   LogLevel = "debug"
	LOG_INFO    LogLevel = "info"
	LOG_WARNING LogLevel = "warning"
	LOG_ERROR   LogLevel = "error"
	LOG_FATAL   LogLevel = "fatal"
)

type LoggerConfig struct {
	Output     LogOutput `mapstructure:"output"`
	Level      LogLevel  `mapstructure:"level"`
	FilePath   string    `mapstructure:"file_path"`
	MaxSize    int       `mapstructure:"max_size"`
	MaxBackups int       `mapstructure:"max_backups"`
	MaxAge     int       `mapstructure:"max_age"`
	Compress   bool      `mapstructure:"compress"`
}
