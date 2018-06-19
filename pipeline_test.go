package addrvrf_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
)

func TestPipeline(t *testing.T) {
	input := NewBufferCloserSpy(&bytes.Buffer{})
	output := NewBufferCloserSpy(&bytes.Buffer{})
	client := &HTTPClientStub{}

	pipeline := addrvrf.NewPipeline(input, output, client, 8)

	input.WriteString("Street1,City,State,ZIPCode\n")
	input.WriteString("A,B,C,D\n")
	input.WriteString("A,B,C,D\n")

	pipeline.Process()

	assert.Equal(t, "Status,DeliveryLine1,LastLine,Street,City,State,ZIPCode\n", readLine(output))
	assert.Equal(t, "Deliverable,delivery line 1,last line,street,city,state,zip code\n", readLine(output))
	assert.Equal(t, "Deliverable,delivery line 1,last line,street,city,state,zip code\n", readLine(output))
}

type HTTPClientStub struct{}

func (c *HTTPClientStub) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       NewBufferCloserSpy(bytes.NewBufferString(SmartyAPISuccessResponse)),
	}, nil
}
