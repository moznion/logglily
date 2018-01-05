package api

import (
	"fmt"
	"testing"
)

func TestInstantiateHTTP(t *testing.T) {
	tags := []string{"test"}
	token := "testToken"
	client := NewSimpleClient(tags, token, false)

	expected := fmt.Sprintf("http://logs-01.loggly.com/bulk/%s/tag/%s/", token, "test")
	if client.logBulkAPIEndpoint != expected {
		t.Errorf("got == `%v` but wants `%v`", client.logBulkAPIEndpoint, expected)
	}

	expected = fmt.Sprintf("http://logs-01.loggly.com/inputs/%s/tag/%s/", token, "test")
	if client.logEventAPIEndpoint != expected {
		t.Errorf("got == `%v` but wants `%v`", client.logBulkAPIEndpoint, expected)
	}
}

func TestInstantiateHTTPS(t *testing.T) {
	tags := []string{"test"}
	token := "testToken"
	client := NewSimpleClient(tags, token, true)

	expected := fmt.Sprintf("https://logs-01.loggly.com/bulk/%s/tag/%s/", token, "test")
	if client.logBulkAPIEndpoint != expected {
		t.Errorf("got == `%v` but wants `%v`", client.logBulkAPIEndpoint, expected)
	}

	expected = fmt.Sprintf("https://logs-01.loggly.com/inputs/%s/tag/%s/", token, "test")
	if client.logEventAPIEndpoint != expected {
		t.Errorf("got == `%v` but wants `%v`", client.logBulkAPIEndpoint, expected)
	}
}
