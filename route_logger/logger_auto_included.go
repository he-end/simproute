package logger

import (
	"sync"

	"github.com/he-end/simproute/goruntime"
)

var (
	loggerRuntimesStore sync.Map
)

type RegisterRuntime struct {
	// its use to key of log
	Key string
	// this is value for key that inputed, so the all this runtime will log with key & value which inputed
	Value string
}

// Generate logger for only this runtime
//
// The value is needed for corelation id for each level logger
func NewLoggerOnRuntime(reg RegisterRuntime) {
	loggerRuntimesStore.Store(goruntime.Goid(), reg)
	// globalLogger = globalLogger.WithOptions(zap.AddCallerSkip(1))
	// defer loggerRuntimesStore.Clear()
}

func DeferDeleteRuntimeValue() {
	loggerRuntimesStore.Delete(goruntime.Goid())
}

func GetLoggerRuntimeStore() *RegisterRuntime {
	store, ok := loggerRuntimesStore.Load(goruntime.Goid())
	if !ok {
		return nil
	}
	result := store.(RegisterRuntime)
	return &result
}

// Info()
// Warn()
// Error()
// Fatal()
// Panic()
// }
