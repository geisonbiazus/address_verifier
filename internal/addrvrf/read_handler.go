package addrvrf

import (
	"encoding/csv"
	"io"
)

type ReadHandler struct {
	reader   *csv.Reader
	closer   io.Closer
	output   chan *Envelope
	sequence int
}

func NewReadHandler(readCloser io.ReadCloser, output chan *Envelope) *ReadHandler {
	return &ReadHandler{
		reader:   csv.NewReader(readCloser),
		closer:   readCloser,
		output:   output,
		sequence: InitialSequence,
	}
}

func (h *ReadHandler) Handle() error {
	defer h.close()

	h.skipHeader()

	for {
		record, err := h.reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		h.sendEnvelope(record)
	}
	h.sendEOF()

	return nil
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

func (h *ReadHandler) sendEOF() {
	h.output <- &Envelope{
		Sequence: h.sequence,
		EOF:      true,
	}
}

func (h *ReadHandler) close() {
	h.closer.Close()
	close(h.output)
}
