package main

import (
	"flag"
	"fmt"

	"github.com/petroslafaz/basic-http-server-go/server"
)

func main() {
	// read the directory to serve files from
	// e.g go run server.go -directory /tmp
	directory := flag.String("directory", "", "the directory to serve files from")
	flag.Parse()

	fmt.Println("Serving directory:", *directory)

	server.StartServer("4221", *directory)
}
