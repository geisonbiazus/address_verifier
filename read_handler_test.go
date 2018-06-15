package addrvrf_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
)

func TestReadHandler(t *testing.T) {
	type fixture struct {
		buffer  *BufferCloserSpy
		output  chan *addrvrf.Envelope
		handler *addrvrf.ReadHandler
	}

	setup := func() *fixture {
		buffer := NewBufferCloserSpy(bytes.NewBufferString(""))
		output := make(chan *addrvrf.Envelope, 10)
		handler := addrvrf.NewReadHandler(buffer, output)

		writeLine(buffer, "Street,City,State,ZIPCode")

		return &fixture{
			buffer:  buffer,
			output:  output,
			handler: handler,
		}
	}

	t.Run("Read a CSV line", func(t *testing.T) {
		f := setup()

		writeLine(f.buffer, "A1,B1,C1,D1")

		err := f.handler.Handle()

		assertEnvelopeSent(t, addrvrf.InitialSequence, <-f.output)
		assert.True(t, f.buffer.Closed)
		assert.Nil(t, err)
	})

	t.Run("Read multiple lines and create Envelopes", func(t *testing.T) {
		f := setup()

		writeLine(f.buffer, "A1,B1,C1,D1")
		writeLine(f.buffer, "A2,B2,C2,D2")
		writeLine(f.buffer, "A3,B3,C3,D3")
		writeLine(f.buffer, "A4,B4,C4,D4")
		writeLine(f.buffer, "A5,B5,C5,D5")

		err := f.handler.Handle()

		assertEnvelopeSent(t, addrvrf.InitialSequence, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+1, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+2, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+3, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+4, <-f.output)
		assert.True(t, f.buffer.Closed)
		assert.Nil(t, err)
	})

	t.Run("Malformed file", func(t *testing.T) {
		f := setup()
		writeLine(f.buffer, "A")

		err := f.handler.Handle()

		assert.NotNil(t, err)
	})

	t.Run("Output channel is closed in the end", func(t *testing.T) {
		f := setup()
		writeLine(f.buffer, "A1,B1,C1,D1")

		f.handler.Handle()

		<-f.output
		_, open := <-f.output

		assert.False(t, open)
	})
}

func writeLine(buffer *BufferCloserSpy, line string) {
	buffer.WriteString(line + "\n")
}

func assertEnvelopeSent(t *testing.T, seq int, actual *addrvrf.Envelope) {
	num := strconv.Itoa(seq)
	expected := &addrvrf.Envelope{
		Sequence: seq,
		Input: addrvrf.AddressInput{
			Street:  "A" + num,
			City:    "B" + num,
			State:   "C" + num,
			ZIPCode: "D" + num,
		},
	}

	assert.DeepEqual(t, expected, actual)
}
