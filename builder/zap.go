package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/senpan/xlogger/core"
	"github.com/senpan/xlogger/logtrace"
	"github.com/senpan/xlogger/stackerr"
)

type ZapBuilder struct {
	logger  *zap.Logger
	version string
}

func NewZapBuilder(logC *core.XLoggerConf, version string) *ZapBuilder {
	logger := initZap(logC)
	return &ZapBuilder{
		logger:  logger,
		version: version,
	}
}

func (zb *ZapBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
	if tag == "" {
		tag = "NoTag"
	}
	tag = core.Filter(tag)
	_, message := zb.Build(args, v...)

	s := zb.logger.Sugar()
	length := 6
	var metadata map[string]string
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx != nil {
		traceNode := logtrace.ExtractTraceNodeFromContext(ctx)
		metadata = traceNode.ForkMap()
		length = length + len(metadata)*2
	}
	keyValues := make([]interface{}, 0, length)
	keyValues = append(keyValues, "tag", tag)
	if ctx != nil {
		for key, value := range metadata {
			keyValues = append(keyValues, key, value)
		}
		logtrace.IncrementRpcId(ctx)
		// 计算使用时间
		if startValue := ctx.Value("__svc_start__"); startValue != nil {
			if start, ok := startValue.(time.Time); ok {
				cost := time.Since(start)
				keyValues = append(keyValues, "cost", fmt.Sprintf("%.2f", cost.Seconds()*1e3))
			}
		}
	}

	switch lvl {
	case "DEBUG":
		s.Debugw(message, keyValues...)
	case "TRACE":
		s.Infow(message, keyValues...)
	case "INFO":
		s.Infow(message, keyValues...)
	case "WARNING":
		s.Warnw(message, keyValues...)
	case "ERROR":
		s.Errorw(message, keyValues...)
	case "FATAL":
		s.Panicw(message, keyValues...)
	case "PANIC":
		s.Panicw(message, keyValues...)
	}
}

func (zb *ZapBuilder) Build(args interface{}, v ...interface{}) (position string, message string) {
	switch t := args.(type) {
	case *stackerr.StackErr:
		message = t.Info
	case error:
		message = t.Error()
	case string:
		if len(v) > 0 {
			message = fmt.Sprintf(t, v...)
		} else {
			message = t
		}
	default:
		message = fmt.Sprint(t)
	}
	message = core.Filter(message)
	return
}

func (zb *ZapBuilder) Close() {
	_ = zb.logger.Sync()
}

func (zb *ZapBuilder) SetVersion(version string) {
	zb.version = version
}

func initZap(conf *core.XLoggerConf) *zap.Logger {
	// 获取配置信息
	var writer io.Writer
	if conf.Mode == "stdout" {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:   conf.Filename,
			MaxSize:    conf.MaxSize,    // 日志文件大小，单位是 MB
			MaxBackups: conf.MaxBackups, // 最大过期日志保留个数
			MaxAge:     conf.MaxAge,     // 保留过期文件最大时间，单位 天
			Compress:   conf.Compress,   // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
		}
	}
	// 动态日志等级
	dynamicLevel := zap.NewAtomicLevel()
	// 设置日志等级
	dynamicLevel.SetLevel(parseLoggerLevel(conf.Level))
	zc := zapcore.NewTee(
		// 有好的格式、输出控制台、动态等级
		zapcore.NewCore(zapcore.NewJSONEncoder(NewEncoderConfig()), zapcore.AddSync(writer), dynamicLevel),
	)
	return zap.New(zc, zap.AddCaller(), zap.AddCallerSkip(2))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",                        // json时时间键
		LevelKey:       "level",                       // json时日志等级键
		NameKey:        "file",                        // json时日志记录器名
		CallerKey:      "caller",                      // json时日志文件信息键
		MessageKey:     "msg",                         // json时日志消息键
		StacktraceKey:  "stack",                       // json时堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 友好日志等级名大小写（info INFO）
		EncodeTime:     zapcore.RFC3339TimeEncoder,    // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 日志文件信息（包/文件.go:行号）
	}
}

// 解析日志等级
func parseLoggerLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	case "error":
		return zap.ErrorLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "info", "trace":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	}
	return zap.InfoLevel
}
