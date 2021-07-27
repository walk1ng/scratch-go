package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type Level int

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (lv Level) string() string {
	switch lv {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func NewLogger(to io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      to,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, properties map[string]interface{}) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]interface{}) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]interface{}) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]interface{}) (int, error) {

	if level < l.minLevel {
		return 0, nil
	}

	_, f, lineno, _ := runtime.Caller(100)

	body := struct {
		Time    string                 `json:"time" yaml:"ts"`
		Level   string                 `json:"level" yaml:"level"`
		Caller  string                 `json:"caller" yaml:"caller"`
		Message string                 `json:"message" yaml:"message"`
		Props   map[string]interface{} `json:"props,omitempty" yaml:"props"`
		Trace   string                 `json:"trace,omitempty" yaml:"trace"`
	}{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level.string(),
		Caller:  fmt.Sprintf("%s:%d", f, lineno),
		Message: message,
		Props:   properties,
	}

	// if above error level
	if level >= LevelError {
		body.Trace = string(debug.Stack())
	}

	// message bytes
	var line []byte

	line, err := json.Marshal(body)
	if err != nil {
		line = []byte(LevelError.string() + ": failed to marshal message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(p []byte) (n int, err error) {
	return l.print(LevelError, string(p), nil)
}
