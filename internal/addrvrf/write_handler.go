package addrvrf

import (
	"encoding/csv"
	"io"
)

type WriteHandler struct {
	input  chan *Envelope
	closer io.Closer
	writer *csv.Writer
}

func NewWriteHandler(input chan *Envelope, writeCloser io.WriteCloser) *WriteHandler {
	return &WriteHandler{
		input:  input,
		closer: writeCloser,
		writer: csv.NewWriter(writeCloser),
	}
}

func (h *WriteHandler) Handle() {
	h.writeLine("Status", "DeliveryLine1", "LastLine", "Street", "City", "State", "ZIPCode")

	for envelope := range h.input {
		h.writeOutput(envelope.Output)
	}

	h.writer.Flush()
	h.closer.Close()
}

func (h *WriteHandler) writeLine(line ...string) {
	h.writer.Write(line)
}

func (h *WriteHandler) writeOutput(o AddressOutput) {
	h.writeLine(o.Status, o.DeliveryLine1, o.LastLine, o.Street, o.City, o.State, o.ZIPCode)
}
