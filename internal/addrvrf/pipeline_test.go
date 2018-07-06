package addrvrf_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/geisonbiazus/addrvrf/internal/addrvrf"
	"github.com/geisonbiazus/addrvrf/internal/assert"
)

func TestPipeline(t *testing.T) {
	type fixture struct {
		input    *BufferCloserSpy
		output   *BufferCloserSpy
		client   *HTTPClientStub
		pipeline *addrvrf.Pipeline
	}

	setup := func() *fixture {
		input := NewBufferCloserSpy(&bytes.Buffer{})
		output := NewBufferCloserSpy(&bytes.Buffer{})
		client := &HTTPClientStub{}

		pipeline := addrvrf.NewPipeline(input, output, client, 8)

		return &fixture{
			input:    input,
			output:   output,
			client:   client,
			pipeline: pipeline,
		}
	}

	t.Run("Process the CSV input and generate the output", func(t *testing.T) {
		f := setup()

		f.input.WriteString("Street1,City,State,ZIPCode\n")
		f.input.WriteString("A,B,C,D\n")
		f.input.WriteString("A,B,C,D\n")

		err := f.pipeline.Process()

		assert.Equal(t, "Status,DeliveryLine1,LastLine,Street,City,State,ZIPCode\n", readLine(f.output))
		assert.Equal(t, "Deliverable,delivery line 1,last line,street,city,state,zip code\n", readLine(f.output))
		assert.Equal(t, "Deliverable,delivery line 1,last line,street,city,state,zip code\n", readLine(f.output))
		assert.Nil(t, err)
	})

	t.Run("Returns error with malformed CSV", func(t *testing.T) {
		f := setup()

		f.input.WriteString("Street1,City,State,ZIPCode\n")
		f.input.WriteString("A\n")

		err := f.pipeline.Process()

		assert.NotNil(t, err)
	})
}

type HTTPClientStub struct{}

func (c *HTTPClientStub) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       NewBufferCloserSpy(bytes.NewBufferString(SmartyAPISuccessResponse)),
	}, nil
}
