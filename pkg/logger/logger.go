// Package logger 提供线程安全的分级日志器。
package logger

import (
	"log"
	"os"
	"sync"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

type Logger struct {
	mu    sync.Mutex
	level Level
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func New() *Logger { return NewLevel(LevelInfo) }

func NewLevel(level Level) *Logger {
	flag := log.LstdFlags | log.Lmicroseconds
	return &Logger{
		level: level,
		debug: log.New(os.Stderr, "[DEBUG] ", flag),
		info:  log.New(os.Stderr, "[INFO ] ", flag),
		warn:  log.New(os.Stderr, "[WARN ] ", flag),
		err:   log.New(os.Stderr, "[ERROR] ", flag),
	}
}

// nil receiver 上的方法均为无操作，使得 *Logger 可作为可选依赖：
// 调用方未配置日志器（传 nil）时不会 panic，仅静默丢弃日志。
func (l *Logger) SetLevel(level Level) {
	if l == nil {
		return
	}
	l.mu.Lock(); defer l.mu.Unlock(); l.level = level
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.mu.Lock(); defer l.mu.Unlock()
	if l.level <= LevelDebug { l.debug.Printf(format, args...) }
}
func (l *Logger) Infof(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.mu.Lock(); defer l.mu.Unlock()
	if l.level <= LevelInfo { l.info.Printf(format, args...) }
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.mu.Lock(); defer l.mu.Unlock()
	if l.level <= LevelWarn { l.warn.Printf(format, args...) }
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.mu.Lock(); defer l.mu.Unlock()
	if l.level <= LevelError { l.err.Printf(format, args...) }
}
func (l *Logger) Printf(format string, args ...interface{}) { l.Infof(format, args...) }
