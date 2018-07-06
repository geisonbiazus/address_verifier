package addrvrf_test

import (
	"bytes"
	"net/http"
)

type HTTPClientSpy struct {
	Request    *http.Request
	Response   *http.Response
	Error      error
	ReadCloser *BufferCloserSpy
}

func (c *HTTPClientSpy) Do(r *http.Request) (*http.Response, error) {
	c.Request = r
	return c.Response, c.Error
}

func (c *HTTPClientSpy) Configure(status int, body string) {
	c.ReadCloser = NewBufferCloserSpy((bytes.NewBufferString(body)))
	c.Response = &http.Response{
		StatusCode: status,
		Body:       c.ReadCloser,
	}
}

func (c *HTTPClientSpy) ConfigureError(err error) {
	c.Error = err
}

type BufferCloserSpy struct {
	*bytes.Buffer
	Closed bool
}

func NewBufferCloserSpy(b *bytes.Buffer) *BufferCloserSpy {
	return &BufferCloserSpy{Buffer: b}
}

func (r *BufferCloserSpy) Close() error {
	r.Closed = true
	return nil
}
