package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkHTTPResponse(res *http.Response) error {
	status := res.StatusCode
	if 200 <= status && status <= 299 {
		return nil
	}

	msg, _ := ioutil.ReadAll(res.Body)
	return fmt.Errorf("failed to call log API [status=%d, msg=%s]", status, msg)
}
