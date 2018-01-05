package logger

import (
	"testing"

	"github.com/moznion/logglily/internal/api"
)

func TestAsyncPoolLoggerLogShouldBeSuccessfully(t *testing.T) {
	l := NewAsyncPoolLogger([]string{"test-tag"}, "test-token", true, 3, 100000)
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

	for _, errChan := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-errChan; err != nil {
			t.Error("unexpected error", err)
		}
	}
}

func TestAsyncPoolLoggerLogShouldBeFailWithError(t *testing.T) {
	l := NewAsyncPoolLogger([]string{"test-tag"}, "test-token", true, 3, 100000)
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

	for _, errChan := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-errChan; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}

func TestAsyncPoolLoggerLogShouldBeFailWhenJSONMarshalingIsFailed(t *testing.T) {
	l := NewAsyncPoolLogger([]string{"test-tag"}, "test-token", true, 3, 100000)
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

	for _, errChan := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-errChan; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}

func TestAsyncPoolLoggerLogShouldBeFailWithHTTPFailed(t *testing.T) {
	l := NewAsyncPoolLogger([]string{"test-tag"}, "test-token", true, 3, 100000)
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

	for _, errChan := range []chan error{result1.AsyncErrChan, result2.AsyncErrChan, result3.AsyncErrChan, result4.AsyncErrChan, result5.AsyncErrChan} {
		if err := <-errChan; err == nil {
			t.Error("err should not be nil, but got nil")
		}
	}
}

func TestAsyncPoolLogger_Shutdown(t *testing.T) {
	l := NewAsyncPoolLogger([]string{"test-tag"}, "test-token", true, 10, 100000)
	l.APIClient = &api.DummySuccClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	l.Log(payload)

	terminatedChan := l.Shutdown()

	<-terminatedChan

	if l.active {
		t.Error("l.active == true but wants false")
	}

	_, chanOpened := <-l.logsQueue
	if chanOpened {
		t.Error("chanOpened == true but wants false")
	}

	result, err := l.Log(payload)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
	if ec := <-result.AsyncErrChan; ec != err {
		t.Error("error and payload of channel error are different")
	}
}
