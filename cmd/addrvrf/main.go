package main

import (
	"log"
	"net/http"
	"os"

	"github.com/geisonbiazus/addrvrf/internal/addrvrf"
)

func main() {
	client := addrvrf.NewAuthorizerClient(http.DefaultClient, "AUTH_ID", "AUTH_TOKEN")
	pipeline := addrvrf.NewPipeline(os.Stdin, os.Stdout, client, 8)
	err := pipeline.Process()
	if err != nil {
		log.Fatal(err.Error())
	}
}
