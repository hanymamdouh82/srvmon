package main

import (
	"flag"

	"github.com/hanymamdouh82/srvmon/internal/client"
	"github.com/hanymamdouh82/srvmon/internal/server"
	t "github.com/hanymamdouh82/srvmon/internal/types"
)

func readFlags() (rf t.RunFlags) {
	// Parse CLI arguments
	flag.BoolVar(&rf.RunAsServer, "server", false, "Run as server")
	// flag.BoolVar(&rf.LogToTerminal, "log", false, "Enable client terminal logging")
	flag.StringVar(&rf.ConfFilePath, "conf", "", "Server configuration file path")

	flag.Parse()

	return rf
}

func main() {
	rf := readFlags()

	if rf.RunAsServer {
		server.Server(rf)
	} else {
		client.Client()
	}

}
