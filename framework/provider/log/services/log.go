package services

import (
	"bytes"
	"devops-http/framework"
	"devops-http/framework/contract"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"runtime"
)

func NewNiceLog(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	config := container.MustMake(contract.ConfigKey).(contract.Config)
	return &NiceLog{c: container, Zap: Zap(config)}, nil
}

// NiceLog 的通用实例
type NiceLog struct {
	// 五个必要参数
	Zap *zap.Logger
	c   framework.Container // 容器
}

func (log *NiceLog) Panic(msg string, fields ...zapcore.Field) {
	log.Zap.Panic(msg, fields...)
}

func (log *NiceLog) GetZap() *zap.Logger {
	return log.Zap
}

func (log *NiceLog) Error(msg string, fields ...zapcore.Field) {
	s := stack(2, 5)
	log.Zap.Error(fmt.Sprintf("[error] %s :\n%s\n%s",
		msg, s, reset), fields...)
}

func (log *NiceLog) Warn(msg string, fields ...zapcore.Field) {
	log.Zap.Warn(msg, fields...)
}

func (log *NiceLog) Info(msg string, fields ...zapcore.Field) {
	log.Zap.Info(msg, fields...)
}

func (log *NiceLog) Debug(msg string, fields ...zapcore.Field) {
	log.Zap.Debug(msg, fields...)
}

func (log *NiceLog) Fatal(msg string, fields ...zapcore.Field) {
	s := stack(2, 5)
	log.Zap.Fatal(fmt.Sprintf("[error] %s :\n%s\n%s",
		msg, s, reset), fields...)
}

func (log *NiceLog) Trace(msg string, fields ...zapcore.Field) {
	log.Zap.Info(msg, fields...)
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = "\033[0m"
)

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int, limit int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; i <= limit; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
