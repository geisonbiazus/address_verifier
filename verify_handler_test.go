package address_verifier_test

import "testing"

func TestVerifyHandler(t *testing.T) {
	t.Run("Something goes in and something goes out", func(t *testing.T) {
		in := make(chan string, 1)
		out := make(chan string, 1)
		handler := NewVerifyHandler(in, out)

		in <- "My String"
		close(in)

		handler.Handle()

		close(out)
		result := <-out

		if result != "My String" {
			t.Errorf("\nGiven: %v\nWant:  %v", result, "My String")
		}
	})
}

type VerifyHandler struct {
	input  chan string
	output chan string
}

func NewVerifyHandler(input, output chan string) *VerifyHandler {
	return &VerifyHandler{input: input, output: output}
}

func (h *VerifyHandler) Handle() {
	h.output <- <-h.input
}
