package logger

import (
	"testing"

	"time"

	"github.com/moznion/logglily/internal/api"
)

func TestNewSyncBulkLogger_BulkSizeThresholdValidation(t *testing.T) {
	_, err := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 0, 10000000)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}

	_, err = NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 5*1024*1024+1, 10000000)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}

func TestSyncBulkLoggerLogShouldBeSuccessfully(t *testing.T) {
	l, err := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummySuccClient{}

	stdout, _, err := captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if stdout != "" {
		// 72 byte
		t.Errorf("stdout should be empty, but something come: %v", stdout)
	}

	stdout, _, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if stdout != "" {
		// 144 byte
		t.Errorf("stdout should be empty, but something come: %v", stdout)
	}

	stdout, _, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	expected := `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}
{"From":"john","Message":"msg2","timestamp":"2018-01-05T17:11:25.495Z"}`
	if stdout != expected {
		// 215 byte
		t.Errorf("stdout == `%v` but wants `%v`", stdout, expected)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected = `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestSyncBulkLoggerLogShouldBeFailWithError(t *testing.T) {
	l, err := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummyErrClient{}

	_, result, err := captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if result.FailedMessages != nil {
		t.Error("result.FailedMessages is not nil but it should be nil")
	}

	_, result, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if len(result.FailedMessages) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be empty", len(result.FailedMessages))
	}

	_, result, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
		return result, err
	})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
	if len(result.FailedMessages) != 2 {
		t.Errorf("len(failedMessagesList) == %v but it should be 2", len(result.FailedMessages))
	}

	expected := `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}`
	if string(result.FailedMessages[0]) != expected {
		t.Errorf("failedMessagesList[0] == %s but it should be %v", result.FailedMessages[0], expected)
	}

	expected = `{"From":"john","Message":"msg2","timestamp":"2018-01-05T17:11:25.495Z"}`
	if string(result.FailedMessages[1]) != expected {
		t.Errorf("failedMessagesList[1] == %s but it should be %v", result.FailedMessages[1], expected)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected = `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestSyncBulkLoggerLogShouldBeFailWithHTTPFailed(t *testing.T) {
	l, err := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummyHTTPFailClient{}

	_, result, err := captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if len(result.FailedMessages) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be empty", len(result.FailedMessages))
	}

	_, result, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
		return result, err
	})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if len(result.FailedMessages) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be empty", len(result.FailedMessages))
	}

	_, result, err = captureLogStdoutCaptureWithFailedMessagesList(func() (*SyncBulkResult, error) {
		result, err := l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
		return result, err
	})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
	if len(result.FailedMessages) != 2 {
		t.Errorf("len(failedMessagesList) == %v but it should be 2", len(result.FailedMessages))
	}

	expected := `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}`
	if string(result.FailedMessages[0]) != expected {
		t.Errorf("failedMessagesList[0] == %s but it should be %v", result.FailedMessages[0], expected)
	}

	expected = `{"From":"john","Message":"msg2","timestamp":"2018-01-05T17:11:25.495Z"}`
	if string(result.FailedMessages[1]) != expected {
		t.Errorf("failedMessagesList[1] == %s but it should be %v", result.FailedMessages[1], expected)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected = `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestSyncBulkLogger_Shutdown(t *testing.T) {
	l, _ := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 1000)

	l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})

	if len(l.logs) != 1 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}

	l.Shutdown()

	if len(l.logs) != 0 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 0)
	}
	if l.currentPayloadSize != 0 {
		t.Errorf("l.currentPayloadSize == %d but wants %d", l.currentPayloadSize, 0)
	}
	if l.active {
		t.Error("l.active == true but wants false")
	}

	_, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}

func TestSyncBulkLogger_PeriodicallyFlush(t *testing.T) {
	l, _ := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 1000)
	l.APIClient = &api.DummySuccClient{}

	l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if len(l.logs) != 1 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}
	if l.currentPayloadSize == 0 {
		t.Error("l.currentPayloadSize should not be 0 but come 0")
	}

	time.Sleep(time.Duration(1500) * time.Millisecond)

	if len(l.logs) != 0 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}
	if l.currentPayloadSize != 0 {
		t.Error("l.currentPayloadSize should be 0 but it is not")
	}

	l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if len(l.logs) != 2 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 2)
	}
	if l.currentPayloadSize == 0 {
		t.Error("l.currentPayloadSize should not be 0 but come 0")
	}

	time.Sleep(time.Duration(1500) * time.Millisecond)

	if len(l.logs) != 0 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 0)
	}
	if l.currentPayloadSize != 0 {
		t.Error("l.currentPayloadSize should be 0 but it is not")
	}
}

func TestSyncBulkLogger_DisabledPeriodicallyFlush(t *testing.T) {
	l, _ := NewSyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 0)
	l.APIClient = &api.DummySuccClient{}

	l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if len(l.logs) != 1 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}
	if l.currentPayloadSize == 0 {
		t.Error("l.currentPayloadSize should not be 0 but come 0")
	}

	time.Sleep(time.Duration(1500) * time.Millisecond)

	if len(l.logs) != 1 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}
	if l.currentPayloadSize == 0 {
		t.Error("l.currentPayloadSize should not be 0 but come 0")
	}
}
