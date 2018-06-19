package addrvrf

import "io"

type Pipeline struct {
	input   io.ReadCloser
	output  io.WriteCloser
	client  HTTPClient
	workers int
}

func NewPipeline(input io.ReadCloser, output io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		input:   input,
		output:  output,
		client:  client,
		workers: workers,
	}
}

func (p *Pipeline) Process() {
	verifyInput := make(chan *Envelope, 1024)
	sequenceInput := make(chan *Envelope, 1024)
	writeInput := make(chan *Envelope, 1024)

	go NewReadHandler(p.input, verifyInput).Handle()
	go NewVerifyHandler(verifyInput, sequenceInput, NewSmartyVerifier(p.client)).Handle()
	go NewSequenceHandler(sequenceInput, writeInput).Handle()

	NewWriteHandler(writeInput, p.output).Handle()
}
