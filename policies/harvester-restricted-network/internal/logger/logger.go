package logger

import (
	"context"

	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context) context.Context {
	logger := newLogger()

	return context.WithValue(ctx, ctxLogger{}, logger)
}

func FromContext(ctx context.Context) *onelog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*onelog.Logger); ok {
		return l
	}

	return newLogger()
}

func newLogger() *onelog.Logger {
	logWriter := kubewarden.KubewardenLogWriter{}
	logger := onelog.New(
		&logWriter,
		onelog.ALL,
	)
	return logger
}
