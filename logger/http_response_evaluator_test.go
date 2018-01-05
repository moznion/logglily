package logger

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

func TestCheckHTTPResponseShouldBeSuccessfully(t *testing.T) {
	var err error

	err = checkHTTPResponse(&http.Response{StatusCode: 200})
	if err != nil {
		t.Error("unexpected err", err)
	}

	err = checkHTTPResponse(&http.Response{StatusCode: rand.Intn(100) + 200})
	if err != nil {
		t.Error("unexpected err", err)
	}

	err = checkHTTPResponse(&http.Response{StatusCode: 299})
	if err != nil {
		t.Error("unexpected err", err)
	}
}

func TestCheckHTTPResponseShouldBeFail(t *testing.T) {
	var err error

	err = checkHTTPResponse(&http.Response{
		StatusCode: 199,
		Body:       ioutil.NopCloser(strings.NewReader("1xx"))},
	)
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}

	err = checkHTTPResponse(&http.Response{
		StatusCode: 300,
		Body:       ioutil.NopCloser(strings.NewReader("3xx")),
	})
	if err == nil {
		t.Error("err should not be nil, but got nil")
	}
}
