package api

import (
	"net/http"

	"fmt"
	"io/ioutil"
	"strings"
)

type DummySuccClient struct {
}

func (c *DummySuccClient) Log(text []byte) (*http.Response, error) {
	fmt.Printf("%s", text)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}, nil
}

func (c *DummySuccClient) LogAsBulk(text []byte) (*http.Response, error) {
	fmt.Printf("%s", text)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}, nil
}

func (c *DummySuccClient) SetHTTPClient(client *http.Client) {
	// NOP
}

type DummyErrClient struct {
}

func (c *DummyErrClient) Log(text []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}, fmt.Errorf("error on logging: %s", text)
}

func (c *DummyErrClient) LogAsBulk(text []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}, fmt.Errorf("error on bulk logging: %s", text)
}

func (c *DummyErrClient) SetHTTPClient(client *http.Client) {
	// NOP
}

type DummyHTTPFailClient struct {
}

func (c *DummyHTTPFailClient) Log(text []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(strings.NewReader("NG")),
	}, nil
}

func (c *DummyHTTPFailClient) LogAsBulk(text []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(strings.NewReader("NG")),
	}, nil
}

func (c *DummyHTTPFailClient) SetHTTPClient(client *http.Client) {
	// NOP
}
