package addrvrf

import (
	"encoding/csv"
	"io"
)

type ReadHandler struct {
	reader   *csv.Reader
	output   chan *Envelope
	sequence int
}

func NewReadHandler(reader io.Reader, output chan *Envelope) *ReadHandler {
	return &ReadHandler{
		reader:   csv.NewReader(reader),
		output:   output,
		sequence: InitialSequence,
	}
}

func (h *ReadHandler) Handle() {
	h.skipHeader()

	for {
		record, err := h.reader.Read()
		if err == io.EOF {
			break
		}
		h.sendEnvelope(record)
	}
}

func (h *ReadHandler) skipHeader() {
	h.reader.Read()
}

func (h *ReadHandler) sendEnvelope(record []string) {
	h.output <- &Envelope{
		Sequence: h.sequence,
		Input: AddressInput{
			Street:  record[0],
			City:    record[1],
			State:   record[2],
			ZIPCode: record[3],
		},
	}
	h.sequence++
}
