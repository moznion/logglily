package logger

import (
	"encoding/json"
	"sync"

	"time"

	"errors"

	"github.com/moznion/logglily/api"
	internalAPI "github.com/moznion/logglily/internal/api"
)

// AsyncPoolLogger is a loggly logger with event API asynchronously that uses goroutine pool.
//
// This logger works asynchronously. This looks similar to AsyncLogger, but that is different.
// This logger calls loggly event API on pre-spawned goroutine(s) like a job-queue system.
// For each worker pickups the message from queue and call event API to log it.
type AsyncPoolLogger struct {
	APIClient   api.Client
	logsQueue   chan *asyncLog
	workerNum   int
	wg          *sync.WaitGroup
	active      bool
	stoppedChan chan struct{}
}

type asyncLog struct {
	body    []byte
	errChan chan error
}

// NewAsyncPoolLogger creates an instance of AsyncPoolLogger.
//
// Please consider the parameter of `workerNum` and `queueSize`.
// `workerNum` is the number of worker goroutines. For each worker takes the message from the queue and call event API.
// `queueSize` is an important parameter. This is the maximum capacity of the queue.
// If the message is enqueued beyond the maximum capacity of the queue, queueing will block!
// Highly recommended: `queueSize` parameter should be mush enough.
func NewAsyncPoolLogger(tags []string, token string, isHTTPS bool, workerNum int, queueSize int) *AsyncPoolLogger {
	l := &AsyncPoolLogger{
		APIClient:   internalAPI.NewSimpleClient(tags, token, isHTTPS),
		logsQueue:   make(chan *asyncLog, queueSize),
		workerNum:   workerNum,
		wg:          &sync.WaitGroup{},
		active:      true,
		stoppedChan: make(chan struct{}, workerNum),
	}

	l.start(workerNum)

	return l
}

// Log logs message into loggly through event API asynchronously.
//
// This method only enqueues the message into the queue and returns error and result immediately.
// The message will be processed by pre-spawned goroutine, like as job-queue.
//
// If it is necessary to check the status of API calls, please check the error channel that is in the result.
//
// Return value of `error` is a foreground error, it's not background/async one.
// Highly recommend: this value should be cared on the production.
func (l *AsyncPoolLogger) Log(message Message) (*AsyncResult, error) {
	asyncErrChan := make(chan error, 1)

	if !l.active {
		err := errors.New("in progress to shutdown. refused the message")
		asyncErrChan <- err
		return &AsyncResult{
			AsyncErrChan: asyncErrChan,
		}, err
	}

	body, err := json.Marshal(message)
	if err != nil {
		asyncErrChan <- err
		return &AsyncResult{
			AsyncErrChan: asyncErrChan,
		}, err
	}

	l.logsQueue <- &asyncLog{
		body:    body,
		errChan: asyncErrChan,
	}

	return &AsyncResult{
		AsyncErrChan: asyncErrChan,
	}, nil
}

// Shutdown shutdowns the logger.
//
// This method returns channel immediately; that means this method doesn't wait for the completion of shutting down.
// If you want to detect whether shutting down is completed or not, please check the channel of return value.
//
// This method works as following;
// 1. Wait for all messages that are in the queue are processed.
// 2. Close the queue (queue is a channel).
// 3. Wait for all workers are terminated.
//
// If it must shutdown immediately without waiting for post-processing, please consider using ShutdownForce().
//
// NOTE: Do not reuse the instances that you shutdown.
func (l *AsyncPoolLogger) Shutdown() chan struct{} {
	shutdownCompletedChan := make(chan struct{}, 1)

	go func() {
		l.active = false
		duration := 100 * time.Millisecond
		for {
			if len(l.logsQueue) <= 0 {
				close(l.logsQueue)

				for i := 0; i < l.workerNum; i++ {
					<-l.stoppedChan
				}

				break
			}

			time.Sleep(duration)
		}
		shutdownCompletedChan <- notifier
	}()

	return shutdownCompletedChan
}

// ShutdownForce shutdowns forcibly.
//
// This method closes the queue.
//
// CAUTION:
// This method dispose remained messages that are in the queue.
// If it is not capable, please consider to use
//
// NOTE: Do not reuse the instances that you shutdown.
func (l *AsyncPoolLogger) ShutdownForce() {
	l.active = true
	close(l.logsQueue)
}

func (l *AsyncPoolLogger) start(workerNum int) {
	for i := 0; i < workerNum; i++ {
		l.wg.Add(1)

		go func() {
			defer l.wg.Done()
			for {
				log, chanOpened := <-l.logsQueue
				if !chanOpened {
					// chan is closed. Abort.
					l.stoppedChan <- notifier
					return
				}

				func() {
					resp, err := l.APIClient.Log(log.body)
					if err != nil {
						log.errChan <- err
						return
					}
					defer resp.Body.Close()

					if err := checkHTTPResponse(resp); err != nil {
						log.errChan <- err
						return
					}

					log.errChan <- nil
				}()
			}
		}()
	}
}
