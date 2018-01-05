package logger

import (
	"encoding/json"

	"github.com/moznion/logglily/api"
	internalAPI "github.com/moznion/logglily/internal/api"
)

// SyncLogger is a loggly logger with event API synchronously.
//
// This logger is really simple.
// Log() method is called, this logger calls loggly event API with blocking.
//
// If performance is required, please consider using asynchronous logger.
type SyncLogger struct {
	APIClient api.Client
}

// NewSyncLogger creates an instance of SyncLogger.
func NewSyncLogger(tags []string, token string, isHTTPS bool) *SyncLogger {
	return &SyncLogger{
		APIClient: internalAPI.NewSimpleClient(tags, token, isHTTPS),
	}
}

// Log logs message into loggly through event API synchronously.
func (l *SyncLogger) Log(message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	res, err := l.APIClient.Log(body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := checkHTTPResponse(res); err != nil {
		return err
	}

	return nil
}
