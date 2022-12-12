package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/senpan/xlogger/core"
	"github.com/senpan/xlogger/logtrace"
)

type LogrusBuilder struct {
	logger  *logrus.Logger
	version string
}

func NewLogrusBuilder(logC *core.XLoggerConf, version string) *LogrusBuilder {
	logger := initLogrus(logC)
	return &LogrusBuilder{
		logger:  logger,
		version: version,
	}
}

func (lb *LogrusBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
	if tag == "" {
		tag = "NoTag"
	}
	tag = core.Filter(tag)
	// 如果需要将ctx中信息，写入日志中，请在这里额外处理
	loggerFields := make(map[string]interface{})
	loggerFields["tag"] = tag
	loggerFields["caller"] = lb.getCaller()

	if ctx != nil {
		traceNode := logtrace.ExtractTraceNodeFromContext(ctx)
		metadata := traceNode.ForkMap()
		for key, val := range metadata {
			loggerFields[key] = val
		}
		logtrace.IncrementRpcId(ctx)
		// 计算使用时间
		if startValue := ctx.Value("__svc_start__"); startValue != nil {
			if start, ok := startValue.(time.Time); ok {
				cost := time.Since(start)
				loggerFields["cost"] = fmt.Sprintf("%.2f", cost.Seconds()*1e3)
			}
		}
	}

	field := lb.logger.WithFields(loggerFields)
	switch lvl {
	case "DEBUG":
		field.Debugf(cast.ToString(args), v...)
	case "TRACE":
		field.Tracef(cast.ToString(args), v...)
	case "INFO":
		field.Infof(cast.ToString(args), v...)
	case "WARNING":
		field.Warnf(cast.ToString(args), v...)
	case "ERROR":
		field.Errorf(cast.ToString(args), v...)
	case "FATAL":
		field.Panicf(cast.ToString(args), v...)
	case "PANIC":
		field.Panicf(cast.ToString(args), v...)
	}
}

func (lb *LogrusBuilder) Close() {
	lb.logger.Exit(0)
}

func (lb *LogrusBuilder) SetVersion(version string) {
	lb.version = version
}

func (lb *LogrusBuilder) getCaller() string {
	caller := ""
	_, file, line, ok := runtime.Caller(3)
	if ok {
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file + ":" + strconv.Itoa(line)
		}
		// Find the penultimate separator.
		idx = strings.LastIndexByte(file[:idx], '/')
		if idx == -1 {
			return file + ":" + strconv.Itoa(line)
		}
		b := strings.Builder{}
		b.WriteString(file[idx+1:])
		b.WriteString(":")
		b.WriteString(strconv.Itoa(line))
		return b.String()
	}

	return caller
}

func initLogrus(conf *core.XLoggerConf) *logrus.Logger {
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
	// 初始化全局logger实例
	l := logrus.New()
	// 添加logger配置项
	logLevel, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "InitLoggerConfig: ParseLevel Error,Set Level to Info;err:%+v", err)
		logLevel = logrus.DebugLevel
	}
	// 设置日志等级
	l.SetLevel(logLevel)
	// 设置格式化方法
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	l.SetOutput(writer)
	return l
}
