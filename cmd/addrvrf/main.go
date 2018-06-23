package main

import (
	"net/http"
	"os"

	"github.com/geisonbiazus/addrvrf"
)

func main() {
	client := addrvrf.NewAuthorizerClient(http.DefaultClient, "AUTH_ID", "AUTH_TOKEN")
	pipeline := addrvrf.NewPipeline(os.Stdin, os.Stdout, client, 8)
	pipeline.Process()
}