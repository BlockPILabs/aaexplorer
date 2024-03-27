package log

import (
	"context"
	"fmt"
	"os"
)

type logContextKey struct {
}

var _defaultLogger = NewTMLogger(NewSyncWriter(os.Stdout))

func SetDefaultLogger(logger Logger) {
	_defaultLogger = logger
}

func Context(ctx context.Context) Logger {
	value := ctx.Value(logContextKey{})
	if value == nil {
		return _defaultLogger
	}
	l, ok := value.(Logger)
	if !ok {
		return _defaultLogger
	}
	return l
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, logger)
}

func With(ctx context.Context, keyvals ...interface{}) (context.Context, Logger) {
	logger := Context(ctx).With(keyvals...)
	ctx = WithContext(ctx, logger)
	return ctx, logger
}

func L() Logger {
	return _defaultLogger
}

func MaskMsg(msg string) string {
	if len(msg) < 512 {
		return msg
	}
	runes := []rune(msg)
	return fmt.Sprintf("%s...{%d}...%s", string(runes[0:255]), len(runes), string(runes[len(runes)-255:]))

}
