package core

import (
	"context"
)

var builder MessageBuilder

type MessageBuilder interface {
	LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{})
	Close()
	SetVersion(version string)
}

func SetBuilder(b MessageBuilder) {
	builder = b
}

func GetBuilder() MessageBuilder {
	return builder
}
