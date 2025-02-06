package internal
import (
	"fmt"
	"log"
	"runtime"
	

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements ILogger using Zap.
type zapLogger struct {
    logger   *zap.Logger
    logLevel zapcore.Level
}

// NewZapLogger initializes and returns a new Zap logger with a configurable level.
func NewZapLogger(level string) *zapLogger {
    // Map string level to Zap's zapcore.Level
    var logLevel zapcore.Level
    switch level {
    case "debug":
        logLevel = zapcore.DebugLevel
    case "info":
        logLevel = zapcore.InfoLevel
    case "warn":
        logLevel = zapcore.WarnLevel
    case "error":
        logLevel = zapcore.ErrorLevel
    default:
        logLevel = zapcore.InfoLevel // Default to info level
    }
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        MessageKey:     "msg",
        CallerKey:      "caller",
        StacktraceKey:  "stacktrace",
        EncodeTime:     zapcore.ISO8601TimeEncoder,       // Human-readable time format
        EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Color for levels
        EncodeCaller:   zapcore.ShortCallerEncoder,       // Shorten caller file path
        EncodeDuration: zapcore.StringDurationEncoder,    // Human-readable duration
    }

    // Create a log file
    // logFilePath := filepath.Join("logs", fmt.Sprintf("app_%s.log", time.Now().Format("20060102_150405")))
    // os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm) // Ensure log directory exists

    // Create a file writer for the log file
    // _, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    // if err != nil {
    //     log.Fatalf("Failed to open log file: %v", err)
    // }

    cfg := zap.Config{
        Encoding:         "console",                      // Switch to console encoding
        Level:            zap.NewAtomicLevelAt(logLevel),
        OutputPaths:      []string{"stdout"},             // Log to console (stdout)
        ErrorOutputPaths: []string{"stderr"},             // Error logs to stderr
        // OutputPaths:      []string{"stdout", logFilePath}, // uncomment this to save log file
        // ErrorOutputPaths: []string{"stderr", logFilePath},  // and this
        EncoderConfig:    encoderConfig,                  // Use the customized encoder config
    }

    logger, err := cfg.Build()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }

    return &zapLogger{logger: logger, logLevel: logLevel}
}

// Helper to convert variadic fields to Zap fields
func toZapFields(fields ...interface{}) []zap.Field {
    zapFields := make([]zap.Field, len(fields))
    for i, field := range fields {
        zapFields[i] = zap.Any(fmt.Sprintf("field_%d", i), field)  // Dynamically assign field names
    }
    return zapFields
}

// Info logs a message at the INFO level with the given fields.
func (z *zapLogger) Info(msg string, fields ...interface{}) {
    z.logger.Info(msg, toZapFields(fields...)...)
}

// Warn logs a message at the WARN level with the given fields.
func (z *zapLogger) Warn(msg string, fields ...interface{}) {
    z.logger.Warn(msg, toZapFields(fields...)...)
}

// Debug logs a message at the DEBUG level with the given fields, if the log level is
// DEBUG or lower.
func (z *zapLogger) Debug(msg string, fields ...interface{}) {
    if z.logLevel <= zapcore.DebugLevel {
        z.logger.Debug(msg, toZapFields(fields...)...)
    }
}

// Error logs a message at the ERROR level with the given fields and a stack trace.
func (z *zapLogger) Error(msg string, fields ...interface{}) {
    z.logger.Error(msg, append(toZapFields(fields...), captureStackTrace())...)
}

// GetLevel returns the string representation of the current log level.
// This is useful for logging and debugging.
func (z *zapLogger) GetLevel() string {
    return z.logLevel.String()
}

// Capture stack trace as Zap field
func captureStackTrace() zap.Field {
    pc := make([]uintptr, 10)
    runtime.Callers(3, pc) // Skip 3 frames
    frames := runtime.CallersFrames(pc)
    var stacktrace string
    for frame, more := frames.Next(); more; frame, more = frames.Next() {
        stacktrace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
    }
    return zap.String("stacktrace", stacktrace)
}