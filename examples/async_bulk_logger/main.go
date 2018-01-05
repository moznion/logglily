package main

import (
	"fmt"
	"os"

	"github.com/moznion/logglily/logger"
)

func main() {
	token := os.Getenv("LOGGLY_TOKEN")
	tag := os.Getenv("LOGGLY_TAG")

	l, err := logger.NewAsyncBulkLogger([]string{tag}, token, true, 1024*1024*3, 10000)
	if err != nil {
		panic(err)
	}

	result, err := l.Log(logger.Message{
		"message": "hello",
		"from":    "moznion",
		//"timestamp": "2018-01-06T06:57:48.165Z", // <= timestamp is optional. please refer to the loggly spec
	})
	if err != nil {
		panic(err)
	}

	if err := <-result.AsyncErrChan; err != nil { // Please check according to the situation.
		failedMessages := <-result.FailedMessagesChan

		// do something with failed messages
		fmt.Printf("%v", failedMessages)

		panic(err)
	}

	result = l.Flush()                            // <= flushes current buffer
	if err := <-result.AsyncErrChan; err != nil { // Please check according to the situation.
		failedMessages := <-result.FailedMessagesChan

		// do something with failed messages
		fmt.Printf("%v", failedMessages)

		panic(err)
	}

	result = l.Shutdown()                         // <= shutdown this logger
	if err := <-result.AsyncErrChan; err != nil { // Please check according to the situation.
		failedMessages := <-result.FailedMessagesChan

		// do something with failed messages
		fmt.Printf("%v", failedMessages)

		panic(err)
	}
}
