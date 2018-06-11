package addrvrf_test

import (
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestSequenceHandler(t *testing.T) {
	type fixture struct {
		input   chan *addrvrf.Envelope
		output  chan *addrvrf.Envelope
		handler *addrvrf.SequenceHandler
	}

	setup := func() *fixture {
		input := make(chan *addrvrf.Envelope, 10)
		output := make(chan *addrvrf.Envelope, 10)
		handler := addrvrf.NewSequenceHandler(input, output)

		return &fixture{
			input:   input,
			output:  output,
			handler: handler,
		}
	}

	t.Run("An Envelope goes in and it goes out", func(t *testing.T) {
		f := setup()

		envelope := newEnvelopeWithSequence(0)
		f.input <- envelope
		close(f.input)

		f.handler.Handle()
		close(f.output)

		assert.Equal(t, envelope, <-f.output)
	})

	t.Run("Many sorted Envelopes goes through the Handler", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithSequence(0)
		f.input <- newEnvelopeWithSequence(1)
		f.input <- newEnvelopeWithSequence(2)
		f.input <- newEnvelopeWithSequence(3)
		f.input <- newEnvelopeWithSequence(4)
		f.input <- newEnvelopeWithSequence(5)
		close(f.input)

		f.handler.Handle()
		close(f.output)

		assert.Equal(t, 0, (<-f.output).Sequence)
		assert.Equal(t, 1, (<-f.output).Sequence)
		assert.Equal(t, 2, (<-f.output).Sequence)
		assert.Equal(t, 3, (<-f.output).Sequence)
		assert.Equal(t, 4, (<-f.output).Sequence)
		assert.Equal(t, 5, (<-f.output).Sequence)
	})

	t.Run("Many unordered Envelopes goes through the Handler", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithSequence(2)
		f.input <- newEnvelopeWithSequence(5)
		f.input <- newEnvelopeWithSequence(4)
		f.input <- newEnvelopeWithSequence(0)
		f.input <- newEnvelopeWithSequence(1)
		f.input <- newEnvelopeWithSequence(3)
		close(f.input)

		f.handler.Handle()
		close(f.output)

		assert.Equal(t, 0, (<-f.output).Sequence)
		assert.Equal(t, 1, (<-f.output).Sequence)
		assert.Equal(t, 2, (<-f.output).Sequence)
		assert.Equal(t, 3, (<-f.output).Sequence)
		assert.Equal(t, 4, (<-f.output).Sequence)
		assert.Equal(t, 5, (<-f.output).Sequence)
	})
}

func newEnvelopeWithSequence(seq int) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Sequence: seq,
	}
}
