package main

import (
	"os"

	"github.com/moznion/logglily/logger"
)

func main() {
	token := os.Getenv("LOGGLY_TOKEN")
	tag := os.Getenv("LOGGLY_TAG")

	l := logger.NewAsyncLogger([]string{tag}, token, true)
	result, err := l.Log(logger.Message{
		"message": "hello",
		"from":    "moznion",
		//"timestamp": "2018-01-06T15:44:48.642Z", // <= timestamp is optional. please refer to the loggly spec
	})
	if err != nil {
		panic(err)
	}

	// Please check according to the situation.
	if err := <-result.AsyncErrChan; err != nil {
		panic(err)
	}
}
