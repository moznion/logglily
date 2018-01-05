package logger

import (
	"testing"

	"github.com/moznion/logglily/internal/api"
)

func TestAsyncLoggerLogShouldBeSuccessfully(t *testing.T) {
	l := NewAsyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummySuccClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	result1, err1 := l.Log(payload)
	result2, err2 := l.Log(payload)
	result3, err3 := l.Log(payload)
	result4, err4 := l.Log(payload)
	result5, err5 := l.Log(payload)

	for _, err := range []error{err1, err2, err3, err4, err5} {
		if err != nil {
			t.Error("unexpected error", err)
		}
	}

	for _, result := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-result; err != nil {
			t.Error("unexpected error", err)
		}
	}
}

func TestAsyncLoggerLogShouldBeFailWithError(t *testing.T) {
	l := NewAsyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummyErrClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	result1, err1 := l.Log(payload)
	result2, err2 := l.Log(payload)
	result3, err3 := l.Log(payload)
	result4, err4 := l.Log(payload)
	result5, err5 := l.Log(payload)

	for _, err := range []error{err1, err2, err3, err4, err5} {
		if err != nil {
			// Check the error that occurs on registering the payload. It must be nil.
			t.Error("unexpected error", err)
		}
	}

	for _, result := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-result; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}

func TestAsyncLoggerLogShouldBeFailWhenJSONMarshalingIsFailed(t *testing.T) {
	l := NewAsyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummySuccClient{}

	payload := Message{
		"chan": make(chan error),
	}

	result1, err1 := l.Log(payload)
	result2, err2 := l.Log(payload)
	result3, err3 := l.Log(payload)
	result4, err4 := l.Log(payload)
	result5, err5 := l.Log(payload)

	for _, err := range []error{err1, err2, err3, err4, err5} {
		if err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}

	for _, result := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-result; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}

func TestAsyncLoggerLogShouldBeFailWithHTTPStatusFailing(t *testing.T) {
	l := NewAsyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummyHTTPFailClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	result1, err1 := l.Log(payload)
	result2, err2 := l.Log(payload)
	result3, err3 := l.Log(payload)
	result4, err4 := l.Log(payload)
	result5, err5 := l.Log(payload)

	for _, err := range []error{err1, err2, err3, err4, err5} {
		if err != nil {
			// Check the error that occurs on registering the payload. It must be nil.
			t.Error("unexpected error", err)
		}
	}

	for _, result := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-result; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}
