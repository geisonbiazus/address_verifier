package addrvrf

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

func NewVerifyHandler(input, output chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:    input,
		output:   output,
		verifier: verifier,
	}
}

func (h *VerifyHandler) Handle() {
	defer h.close()

	for envelope := range h.input {
		envelope.Output = h.verifier.Verify(envelope.Input)
		h.output <- envelope
	}
}

func (h *VerifyHandler) close() {
	close(h.output)
}
