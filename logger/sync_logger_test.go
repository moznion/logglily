package logger

import (
	"testing"

	"encoding/json"

	"github.com/moznion/logglily/internal/api"
)

func TestSyncLogShouldBeSuccessful(t *testing.T) {
	l := NewSyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummySuccClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	stdout, err := captureLogStdoutCapture(func() error {
		return l.Log(payload)
	})
	if err != nil {
		t.Error("unexpected err", err)
	}

	expected, _ := json.Marshal(payload)
	if stdout != string(expected) {
		t.Errorf("stdout == `%v` want `%v`", stdout, expected)
	}
}

func TestSyncLogShouldBeFailWithError(t *testing.T) {
	l := NewSyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummyErrClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	_, err := captureLogStdoutCapture(func() error {
		return l.Log(payload)
	})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}

func TestSyncLogShouldBeFailWithHTTPFail(t *testing.T) {
	l := NewSyncLogger([]string{"test-tag"}, "test-token", true)
	l.APIClient = &api.DummyHTTPFailClient{}

	payload := Message{
		"Message": "test-msg",
		"From":    "john doe",
	}

	_, err := captureLogStdoutCapture(func() error {
		return l.Log(payload)
	})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}
