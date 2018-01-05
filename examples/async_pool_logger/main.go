package main

import (
	"os"

	"github.com/moznion/logglily/logger"
)

func main() {
	token := os.Getenv("LOGGLY_TOKEN")
	tag := os.Getenv("LOGGLY_TAG")

	l := logger.NewAsyncPoolLogger([]string{tag}, token, true, 5, 500000)
	result, err := l.Log(logger.Message{
		"message": "hello",
		"from":    "moznion",
		//"timestamp": "2018-01-06T06:57:48.165Z", // <= timestamp is optional. please refer to the loggly spec
	})
	if err != nil {
		panic(err)
	}

	// Please check according to the situation.
	if err := <-result.AsyncErrChan; err != nil {
		panic(err)
	}

	shutdownChannel := l.Shutdown()
	<-shutdownChannel // wait for completion the termination

	// l.ShutdownForce() // Other way: force shutdown
}
