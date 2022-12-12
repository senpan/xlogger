package logger

import (
	"context"
	"sync"

	"github.com/senpan/xlogger/builder"
	"github.com/senpan/xlogger/core"
)

func init() {
	core.SetBuilder(new(builder.ZapBuilder))
}

var ilOnce sync.Once

// InitXLogger 初始化日志组件
func InitXLogger(version string) {
	ilOnce.Do(func() {
		logConf := core.GetXLoggerConf()
		var mb core.MessageBuilder
		switch logConf.Lib {
		case "logrus":
			mb = builder.NewLogrusBuilder(logConf, version)
		default:
			mb = builder.NewZapBuilder(logConf, version)
		}
		core.SetBuilder(mb)
	})
}

// InitXLoggerFor 使用配置文件初始化日志组件
func InitXLoggerFor(version string, logConf *core.XLoggerConf) {
	var mb core.MessageBuilder
	switch logConf.Lib {
	case "logrus":
		mb = builder.NewLogrusBuilder(logConf, version)
	default:
		mb = builder.NewZapBuilder(logConf, version)
	}
	core.SetBuilder(mb)
}

func D(tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(context.TODO(), "DEBUG", tag, args, v...)
}

func Dx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(ctx, "DEBUG", tag, args, v...)
}

func I(tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(context.TODO(), "INFO", tag, args, v...)
}
func Ix(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(ctx, "INFO", tag, args, v...)
}

func W(tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(context.TODO(), "WARNING", tag, args, v...)
}

func Wx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(ctx, "WARNING", tag, args, v...)
}

func E(tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(context.TODO(), "ERROR", tag, args, v...)
}

func Ex(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(ctx, "ERROR", tag, args, v...)
}

func F(tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(context.TODO(), "FATAL", tag, args, v...)
}

func Fx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	core.GetBuilder().LoggerX(ctx, "FATAL", tag, args, v...)
}

func SetVersion(version string) {
	core.GetBuilder().SetVersion(version)
}

func Close() {
	core.GetBuilder().Close()
}
