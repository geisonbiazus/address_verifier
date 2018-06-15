package addrvrf_test

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
)

func TestWriteHandler(t *testing.T) {
	type fixture struct {
		input   chan *addrvrf.Envelope
		buffer  *BufferCloserSpy
		handler *addrvrf.WriteHandler
	}

	setup := func() *fixture {
		input := make(chan *addrvrf.Envelope, 10)
		buffer := NewBufferCloserSpy(&bytes.Buffer{})
		handler := addrvrf.NewWriteHandler(input, buffer)

		return &fixture{
			input:   input,
			buffer:  buffer,
			handler: handler,
		}
	}

	t.Run("Write CSV header", func(t *testing.T) {
		f := setup()
		close(f.input)

		f.handler.Handle()

		assert.Equal(t, "Status,DeliveryLine1,LastLine,Street,City,State,ZIPCode\n", readLine(f.buffer))
		assert.True(t, f.buffer.Closed)
	})

	t.Run("Write an envelope data", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithOutput("1")
		close(f.input)

		f.handler.Handle()
		readLine(f.buffer)

		assert.Equal(t, "A1,B1,C1,D1,E1,F1,G1\n", readLine(f.buffer))
		assert.True(t, f.buffer.Closed)
	})

	t.Run("Write multiple envelopes", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithOutput("1")
		f.input <- newEnvelopeWithOutput("2")
		f.input <- newEnvelopeWithOutput("3")
		f.input <- newEnvelopeWithOutput("4")
		f.input <- newEnvelopeWithOutput("5")
		close(f.input)

		f.handler.Handle()
		readLine(f.buffer)

		assert.Equal(t, "A1,B1,C1,D1,E1,F1,G1\n", readLine(f.buffer))
		assert.Equal(t, "A2,B2,C2,D2,E2,F2,G2\n", readLine(f.buffer))
		assert.Equal(t, "A3,B3,C3,D3,E3,F3,G3\n", readLine(f.buffer))
		assert.Equal(t, "A4,B4,C4,D4,E4,F4,G4\n", readLine(f.buffer))
		assert.Equal(t, "A5,B5,C5,D5,E5,F5,G5\n", readLine(f.buffer))
		assert.True(t, f.buffer.Closed)
	})
}

func readLine(buffer *BufferCloserSpy) string {
	content, _ := buffer.ReadString('\n')
	return content
}

func newEnvelopeWithOutput(seq string) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Output: addrvrf.AddressOutput{
			Status:        "A" + seq,
			DeliveryLine1: "B" + seq,
			LastLine:      "C" + seq,
			Street:        "D" + seq,
			City:          "E" + seq,
			State:         "F" + seq,
			ZIPCode:       "G" + seq,
		},
	}
}
