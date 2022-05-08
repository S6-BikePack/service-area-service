package mock

import "context"

type Logger struct{}

func (Logger) Close() error {
	return nil
}

func (Logger) Panic(ctx context.Context, args ...interface{}) {
	return
}

func (Logger) Fatal(ctx context.Context, args ...interface{}) {
	return
}

func (Logger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (Logger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (Logger) Warning(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}

func (Logger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	return
}
