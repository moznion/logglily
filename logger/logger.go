package logger

var notifier struct{}
var newlineCharByte = []byte{10}

// Message is a structure that represents the message payload.
type Message map[string]interface{}

// AsyncBulkResult is a result structure of asynchronously bulk API calling.
type AsyncBulkResult struct {
	// AsyncErrChan is a channel that notifies the error of asynchronously processing.
	// If it is necessary to check the result status of async, please use this.
	//
	// This channel contains `nil` exactly when logger only buffers.
	// There is any possibilities that the channel contains something error,
	// in case of called loggly bulk API.
	//
	// This channel is effective to make wait/sync the processing.
	AsyncErrChan chan error

	// FailedMessagesChan is a channel that contains messages that are failed to log into loggly.
	// Recommend: these messages should be cared by high-level layer.
	//
	// Example
	//
	//    if result.Err != nil {
	//        failedMessages := <-result.FailedMessagesChan
	//        // Do something
	//    }
	//
	//    if asyncErr := <-result.AsyncErrChan {
	//        failedMessages := <-result.FailedMessagesChan
	//        // Do something
	//    }
	FailedMessagesChan chan [][]byte
}

// AsyncResult is a result structure of asynchronously event API calling.
type AsyncResult struct {
	// AsyncErrChan is a channel that notifies the error of asynchronously processing.
	// If it is necessary to check the result status of async, please use this.
	//
	// This channel is effective to make wait/sync the processing.
	AsyncErrChan chan error
}

// SyncBulkResult is a result structure of synchronously bulk API calling.
type SyncBulkResult struct {
	// FailedMessages represents messages that are failed to log into loggly.
	// Recommend: these messages should be cared by high-level layer.
	//
	// Example
	//
	//    if err != nil {
	//        failedMessages := result.FailedMessages
	//        // Do something
	//    }
	FailedMessages [][]byte
}
