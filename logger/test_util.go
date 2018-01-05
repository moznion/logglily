package logger

import (
	"bytes"
	"io"
	"os"
)

func captureLogStdoutCapture(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	return out, err
}

func captureLogStdoutCaptureWithFailedMessagesList(f func() (*SyncBulkResult, error)) (string, *SyncBulkResult, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result, err := f()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	return out, result, err
}
