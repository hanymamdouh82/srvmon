package main

import (
	"flag"
)

var runAsServer bool
var logToTerminal bool

func main() {
	// Parse CLI arguments
	flag.BoolVar(&runAsServer, "server", false, "Run as server")
	flag.BoolVar(&logToTerminal, "log", false, "Enable client terminal logging")
	flag.Parse()

	if runAsServer {
		server()
	} else {
		client()
	}

}
