package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/hjoshi123/WaaS/util"
	"go.uber.org/zap"
)

type Logger struct {
	ctx context.Context
}

// NewLogger creates a new logger from an io.Writer
func NewLogger() Logger {
	return Logger{
		ctx: context.Background(),
	}
}

// Info logs go here
func (l Logger) Info(s string) {
	util.Logger(l.ctx).Info("authentication", zap.String("info", s))
}

// Error logs go here
func (l Logger) Error(s string) {
	util.Logger(l.ctx).Error("authentication", zap.Error(errors.New(s)))
}

func (l Logger) FromRequest(req *http.Request) Logger {
	return Logger{
		ctx: req.Context(),
	}
}

func (l Logger) ContextLogger(ctx context.Context) Logger {
	return Logger{
		ctx: ctx,
	}
}
