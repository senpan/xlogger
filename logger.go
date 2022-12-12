package logger

import (
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
