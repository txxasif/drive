package util

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger that provides a consistent logging interface
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(level zapcore.Level) *Logger {
	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		level,
	)

	// Create logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{logger}
}

// timeEncoder encodes the time as RFC3339 with milliseconds precision
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// WithRequestID creates a child logger with request ID
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.With(zap.String("request_id", requestID))
}

// WithError creates a child logger with error
func (l *Logger) WithError(err error) *Logger {
	return l.With(zap.Error(err))
}

// WithDuration creates a child logger with duration
func (l *Logger) WithDuration(duration time.Duration) *Logger {
	return l.With(zap.Duration("duration", duration))
}

// WithUserID creates a child logger with user ID
func (l *Logger) WithUserID(userID uint) *Logger {
	return l.With(zap.Uint("user_id", userID))
}

// WithEmail creates a child logger with email
func (l *Logger) WithEmail(email string) *Logger {
	return l.With(zap.String("email", email))
}

// WithMethod creates a child logger with HTTP method
func (l *Logger) WithMethod(method string) *Logger {
	return l.With(zap.String("method", method))
}

// WithPath creates a child logger with HTTP path
func (l *Logger) WithPath(path string) *Logger {
	return l.With(zap.String("path", path))
}

// WithRemoteAddr creates a child logger with remote address
func (l *Logger) WithRemoteAddr(addr string) *Logger {
	return l.With(zap.String("remote_addr", addr))
}

// WithStatusCode creates a child logger with status code
func (l *Logger) WithStatusCode(code int) *Logger {
	return l.With(zap.Int("status_code", code))
}

// Helper functions for use with zap fields
func WithError(err error) zapcore.Field {
	return zap.Error(err)
}

func WithUserID(userID uint) zapcore.Field {
	return zap.Uint("user_id", userID)
}

func WithPath(path string) zapcore.Field {
	return zap.String("path", path)
}
