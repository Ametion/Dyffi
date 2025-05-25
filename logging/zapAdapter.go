package logging

import (
	"go.uber.org/zap"
)

// zapAdapter wraps a *zap.Logger to implement dyffi.Logger.
type zapAdapter struct {
	logger *zap.Logger
}

// NewZapAdapter creates a production Zap logger and wraps it.
func NewZapAdapter() (*zapAdapter) {
	l, err := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic("failed to create zap logger, error: " + err.Error())
	}
	return &zapAdapter{logger: l}
}

func (z *zapAdapter) toZapFields(fields []Field) []zap.Field {
	out := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		// let zap infer the type
		out = append(out, zap.Any(f.Key, f.Value))
	}
	return out
}

func (z *zapAdapter) LogInfo(msg string, fields ...Field) {
	z.logger.Info(msg, z.toZapFields(fields)...)
}

func (z *zapAdapter) LogWarning(msg string, fields ...Field) {
	z.logger.Warn(msg, z.toZapFields(fields)...)
}

func (z *zapAdapter) LogError(msg string, fields ...Field) {
	z.logger.Error(msg, z.toZapFields(fields)...)
}

func (z *zapAdapter) LogPanic(msg string, fields ...Field) {
	z.logger.Panic(msg, z.toZapFields(fields)...)
}
