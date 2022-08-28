package systembolaget

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type logKeyType struct{}

var logKey = logKeyType{}

func SetLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}

func GetLogger(ctx context.Context) *zap.Logger {
	value := ctx.Value(logKey)
	if value == nil {
		return zap.NewNop()
	}

	log, ok := value.(*zap.Logger)
	if !ok {
		panic(fmt.Errorf("invalid value for logger in logger"))
	}

	return log
}
