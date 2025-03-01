package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
)

func getData() (data MonData) {
	data.Timestamp = time.Now()

	hostname, _ := os.Hostname()
	data.HostName = hostname

	data.FreeMem = memory.FreeMemory()

	// Get CPU load (1 second average)
	cpuPercentages, _ := cpu.Percent(time.Second, false)
	data.CPULoad = int(cpuPercentages[0]) // Convert float to int for simplicity

	// Get Disk Usage (%)
	diskStats, _ := disk.Usage("/")
	data.DiskUsage = int(diskStats.UsedPercent)

	return data
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected: ", conn.RemoteAddr())

	encoder := json.NewEncoder(conn)

	for {
		data := getData()

		err := encoder.Encode(data)

		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func server() {
	fmt.Println("Server started...!")
	listener, err := net.Listen("tcp", ":9090")
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
