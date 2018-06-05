package addrvrf_test

import (
	"strings"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
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

		envelope := &addrvrf.Envelope{
			Input: addrvrf.AddressInput{
				City: "city",
			},
		}

		return &fixture{
			handler:       handler,
			envelope:      envelope,
			inputChannel:  in,
			outputChannel: out,
		}
	}

	t.Run("Envelope goes in, is processed and goes out", func(t *testing.T) {
		f := setup()

		f.inputChannel <- f.envelope
		close(f.inputChannel)

		f.handler.Handle()
		close(f.outputChannel)

		processedEnvelope := <-f.outputChannel
		assert.Equal(t, f.envelope, processedEnvelope)
		assert.Equal(t, strings.ToUpper(f.envelope.Input.City), processedEnvelope.Output.City)
	})
}

type fakeVerifier struct{}

func (v *fakeVerifier) Verify(i addrvrf.AddressInput) addrvrf.AddressOutput {
	return addrvrf.AddressOutput{
		City: strings.ToUpper(i.City),
	}
}
