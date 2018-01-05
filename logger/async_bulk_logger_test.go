package logger

import (
	"testing"

	"time"

	"github.com/moznion/logglily/internal/api"
)

func TestNewAsyncBulkLogger_ValidateBulkByteSizeThreshold(t *testing.T) {
	_, err := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 0, 10000000)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}

	_, err = NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 5*1024*1024+1, 10000000)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}

func TestAsyncBulkLoggerLogShouldBeSuccessfully(t *testing.T) {
	l, err := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummySuccClient{}

	result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected := `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestAsyncBulkLoggerLogShouldBeFailWithErr(t *testing.T) {
	l, err := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummyErrClient{}

	result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err == nil {
		t.Error("err should not be nil, but got nil")
	}

	failedMessagesList := <-result.FailedMessagesChan
	if len(failedMessagesList) != 2 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 2)
	}
	if string(failedMessagesList[0]) != `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}` {
		t.Errorf("failedMessageList[0] == %s but it should be %s", failedMessagesList[0], `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}`)
	}
	if string(failedMessagesList[1]) != `{"From":"john","Message":"msg2","timestamp":"2018-01-05T17:11:25.495Z"}` {
		t.Errorf("failedMessageList[1] == %s but it should be %s", failedMessagesList[1], `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.495Z"}`)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected := `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestAsyncBulkLoggerLogShouldBeFailWithHTTPFailed(t *testing.T) {
	l, err := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 10000000)
	if err != nil {
		t.Error("unexpected err", err)
	}

	l.APIClient = &api.DummyHTTPFailClient{}

	result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.495Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err != nil {
		t.Error("unexpected err", err)
	}
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 0)
	}

	result, err = l.Log(Message{"Message": "msg3", "From": "john", "timestamp": "2018-01-05T17:11:25.496Z"})
	if err != nil {
		t.Error("unexpected err", err)
	}
	if err := <-result.AsyncErrChan; err == nil {
		t.Error("err should not be nil, but got nil")
	}

	failedMessagesList := <-result.FailedMessagesChan
	if len(failedMessagesList) != 2 {
		t.Errorf("len(failedMessagesList) == %v but it should be %v", len(failedMessagesList), 2)
	}
	if string(failedMessagesList[0]) != `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}` {
		t.Errorf("failedMessageList[0] == %s but it should be %s", failedMessagesList[0], `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.494Z"}`)
	}
	if string(failedMessagesList[1]) != `{"From":"john","Message":"msg2","timestamp":"2018-01-05T17:11:25.495Z"}` {
		t.Errorf("failedMessageList[1] == %s but it should be %s", failedMessagesList[1], `{"From":"john","Message":"msg1","timestamp":"2018-01-05T17:11:25.495Z"}`)
	}

	if l.currentPayloadSize != 72 {
		t.Errorf("currentPayloadSize == %v, but wants %v", l.currentPayloadSize, 72)
	}

	if len(l.logs) != 1 {
		t.Errorf("size of logs == %v but wants %v", len(l.logs), 1)
	}

	expected := `{"From":"john","Message":"msg3","timestamp":"2018-01-05T17:11:25.496Z"}`
	if string(l.logs[0]) != expected {
		t.Errorf("log == %v but wants %v", string(l.logs[0]), expected)
	}
}

func TestAsyncBulkLogger_Shutdown(t *testing.T) {
	l, _ := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 1000)
	l.APIClient = &api.DummySuccClient{}

	result, _ := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})

	<-result.AsyncErrChan
	if len(l.logs) != 1 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 1)
	}

	result = l.Shutdown()
	<-result.AsyncErrChan
	if failedMessagesList := <-result.FailedMessagesChan; len(failedMessagesList) != 0 {
		t.Errorf("len(failedMessageList) == %d but wants %d", len(failedMessagesList), 0)
	}

	if len(l.logs) != 0 {
		t.Errorf("len(l.logs) == %d but wants %d", len(l.logs), 0)
	}
	if l.currentPayloadSize != 0 {
		t.Errorf("l.currentPayloadSize == %d but wants %d", l.currentPayloadSize, 0)
	}
	if l.active {
		t.Error("l.active == true but wants false")
	}

	result, err := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
	if ec := <-result.AsyncErrChan; ec != err {
		t.Error("error and payload of channel error are different")
	}
}

func TestAsyncBulkLogger_PeriodicallyFlushing(t *testing.T) {
	l, _ := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 1000)
	l.APIClient = &api.DummySuccClient{}

	result, _ := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})

	<-result.AsyncErrChan
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

	result1, _ := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	result2, _ := l.Log(Message{"Message": "msg2", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})
	<-result1.AsyncErrChan
	<-result2.AsyncErrChan
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

func TestAsyncBulkLogger_DisabledPeriodicallyFlushing(t *testing.T) {
	l, _ := NewAsyncBulkLogger([]string{"test-tag"}, "test-token", true, 215, 0)
	l.APIClient = &api.DummySuccClient{}

	result, _ := l.Log(Message{"Message": "msg1", "From": "john", "timestamp": "2018-01-05T17:11:25.494Z"})

	<-result.AsyncErrChan
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
