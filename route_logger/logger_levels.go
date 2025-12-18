package logger

import (
	"github.com/he-end/simproute/goruntime"
	"go.uber.org/zap"
)

func Info(message string, fields ...zap.Field) {
	logSkipCaller := globalLogger.WithOptions(zap.AddCallerSkip(1))
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok || store == nil {
		logSkipCaller.Info(message, fields...)
		return
	}

	withValue := store.(RegisterRuntime)
	zapRequestID := zap.String(withValue.Key, withValue.Value)
	fields = append(fields, zapRequestID)
	logSkipCaller.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logSkipCaller := globalLogger.WithOptions(zap.AddCallerSkip(1))
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok || store == nil {
		logSkipCaller.Warn(message, fields...)
		return
	}
	withValue := store.(RegisterRuntime)
	zapRequestID := zap.String(withValue.Key, withValue.Value)
	fields = append(fields, zapRequestID)
	logSkipCaller.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logSkipCaller := globalLogger.WithOptions(zap.AddCallerSkip(1))
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok || store == nil {
		logSkipCaller.Error(message, fields...)
		return
	}
	withValue := store.(RegisterRuntime)
	zapRequestID := zap.String(withValue.Key, withValue.Value)
	fields = append(fields, zapRequestID)
	logSkipCaller.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logSkipCaller := globalLogger.WithOptions(zap.AddCallerSkip(1))
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok || store == nil {
		logSkipCaller.Fatal(message, fields...)
		return
	}

	withValue := store.(RegisterRuntime)
	zapRequestID := zap.String(withValue.Key, withValue.Value)
	fields = append(fields, zapRequestID)
	logSkipCaller.Fatal(message, fields...)
}

func Panic(message string, fields ...zap.Field) {
	logSkipCaller := globalLogger.WithOptions(zap.AddCallerSkip(1))
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok || store == nil {
		logSkipCaller.Panic(message, fields...)
		return
	}

	withValue := store.(RegisterRuntime)
	zapRequestID := zap.String(withValue.Key, withValue.Value)
	fields = append(fields, zapRequestID)
	logSkipCaller.Panic(message, fields...)
}
