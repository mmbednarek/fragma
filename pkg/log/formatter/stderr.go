package formatter

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mmbednarek/fragma/pkg/log"
)

const OpaqueFmt = "\033[0m"
const GreenFmt = "\033[0;32m"
const RedFmt = "\033[0;31m"
const YellowFmt = "\033[0;33m"
const UnderlineFmt = "\033[1m"

func levelToFmt(level log.Level) string {
	switch level {
	case log.LevelInfo:
		return GreenFmt
	case log.LevelWarning:
		return YellowFmt
	case log.LevelError:
		return RedFmt
	}
	return OpaqueFmt
}

type BasicFormatter struct {
}

func writeObject(b *strings.Builder, obj interface{}) {
	if s, ok := obj.(string); ok {
		b.WriteString(RedFmt + "\"")
		b.WriteString(s)
		b.WriteString("\"" + OpaqueFmt)
		return
	}
	if err, ok := obj.(error); ok {
		b.WriteString(RedFmt + "\"")
		b.WriteString(err.Error())
		b.WriteString("\"" + OpaqueFmt)
		return
	}
	if num, ok := obj.(int); ok {
		b.WriteString(YellowFmt)
		b.WriteString(strconv.Itoa(num))
		b.WriteString(OpaqueFmt)
		return
	}
	b.WriteString(fmt.Sprintf("%v", obj))
}

func (s BasicFormatter) Write(sink io.Writer, level log.Level, time time.Time, fields map[string]interface{}, msg string) {
	levelStr := level.String()

	b := strings.Builder{}
	b.WriteString("[")
	b.WriteString(time.Format("2006-01-02 15:04:05"))
	b.WriteString(OpaqueFmt)
	b.WriteString("] ")
	b.WriteString("[")
	b.WriteString(levelToFmt(level))
	b.WriteString(levelStr)
	b.WriteString(OpaqueFmt)
	b.WriteString("]")
	b.WriteString(strings.Repeat(" ", 6-len(levelStr)))
	b.WriteString(msg)

	for key, value := range fields {
		b.WriteString(" " + UnderlineFmt)
		b.WriteString(key)
		b.WriteString(OpaqueFmt + "=")
		writeObject(&b, value)
	}
	b.WriteString("\n")

	if _, err := sink.Write([]byte(b.String())); err != nil {
		panic(err)
	}
}

func init() {
	log.SetGlobalFormatter(BasicFormatter{})
}
