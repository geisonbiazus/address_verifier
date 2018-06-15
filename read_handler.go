package addrvrf

import (
	"encoding/csv"
	"io"
)

type ReadHandler struct {
	reader *csv.Reader
	output chan *Envelope
}

func NewReadHandler(reader io.Reader, output chan *Envelope) *ReadHandler {
	return &ReadHandler{
		reader: csv.NewReader(reader),
		output: output,
	}
}

func (h *ReadHandler) Handle() {
	record, _ := h.reader.Read()
	h.output <- &Envelope{
		Sequence: 0,
		Input: AddressInput{
			Street:  record[0],
			City:    record[1],
			State:   record[2],
			ZIPCode: record[3],
		},
	}
}
