package addrvrf_test

import (
	"testing"

	"github.com/geisonbiazus/addrvrf/internal/addrvrf"
	"github.com/geisonbiazus/addrvrf/internal/assert"
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
		f.input <- newEOFEnvelope(addrvrf.InitialSequence + 6)

		f.handler.Handle()

		assert.Equal(t, addrvrf.InitialSequence+0, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+1, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+2, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+3, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+4, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+5, (<-f.output).Sequence)
		assertClosed(t, f.output)
	})

	t.Run("Many unordered Envelopes goes through the Handler", func(t *testing.T) {
		f := setup()

		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 2)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 5)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 4)
		f.input <- newEOFEnvelope(addrvrf.InitialSequence + 6)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 0)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 1)
		f.input <- newEnvelopeWithSequence(addrvrf.InitialSequence + 3)

		f.handler.Handle()

		assert.Equal(t, addrvrf.InitialSequence+0, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+1, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+2, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+3, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+4, (<-f.output).Sequence)
		assert.Equal(t, addrvrf.InitialSequence+5, (<-f.output).Sequence)
		assertClosed(t, f.output)
	})

	t.Run("Output and Input channels are closed when EOF is sent", func(t *testing.T) {
		f := setup()

		f.input <- newEOFEnvelope(addrvrf.InitialSequence)
		f.handler.Handle()

		assertClosed(t, f.input)
		assertClosed(t, f.output)
	})
}

func newEnvelopeWithSequence(seq int) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Sequence: seq,
	}
}

func newEOFEnvelope(seq int) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Sequence: seq,
		EOF:      true,
	}
}

func assertClosed(t *testing.T, channel chan *addrvrf.Envelope) {
	t.Helper()
	_, open := <-channel

	assert.False(t, open)
}
