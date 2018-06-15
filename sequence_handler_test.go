package addrvrf_test

import (
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
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

		envelope := newEnvelopeWithSequence(addrvrf.InitialSequence)
		f.input <- envelope
		close(f.input)

		f.handler.Handle()

		assert.Equal(t, envelope, <-f.output)
	})

	t.Run("Many sorted Envelopes goes through the Handler", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 0)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 1)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 2)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 3)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 4)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 5)
		close(f.input)

		f.handler.Handle()

		assert.Equal(t, addrvrf.InitialSequence+0, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+1, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+2, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+3, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+4, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+5, (<-f.output).Sequence)
	})

	t.Run("Many unordered Envelopes goes through the Handler", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 2)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 5)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 4)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 0)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 1)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 3)
		close(f.input)

		f.handler.Handle()

		assert.Equal(t, addrvrf.InitialSequence+0, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+1, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+2, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+3, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+4, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+5, (<-f.output).Sequence)
	})

	t.Run("Output channel is closed when there is no more input", func(t *testing.T) {
		f := setup()
		close(f.input)

		f.handler.Handle()

		_, open := <-f.output

		assert.False(t, open)
	})
}

func newEnvelopeWithSequence(seq int) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Sequence: seq,
	}
}
