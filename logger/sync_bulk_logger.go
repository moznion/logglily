package logger

import (
	"bytes"
	"encoding/json"

	"sync"

	"time"

	"errors"

	"github.com/moznion/logglily/api"
	internalAPI "github.com/moznion/logglily/internal/api"
)

// SyncBulkLogger is a loggly logger with bulk API synchronously.
//
// This logger works as like following;
// 1. Log() method buffers the message.
// 2. If bulkSizeThreshold of messages that are in the buffer is exceeded, logger calls loggly bulk logging API.
// 3. Else, logger only buffers the message. It postpones calling API.
//
// And if flushIntervalMillis is set, this logger flushes periodically according to the interval.
// This means this logger flushes periodically even if Log() is not called.
//
// CAUTION:
// This logger has an ability to flush periodically according to the interval.
// If periodically flushing is failed, the messages that are failed to log to loggly are lost.
// I want to fix this problem in the future, but now there is not any solution.
// If it is not allowable, please consider to stop using the periodically flushing.
type SyncBulkLogger struct {
	APIClient            api.Client
	currentPayloadSize   int
	logs                 [][]byte
	mutex                *sync.Mutex
	flushMutex           *sync.Mutex
	bulkSizeThreshold    int
	active               bool
	flushTickerStoppedCh chan struct{}
	stopFlushTickerCh    chan struct{}
}

// NewSyncBulkLogger creates an instance of SyncBulkLogger.
//
// `flushIntervalMillis` is an interval millisecond that is used to determine the interval to flush periodically.
// If `flushIntervalMillis` is less or equal to 0, periodically flushing is disabled.
//
// `bulkByteSizeThreshold` is a threshold byte size that is used to split a chunk of bulk API payload.
// Ref: https://www.loggly.com/docs/http-bulk-endpoint/
func NewSyncBulkLogger(
	tags []string,
	token string,
	isHTTPS bool,
	bulkByteSizeThreshold int,
	flushIntervalMillis int,
) (*SyncBulkLogger, error) {
	if err := validateBulkByteSizeThreshold(bulkByteSizeThreshold); err != nil {
		return nil, err
	}

	l := &SyncBulkLogger{
		APIClient:            internalAPI.NewSimpleClient(tags, token, isHTTPS),
		currentPayloadSize:   0,
		bulkSizeThreshold:    bulkByteSizeThreshold,
		mutex:                &sync.Mutex{},
		flushMutex:           &sync.Mutex{},
		active:               true,
		flushTickerStoppedCh: make(chan struct{}, 1),
		stopFlushTickerCh:    make(chan struct{}, 1),
	}

	l.startPeriodicallyFlushing(flushIntervalMillis)

	return l, nil
}

// Log logs the message into loggly as a bulk synchronously.
func (l *SyncBulkLogger) Log(message Message) (*SyncBulkResult, error) {
	if !l.active {
		return &SyncBulkResult{
			FailedMessages: nil,
		}, errors.New("in progress to shutdown. refused the message")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return &SyncBulkResult{
			FailedMessages: nil,
		}, err
	}

	return l.post(body)
}

// Flush flushes remained messages that are in the buffer.
func (l *SyncBulkLogger) Flush() (*SyncBulkResult, error) {
	return l.flush(l.bufferInitializer)
}

// Shutdown attempts to shutting down.
func (l *SyncBulkLogger) Shutdown() {
	l.active = false
	l.stopFlushTickerCh <- notifier
	<-l.flushTickerStoppedCh
	l.flush(l.bufferInitializer)
}

func (l *SyncBulkLogger) post(body []byte) (*SyncBulkResult, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	bodySize := len(body)
	if bodySize+l.currentPayloadSize < l.bulkSizeThreshold {
		l.logs = append(l.logs, body)
		l.currentPayloadSize += bodySize + 1
		//                               ~~~ size of newline character

		return &SyncBulkResult{
			FailedMessages: nil,
		}, nil
	}

	// Over the threshold. Post payloads.
	return l.flush(func() {
		l.logs = [][]byte{body}
		l.currentPayloadSize = bodySize + 1
	})
}

func (l *SyncBulkLogger) flush(bufferSweeper func()) (*SyncBulkResult, error) {
	l.flushMutex.Lock()
	defer l.flushMutex.Unlock()
	defer bufferSweeper()

	if len(l.logs) <= 0 {
		// NOP
		return &SyncBulkResult{
			FailedMessages: nil,
		}, nil
	}

	payload := bytes.Join(l.logs, newlineCharByte)
	resp, err := l.APIClient.LogAsBulk(payload)
	if err != nil {
		return &SyncBulkResult{
			FailedMessages: l.logs,
		}, err
	}
	defer resp.Body.Close()

	if err := checkHTTPResponse(resp); err != nil {
		return &SyncBulkResult{
			FailedMessages: l.logs,
		}, err
	}

	return &SyncBulkResult{
		FailedMessages: nil,
	}, nil
}

func (l *SyncBulkLogger) bufferInitializer() {
	l.logs = l.logs[:0]
	l.currentPayloadSize = 0
}

func (l *SyncBulkLogger) startPeriodicallyFlushing(flushIntervalMillis int) {
	if flushIntervalMillis <= 0 {
		// Disable frequently flushing
		l.flushTickerStoppedCh <- notifier
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(flushIntervalMillis) * time.Millisecond)

	loop:
		for {
			select {
			case <-ticker.C:
				l.flush(l.bufferInitializer)
			case <-l.stopFlushTickerCh:
				break loop
			}
		}

		l.flushTickerStoppedCh <- notifier
	}()
}
