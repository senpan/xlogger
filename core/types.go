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

func D(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(context.TODO(), "DEBUG", tag, args, v...)
}

func Dx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "DEBUG", tag, args, v...)
}

func I(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(context.TODO(), "INFO", tag, args, v...)
}
func Ix(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "INFO", tag, args, v...)
}

func W(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(context.TODO(), "WARNING", tag, args, v...)
}

func Wx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "WARNING", tag, args, v...)
}

func E(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(context.TODO(), "ERROR", tag, args, v...)
}

func Ex(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "ERROR", tag, args, v...)
}

func F(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(context.TODO(), "FATAL", tag, args, v...)
}

func Fx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "FATAL", tag, args, v...)
}

func SetVersion(version string) {
	builder.SetVersion(version)
}

func Close() {
	builder.Close()
}
