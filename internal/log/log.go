package log

import (
	"os"
	"path/filepath"

	"github.com/hankmor/wechat-publisher/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化日志
func InitLogger(cfg *config.LogConfig) error {
	// 创建日志目录
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return err
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 配置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置 lumberjack 写入器
	logFile := filepath.Join(cfg.Path, "wechat-publisher.log")
	writer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.MaxSize,    // 单个日志文件最大大小（MB）
		MaxBackups: cfg.MaxBackups, // 最大备份文件数
		MaxAge:     cfg.MaxAge,     // 日志文件最大保存天数
		Compress:   cfg.Compress,   // 是否压缩日志文件
	}

	// 配置输出
	var cores []zapcore.Core

	// 文件输出
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	)
	cores = append(cores, fileCore)

	// 控制台输出（如果配置了）
	if cfg.Console {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建 logger
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Info 打印 info 级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug 打印 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Error 打印 error 级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 打印 fatal 级别日志并退出
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
