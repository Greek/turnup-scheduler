package retry

import (
	"runtime"
	"time"
	"turnup-scheduler/internal/logging"
)

// Retry retries the given function `fn` up to `retries` times, waiting `delay` between attempts.
// If fn returns nil, Retry returns nil immediately. If all attempts fail, returns the last error.
func Retry[T any](retries int, delay time.Duration, fn func() (T, error)) (T, error) {
	pc, _, _, _ := runtime.Caller(0)
	runtimeFn := runtime.FuncForPC(pc)
	var (
		val    T
		err    error
		fnName string
	)

	if runtimeFn == nil {
		fnName = "method unknown"
	} else {
		fnName = runtimeFn.Name()
	}

	logger := logging.BuildLogger(fnName)
	for range retries {
		val, err = fn()
		if err == nil {
			return val, nil
		}

		logger.Warn("Retrying " + fnName)
	}
	return val, err
}
