package logging

import (
	"log/slog"
	"os"
	"unicode"
)

// slogAdapter wraps a *slog.Logger to implement ILogger.
type slogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter constructs a new SlogAdapter using a text handler
// that writes to stdout and includes caller information.
func NewSlogAdapter() *slogAdapter {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	return &slogAdapter{
		logger: slog.New(handler),
	}
}

// sanitizeKey replaces any rune that isn’t letter, digit, '_' or '-'
// with an underscore, so we never get "!BADKEY." prefixes.
func sanitizeKey(key string) string {
	var out []rune
	for _, r := range key {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '_', r == '-':
			out = append(out, r)
		default:
			out = append(out, '_')
		}
	}
	return string(out)
}

// toKeyvals turns our []Field into a slice of alternating key, value,
// suitable for slog.Info(msg, keyvals…).
func (s *slogAdapter) toKeyvals(fields []Field) []any {
	out := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		out = append(out, sanitizeKey(f.Key), f.Value)
	}
	return out
}

func (s *slogAdapter) LogInfo(msg string, fields ...Field) {
	s.logger.Info(msg, s.toKeyvals(fields)...)
}

func (s *slogAdapter) LogWarning(msg string, fields ...Field) {
	s.logger.Warn(msg, s.toKeyvals(fields)...)
}

func (s *slogAdapter) LogError(msg string, fields ...Field) {
	s.logger.Error(msg, s.toKeyvals(fields)...)
}

func (s *slogAdapter) LogPanic(msg string, fields ...Field) {
	s.logger.Error(msg, s.toKeyvals(fields)...)
	panic(msg)
}
