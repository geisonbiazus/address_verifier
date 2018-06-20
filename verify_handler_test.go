package addrvrf_test

import (
	"strings"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
)

func TestVerifyHandler(t *testing.T) {
	type fixture struct {
		handler  *addrvrf.VerifyHandler
		verifier *fakeVerifier
		input    chan *addrvrf.Envelope
		output   chan *addrvrf.Envelope
	}

	setup := func() *fixture {
		in := make(chan *addrvrf.Envelope, 10)
		out := make(chan *addrvrf.Envelope, 10)
		verifier := &fakeVerifier{}

		handler := addrvrf.NewVerifyHandler(in, out, verifier)

		return &fixture{
			handler:  handler,
			verifier: verifier,
			input:    in,
			output:   out,
		}
	}

	t.Run("Envelope goes in, is processed and goes out", func(t *testing.T) {
		f := setup()

		envelope := createEnvelope("city")

		f.input <- envelope
		close(f.input)

		f.handler.Handle()
		close(f.output)

		processedEnvelope := <-f.output
		assert.Equal(t, envelope, processedEnvelope)
		assert.Equal(t, strings.ToUpper(envelope.Input.City), processedEnvelope.Output.City)
	})

	t.Run("Asynchronously process all that goes through the channel", func(t *testing.T) {
		f := setup()

		go func() {
			f.handler.Handle()
		}()

		f.input <- createEnvelope("city1")
		f.input <- createEnvelope("city2")
		f.input <- createEnvelope("city3")
		close(f.input)

		assert.Equal(t, "CITY1", (<-f.output).Output.City)
		assert.Equal(t, "CITY2", (<-f.output).Output.City)
		assert.Equal(t, "CITY3", (<-f.output).Output.City)
		close(f.output)
	})

	t.Run("EOF is ignored and passed to output directly", func(t *testing.T) {
		f := setup()

		f.input <- &addrvrf.Envelope{EOF: true}
		close(f.input)

		f.handler.Handle()
		close(f.output)

		assert.DeepEqual(t, &addrvrf.Envelope{EOF: true}, <-f.output)
		assert.False(t, f.verifier.called)
	})
}

func createEnvelope(city string) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Input: addrvrf.AddressInput{
			City: city,
		},
	}
}

type fakeVerifier struct {
	called bool
}

func (v *fakeVerifier) Verify(i addrvrf.AddressInput) addrvrf.AddressOutput {
	v.called = true
	return addrvrf.AddressOutput{
		City: strings.ToUpper(i.City),
	}
}
