package core

import (
	"sync"

	"github.com/spf13/cast"

	"github.com/senpan/xtools/confx"
)

type XLoggerConf struct {
	Lib        string `ini:"lib" yaml:"lib"`
	Mode       string `ini:"mode" yaml:"mode"`
	Level      string `ini:"level" yaml:"level"`
	Filename   string `ini:"filename" yaml:"filename"`
	MaxSize    int    `ini:"maxSize" yaml:"maxSize"`
	MaxBackups int    `ini:"maxBackups" yaml:"maxBackups"`
	MaxAge     int    `ini:"maxAge" yaml:"maxAge"`
	Compress   bool   `ini:"compress" yaml:"compress"`
}

var loggerConf *XLoggerConf
var loggerOnce sync.Once

func GetXLoggerConf() *XLoggerConf {
	if loggerConf != nil {
		return loggerConf
	}
	// 初始化一次
	loggerOnce.Do(func() {
		logMap := make(map[string]string)
		logMap = confx.GetConfToMap("Logger")
		loggerConf = new(XLoggerConf)
		loggerConf.Mode = "stdout"
		loggerConf.Level = "DEBUG"
		for k, v := range logMap {
			switch k {
			case "Lib":
				loggerConf.Lib = v
			case "mode":
				loggerConf.Mode = v
			case "level":
				loggerConf.Level = v
			case "filename":
				loggerConf.Filename = v
			case "maxSize":
				loggerConf.MaxSize = cast.ToInt(v)
			case "maxBackups":
				loggerConf.MaxBackups = cast.ToInt(v)
			case "maxAge":
				loggerConf.MaxAge = cast.ToInt(v)
			case "compress":
				loggerConf.Compress = v == "true"
			}
		}
	})

	return loggerConf
}

func GetDefaultConf() *XLoggerConf {
	return &XLoggerConf{
		Lib:        "zap",
		Mode:       "stdout",
		Level:      "DEBUG",
		Filename:   "./logs/xlogger.log",
		MaxSize:    128,
		MaxBackups: 2,
		MaxAge:     3,
		Compress:   false,
	}
}
