package addrvrf

import "io"

type Pipeline struct {
	readHandler     *ReadHandler
	verifyHandler   *VerifyHandler
	sequenceHandler *SequenceHandler
	writeHandler    *WriteHandler
	workers         int
}

func NewPipeline(input io.ReadCloser, output io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	verifyInput := make(chan *Envelope, 1024)
	sequenceInput := make(chan *Envelope, 1024)
	writeInput := make(chan *Envelope, 1024)

	return &Pipeline{
		readHandler:     NewReadHandler(input, verifyInput),
		verifyHandler:   NewVerifyHandler(verifyInput, sequenceInput, NewSmartyVerifier(client)),
		sequenceHandler: NewSequenceHandler(sequenceInput, writeInput),
		writeHandler:    NewWriteHandler(writeInput, output),
		workers:         workers,
	}
}

func (p *Pipeline) Process() {
	go p.readHandler.Handle()

	for i := 0; i < p.workers; i++ {
		go p.verifyHandler.Handle()
	}

	go p.sequenceHandler.Handle()
	p.writeHandler.Handle()
}
