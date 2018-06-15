package addrvrf_test

import (
	"strings"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/geisonbiazus/addrvrf/assert"
)

func TestVerifyHandler(t *testing.T) {
	type fixture struct {
		handler       *addrvrf.VerifyHandler
		envelope      *addrvrf.Envelope
		inputChannel  chan *addrvrf.Envelope
		outputChannel chan *addrvrf.Envelope
	}

	setup := func() *fixture {
		in := make(chan *addrvrf.Envelope, 10)
		out := make(chan *addrvrf.Envelope, 10)
		verifier := &fakeVerifier{}

		handler := addrvrf.NewVerifyHandler(in, out, verifier)

		envelope := createEnvelope("city")

		return &fixture{
			handler:       handler,
			envelope:      envelope,
			inputChannel:  in,
			outputChannel: out,
		}
	}

	t.Run("Envelope goes in, is processed and goes out", func(t *testing.T) {
		f := setup()

		envelope := createEnvelope("city")

		f.inputChannel <- envelope
		close(f.inputChannel)

		f.handler.Handle()

		processedEnvelope := <-f.outputChannel
		assert.Equal(t, envelope, processedEnvelope)
		assert.Equal(t, strings.ToUpper(envelope.Input.City), processedEnvelope.Output.City)
	})

	t.Run("Asynchronously process all that goes through the channel", func(t *testing.T) {
		f := setup()

		go func() {
			f.handler.Handle()
		}()

		f.inputChannel <- createEnvelope("city1")
		f.inputChannel <- createEnvelope("city2")
		f.inputChannel <- createEnvelope("city3")
		close(f.inputChannel)

		assert.Equal(t, "CITY1", (<-f.outputChannel).Output.City)
		assert.Equal(t, "CITY2", (<-f.outputChannel).Output.City)
		assert.Equal(t, "CITY3", (<-f.outputChannel).Output.City)
	})

	t.Run("Output channel is closed when there is no more input", func(t *testing.T) {
		f := setup()
		close(f.inputChannel)

		f.handler.Handle()

		_, open := <-f.outputChannel

		assert.False(t, open)
	})
}

func createEnvelope(city string) *addrvrf.Envelope {
	return &addrvrf.Envelope{
		Input: addrvrf.AddressInput{
			City: city,
		},
	}
}

type fakeVerifier struct{}

func (v *fakeVerifier) Verify(i addrvrf.AddressInput) addrvrf.AddressOutput {
	return addrvrf.AddressOutput{
		City: strings.ToUpper(i.City),
	}
}
