package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// zerologAdapter wraps a zerolog.Logger to implement ILogger.
type zerologAdapter struct {
    logger zerolog.Logger
}

func NewZerologAdapter() *zerologAdapter {
    return &zerologAdapter{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger().With().Caller().Logger(),
	}
}

func (z *zerologAdapter) LogInfo(msg string, fields ...Field) {
    e := z.logger.Info()
    for _, f := range fields {
        e = e.Interface(f.Key, f.Value)
    }
    e.Msg(msg)
}

func (z *zerologAdapter) LogWarning(msg string, fields ...Field) {
    e := z.logger.Warn()
    for _, f := range fields {
        e = e.Interface(f.Key, f.Value)
    }
    e.Msg(msg)
}

func (z *zerologAdapter) LogError(msg string, fields ...Field) {
    e := z.logger.Error()
    for _, f := range fields {
        e = e.Interface(f.Key, f.Value)
    }
    e.Msg(msg)
}

func (z *zerologAdapter) LogPanic(msg string, fields ...Field) {
    e := z.logger.Panic()
    for _, f := range fields {
        e = e.Interface(f.Key, f.Value)
    }
    e.Msg(msg)
}
