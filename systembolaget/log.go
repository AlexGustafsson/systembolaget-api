package systembolaget

import (
	"context"
	"fmt"
	"log/slog"
)

type logKeyType struct{}

var logKey = logKeyType{}

func SetLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}

func GetLogger(ctx context.Context) *slog.Logger {
	value := ctx.Value(logKey)
	if value == nil {
		return slog.Default()
	}

	log, ok := value.(*slog.Logger)
	if !ok {
		panic(fmt.Errorf("invalid value for logger in logger"))
	}

	return log
}
