package addrvrf_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type HTTPClientSpy struct {
	Request    *http.Request
	Response   *http.Response
	Error      error
	ReadCloser *ReadCloserSpy
}

func (c *HTTPClientSpy) Do(r *http.Request) (*http.Response, error) {
	c.Request = r
	return c.Response, c.Error
}

func (c *HTTPClientSpy) Configure(status int, body string) {
	c.ReadCloser = NewReadCloserSpy((bytes.NewBufferString(body)))
	c.Response = &http.Response{
		StatusCode: status,
		Body:       c.ReadCloser,
	}
}

func (c *HTTPClientSpy) ConfigureError(err error) {
	c.Error = err
}

type ReadCloserSpy struct {
	Reader io.Reader
	Closed bool
}

func NewReadCloserSpy(r io.Reader) *ReadCloserSpy {
	return &ReadCloserSpy{Reader: r}
}

func (r *ReadCloserSpy) Close() error {
	r.Closed = true
	return nil
}

func (r *ReadCloserSpy) Read(p []byte) (int, error) {
	if r.Closed {
		return 0, errors.New("Already closed")
	}
	return r.Reader.Read(p)
}
