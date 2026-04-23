package logger

import (
	"context"
	"log/slog"
)

func AddFields(ctx context.Context, fields ...any) context.Context {
	return WithContext(ctx, FromContext(ctx).With(fields...))
}

func Error(err error) slog.Attr {
	return slog.Any("error", err)
}
