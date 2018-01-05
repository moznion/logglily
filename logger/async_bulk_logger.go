package logger

import (
	"sync"

	"bytes"

	"encoding/json"

	"time"

	"errors"

	"github.com/moznion/logglily/api"
	internalAPI "github.com/moznion/logglily/internal/api"
)

// AsyncBulkLogger is a loggly logger with bulk API asynchronously.
//
// This logger works as like following;
// 1. Log() method buffers the message. This method returns error, error channel and failed messages list channel immediately.
// 2. If bulkSizeThreshold of messages that are in the buffer is exceeded, logger calls loggly bulk logging API on the background goroutine.
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
type AsyncBulkLogger struct {
	APIClient              api.Client
	currentPayloadSize     int
	logs                   [][]byte
	mutex                  *sync.Mutex
	flushMutex             *sync.Mutex
	bulkSizeThreshold      int
	active                 bool
	flushTickerStoppedChan chan struct{}
	stopFlushTickerChan    chan struct{}
}

// NewAsyncBulkLogger creates an instance of AsyncBulkLogger.
//
// `flushIntervalMillis` is an interval millisecond that is used to determine the interval to flush periodically.
// If `flushIntervalMillis` is less or equal to 0, periodically flushing is disabled.
//
// `bulkByteSizeThreshold` is a threshold byte size that is used to split a chunk of bulk API payload.
// Ref: https://www.loggly.com/docs/http-bulk-endpoint/
func NewAsyncBulkLogger(tags []string, token string, isHTTPS bool, bulkByteSizeThreshold int, flushIntervalMillis int) (*AsyncBulkLogger, error) {
	if err := validateBulkByteSizeThreshold(bulkByteSizeThreshold); err != nil {
		return nil, err
	}

	l := &AsyncBulkLogger{
		APIClient:              internalAPI.NewSimpleClient(tags, token, isHTTPS),
		currentPayloadSize:     0,
		bulkSizeThreshold:      bulkByteSizeThreshold,
		active:                 true,
		mutex:                  &sync.Mutex{},
		flushMutex:             &sync.Mutex{},
		flushTickerStoppedChan: make(chan struct{}, 1),
		stopFlushTickerChan:    make(chan struct{}, 1),
	}

	l.startPeriodicallyFlushing(flushIntervalMillis)

	return l, nil
}

// Log logs the message into loggly as a bulk asynchronously.
//
// Return value of `error` is a foreground error, it's not background/async one.
// Highly recommend: this value should be cared on the production.
func (l *AsyncBulkLogger) Log(message Message) (*AsyncBulkResult, error) {
	asyncErrChan := make(chan error, 1)
	failedMessagesChan := make(chan [][]byte, 1)

	if !l.active {
		err := errors.New("in progress to shutdown. refused the message")
		asyncErrChan <- err
		failedMessagesChan <- nil
		return &AsyncBulkResult{
			AsyncErrChan:       asyncErrChan,
			FailedMessagesChan: failedMessagesChan,
		}, err
	}

	body, err := json.Marshal(message)
	if err != nil {
		asyncErrChan <- err
		failedMessagesChan <- nil
		return &AsyncBulkResult{
			AsyncErrChan:       asyncErrChan,
			FailedMessagesChan: failedMessagesChan,
		}, err
	}

	go func() {
		l.post(body, asyncErrChan, failedMessagesChan)
	}()

	return &AsyncBulkResult{
		AsyncErrChan:       asyncErrChan,
		FailedMessagesChan: failedMessagesChan,
	}, nil
}

// Flush flushes remained messages that are in the buffer.
func (l *AsyncBulkLogger) Flush() *AsyncBulkResult {
	asyncErrChan := make(chan error, 1)
	failedMessagesChan := make(chan [][]byte, 1)

	go l.flush(asyncErrChan, failedMessagesChan, l.bufferInitializer)

	return &AsyncBulkResult{
		AsyncErrChan:       asyncErrChan,
		FailedMessagesChan: failedMessagesChan,
	}
}

// Shutdown attempts to shutting down.
func (l *AsyncBulkLogger) Shutdown() *AsyncBulkResult {
	asyncErrorChan := make(chan error, 1)
	failedMessageChan := make(chan [][]byte, 1)

	go func() {
		l.active = false
		l.stopFlushTickerChan <- notifier
		<-l.flushTickerStoppedChan

		l.flush(asyncErrorChan, failedMessageChan, l.bufferInitializer)
	}()

	return &AsyncBulkResult{
		AsyncErrChan:       asyncErrorChan,
		FailedMessagesChan: failedMessageChan,
	}
}

func (l *AsyncBulkLogger) bufferInitializer() {
	l.logs = l.logs[:0]
	l.currentPayloadSize = 0
}

func (l *AsyncBulkLogger) startPeriodicallyFlushing(flushIntervalMillis int) {
	if flushIntervalMillis <= 0 {
		// Disable frequently flushing
		l.flushTickerStoppedChan <- notifier
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(flushIntervalMillis) * time.Millisecond)

	loop:
		for {
			select {
			case <-ticker.C:
				// TODO: Should it be notifier channel that connect to outer?
				errChan := make(chan error, 1)
				failedMessageChan := make(chan [][]byte, 1)
				l.flush(errChan, failedMessageChan, l.bufferInitializer)
			case <-l.stopFlushTickerChan:
				break loop
			}
		}

		l.flushTickerStoppedChan <- notifier
	}()
}

func (l *AsyncBulkLogger) flush(errChan chan error, failedMessageChan chan [][]byte, bufferSweeper func()) {
	l.flushMutex.Lock()
	defer l.flushMutex.Unlock()
	defer bufferSweeper()

	if len(l.logs) <= 0 {
		// NOP
		errChan <- nil
		failedMessageChan <- nil
		return
	}

	payload := bytes.Join(l.logs, newlineCharByte)
	resp, err := l.APIClient.LogAsBulk(payload)
	if err != nil {
		errChan <- err
		failedMessageChan <- l.logs
		return
	}
	defer resp.Body.Close()

	if err := checkHTTPResponse(resp); err != nil {
		errChan <- err
		failedMessageChan <- l.logs
		return
	}

	errChan <- nil
	failedMessageChan <- nil
}

func (l *AsyncBulkLogger) post(body []byte, errChan chan error, failedMessagesChan chan [][]byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	bodySize := len(body)
	if bodySize+l.currentPayloadSize < l.bulkSizeThreshold {
		l.logs = append(l.logs, body)
		l.currentPayloadSize += bodySize + 1
		//                               ~~~ size of newline character

		errChan <- nil
		failedMessagesChan <- nil
		return
	}

	// Over the threshold. Post payloads.
	l.flush(errChan, failedMessagesChan, func() {
		l.logs = [][]byte{body}
		l.currentPayloadSize = bodySize + 1
	})
}
