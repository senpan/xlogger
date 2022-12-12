package stackerr

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

type StackErr struct {
	Filename   string
	Line       int
	Message    string // 标准输出报错信息
	StackTrace string
	Code       int    // 错误码
	Info       string // 错误详情
	Position   string
	Level      int // 0最高优先级 1-4 普通优先级 5 可不关注的异常
}

func (se *StackErr) ErrorInfo() string {
	return se.Info
}

func (se *StackErr) Error() string {
	return fmt.Sprintf("%d|%s", se.Code, se.Message)
}

func (se *StackErr) Stack() string {
	return fmt.Sprintf("(%s:%d)%s\tStack: %s", se.Filename, se.Line, se.Info, se.StackTrace)
}

func (se *StackErr) Detail() string {
	return fmt.Sprintf("(%s:%d)%s", se.Filename, se.Line, se.Info)
}

func (se *StackErr) Format(tag ...string) (data string) {
	var strVal []string
	strVal = append(strVal, cast.ToString(se.Code))
	strVal = append(strVal, se.Message)
	strVal = append(strVal, se.Filename)
	strVal = append(strVal, cast.ToString(se.Line))
	strVal = append(strVal, se.Info)
	data = strings.Join(strVal, "\t")
	return
}

func (se *StackErr) SetLevel(lvl int) {
	se.Level = lvl
}

func (se *StackErr) GetLevel() int {
	return se.Level
}
