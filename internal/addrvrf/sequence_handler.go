package addrvrf

type SequenceHandler struct {
	input      chan *Envelope
	output     chan *Envelope
	currentSeq int
	buffer     map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:      input,
		output:     output,
		currentSeq: InitialSequence,
		buffer:     make(map[int]*Envelope),
	}
}

func (h *SequenceHandler) Handle() {
	for e := range h.input {
		h.buffer[e.Sequence] = e
		h.sendBufferedEnvelopesInOrder()
	}
}

func (h *SequenceHandler) sendBufferedEnvelopesInOrder() {
	for {
		e, ok := h.buffer[h.currentSeq]
		if !ok {
			break
		}
		delete(h.buffer, e.Sequence)

		if e.EOF {
			h.close()
			break
		}

		h.currentSeq++
		h.output <- e
	}
}

func (h *SequenceHandler) close() {
	close(h.input)
	close(h.output)
}
