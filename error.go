package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	"github.com/senpan/xlogger/core"
	"github.com/senpan/xlogger/stackerr"
)

type XError struct {
	Code int
	Msg  string
}

const (
	SystemDefaultErrorCode = 50000
)

// NewError 构造错误
// err 如果err的类型是err或string,将错误信息写入ErrorMessage
// 如果err是StackErr,直接返回
// ext ext[0]:错误XError
// ext ext[0]:错误code  ext[1]:返回给调用端的错误信息
func NewError(err interface{}, ext ...XError) *stackerr.StackErr {
	return newError(err, ext...)
}

// NewXError 自定义构成错误
// 按照用户给定错误码和错误信息构造
func NewXError(code int, msg string) error {
	return XError{Code: code, Msg: msg}
}

// NewDiyError 基于错误构建自定义错误
// 按照用户给定错误码和错误对象构造
func NewDiyError(code int, err error) error {
	msg := err.Error()
	if errCode, ok := transfer(err); ok {
		msg = errCode.Msg
	}
	return fmt.Errorf("%d|%s", code, msg)
}

func newError(err interface{}, ext ...XError) *stackerr.StackErr {
	var errInfo string
	switch t := err.(type) {
	case *stackerr.StackErr:
		return t
	case string:
		errInfo = core.Filter(t)
	case error:
		errInfo = core.Filter(t.Error())
		if lhCode, ok := transfer(t); ok {
			ext = make([]XError, 1)
			ext[0] = lhCode
		}
	default:
		errInfo = core.Filter(fmt.Sprintf("%v", t))
	}
	stackErr := &stackerr.StackErr{}

	stackErr.Info = errInfo
	_, file, line, ok := runtime.Caller(2)
	if ok {
		stackErr.Line = line
		components := strings.Split(file, "/")
		stackErr.Filename = components[(len(components) - 1)]
		stackErr.Position = filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)
	stackErr.StackTrace = core.Filter(string(buf[:n]), " ")

	if len(ext) >= 1 {
		c := ext[0]
		stackErr.Code = c.Code
		stackErr.Message = c.Msg
	} else {
		stackErr.Code = SystemDefaultErrorCode
		stackErr.Message = errInfo
	}

	return stackErr
}

func transfer(err error) (le XError, ok bool) {
	segments := strings.SplitN(err.Error(), "|", 2)
	code := cast.ToInt(segments[0])
	if len(segments) > 1 && code > 0 {
		le.Code = code
		le.Msg = segments[1]
		ok = true
	}
	return
}

func (le XError) Error() string {
	return fmt.Sprintf("%d|%s", le.Code, le.Msg)
}
