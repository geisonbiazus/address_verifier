package addrvrf_test

import (
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHandler(t *testing.T) {
	t.Run("Something goes in and something goes out", func(t *testing.T) {
		in := make(chan string, 10)
		out := make(chan string, 10)
		handler := addrvrf.NewVerifyHandler(in, out)

		in <- "My String"
		close(in)

		handler.Handle()

		close(out)

		assert.Equal(t, <-out, "My String")
	})
}
