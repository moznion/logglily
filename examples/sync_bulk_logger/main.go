package main

import (
	"os"

	"fmt"

	"github.com/moznion/logglily/logger"
)

func main() {
	token := os.Getenv("LOGGLY_TOKEN")
	tag := os.Getenv("LOGGLY_TAG")

	l, err := logger.NewSyncBulkLogger([]string{tag}, token, true, 100, 5000)
	if err != nil {
		panic(err)
	}

	result, err := l.Log(logger.Message{
		"message": "hello",
		"from":    "moznion",
		//"timestamp": "2018-01-06T06:57:48.165Z", // <= timestamp is optional. please refer to the loggly spec
	})
	if err != nil {
		// do something with failedMessages
		fmt.Println(result.FailedMessages)

		panic(err)
	}

	result, err = l.Flush()
	if err != nil {
		// do something with failedMessages
		fmt.Println(result.FailedMessages)

		panic(err)
	}

	l.Shutdown()
}
