package logger

import (
	"proposal-template/pkg/logger/internal"
)
type ILogger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    GetLevel() string
}

func NewLogger(level string) ILogger {
    return internal.NewZapLogger(level)
}