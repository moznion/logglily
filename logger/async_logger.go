package logger

import (
	"encoding/json"

	"github.com/moznion/logglily/api"
	internalAPI "github.com/moznion/logglily/internal/api"
)

// AsyncLogger is a loggly logger with event API asynchronously.
//
// This logger works asynchronously.
// Log() method is called, this logger calls loggly event API on background;
// this means Log() method spawns gorountine and delegate API calling to that.
// *Thus if it is necessary to control the capacity of goroutines, please consider using AsyncPoolLogger.*
type AsyncLogger struct {
	APIClient api.Client
}

// NewAsyncLogger creates an instance of AsyncLogger.
func NewAsyncLogger(tags []string, token string, isHTTPS bool) *AsyncLogger {
	return &AsyncLogger{
		APIClient: internalAPI.NewSimpleClient(tags, token, isHTTPS),
	}
}

// Log logs message into loggly through event API asynchronously.
//
// This method calls loggly event API on the spawned goroutine.
// This method doesn't block while API calling; so returns result and error immediately.
// If it is necessary to check the status of async API calls, please check the error channel that is in result.
//
// Return value of `error` is a foreground error, it's not background/async one.
// Highly recommend: this value should be cared on the production.
//
// This method spawns goroutine each time this method is called.
// It cannot control the capacity of goroutines, so if it is necessary to control that, please consider using AsyncPoolLogger.
func (l *AsyncLogger) Log(message Message) (*AsyncResult, error) {
	asyncErrChan := make(chan error, 1)

	body, err := json.Marshal(message)
	if err != nil {
		asyncErrChan <- err
		return &AsyncResult{
			AsyncErrChan: asyncErrChan,
		}, err
	}

	go func() {
		resp, err := l.APIClient.Log(body)
		if err != nil {
			asyncErrChan <- err
			return
		}
		defer resp.Body.Close()

		if err := checkHTTPResponse(resp); err != nil {
			asyncErrChan <- err
			return
		}

		asyncErrChan <- nil
	}()

	return &AsyncResult{
		AsyncErrChan: asyncErrChan,
	}, nil
}
