package main

import (
	"flag"

	"github.com/deepthawtz/paymail-go/api"
)

var baseURL, port string

func init() {
	flag.StringVar(&baseURL, "base-url", "localhost", "base URL to serve requests from")
	flag.StringVar(&port, "port", "8080", "port to listen on")
}

func main() {
	flag.Parse()
	s := api.NewServer(baseURL, port)
	s.Start()
}
