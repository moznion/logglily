package api

import (
	"bytes"
	"net/http"

	"fmt"

	"strings"

	"github.com/moznion/logglily/internal"
	"github.com/moznion/logglily/internal/constant/content_type"
)

type SimpleClient struct {
	logEventAPIEndpoint string
	logBulkAPIEndpoint  string
	client              *http.Client
}

func NewSimpleClient(tags []string, token string, isHTTPS bool) *SimpleClient {
	tagUnit := strings.Join(tags, ",")
	return &SimpleClient{
		logEventAPIEndpoint: buildEventAPIEndpoint(tagUnit, token, isHTTPS),
		logBulkAPIEndpoint:  buildBulkAPIEndpoint(tagUnit, token, isHTTPS),
		client:              http.DefaultClient,
	}
}

func (c *SimpleClient) Log(text []byte) (*http.Response, error) {
	return c.post(c.logEventAPIEndpoint, text)
}

func (c *SimpleClient) LogAsBulk(text []byte) (*http.Response, error) {
	return c.post(c.logBulkAPIEndpoint, text)
}

func (c *SimpleClient) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *SimpleClient) post(url string, text []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(text))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", content_type.PlainText)
	req.Header.Set("User-Agent", fmt.Sprintf("logglily/%s; https://github.com/moznion/logglily", internal.Version))
	req.Header.Set("Content-Length", string(len(text)))

	return c.client.Do(req)
}
