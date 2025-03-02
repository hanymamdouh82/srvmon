package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hanymamdouh82/srvmon/internal/data"
	t "github.com/hanymamdouh82/srvmon/internal/types"
	"github.com/hanymamdouh82/srvmon/internal/utils"
	"gopkg.in/yaml.v3"
)

var conf t.ServerConfig

func loadConf(rf t.RunFlags) {
	conf.Interval = 5
	conf.Port = 9090

	data, err := os.ReadFile(rf.ConfFilePath)
	if err != nil {
		fmt.Println("failed to load conf.yaml, fallback to default configuration...!")
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		fmt.Println("Failed to parse conf.yaml, fallback to defult configuration...!")
	}

	utils.LogJSON(conf)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected: ", conn.RemoteAddr())

	encoder := json.NewEncoder(conn)

	for {
		data := data.GetData(&conf)

		err := encoder.Encode(data)

		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		time.Sleep(time.Duration(conf.Interval) * time.Second)
	}
}

func Server(rf t.RunFlags) {
	loadConf(rf)

	fmt.Println("Server started...!")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.Port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		go handleConnection(conn)
	}
}
