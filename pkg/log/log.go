package log

import (
	"context"
	"io"
	"log"
	"os"
	"time"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	}
	return ""
}

type Formatter interface {
	Write(io.Writer, Level, time.Time, map[string]interface{}, string)
}

var globalFormatter Formatter

func SetGlobalFormatter(f Formatter) {
	globalFormatter = f
}

type ContextLog struct {
	context.Context
	fmt  Formatter
	sink io.Writer
}

type contextKey string

var logValues contextKey = "fragma_log"

func CtxFields(ctx context.Context) map[string]interface{} {
	value := ctx.Value(logValues)
	if value == nil {
		return map[string]interface{}{}
	}
	fields, ok := value.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return fields
}

func With(ctx context.Context, values ...interface{}) ContextLog {
	fields := CtxFields(ctx)
	for i := 0; i < len(values)-1; i += 2 {
		fields[values[i].(string)] = values[i+1]
	}

	return ContextLog{
		Context: context.WithValue(ctx, logValues, fields),
		fmt:     globalFormatter,
		sink:    os.Stderr,
	}
}

func (c ContextLog) Info(s string) {
	if c.fmt == nil {
		log.Println(s)
		return
	}
	c.fmt.Write(c.sink, LevelInfo, time.Now(), CtxFields(c), s)
}

func (c ContextLog) Warn(s string) {
	if c.fmt == nil {
		log.Println(s)
		return
	}
	c.fmt.Write(c.sink, LevelWarning, time.Now(), CtxFields(c), s)
}

func (c ContextLog) Error(s string) {
	if c.fmt == nil {
		log.Println(s)
		return
	}
	c.fmt.Write(c.sink, LevelError, time.Now(), CtxFields(c), s)
}
