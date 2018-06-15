package addrvrf_test

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestReadHandler(t *testing.T) {
	t.Run("Read a CSV line", func(t *testing.T) {
		buffer := bytes.NewBufferString("A,B,C,D")
		output := make(chan *addrvrf.Envelope, 10)
		handler := addrvrf.NewReadHandler(buffer, output)

		handler.Handle()
		close(output)

		expected := &addrvrf.Envelope{
			Sequence: 0,
			Input: addrvrf.AddressInput{
				Street:  "A",
				City:    "B",
				State:   "C",
				ZIPCode: "D",
			},
		}

		assert.Equal(t, expected, <-output)
	})
}
