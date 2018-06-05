package address_verifier

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
