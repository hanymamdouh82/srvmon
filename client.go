package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	serverDataList []MonData
	mu             sync.Mutex
)

func monServer(srv Server, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", srv.Address)
	if err != nil {
		fmt.Println("Server ", srv.ServerName, " - is not reachable...!")
		return
	}
	defer conn.Close()

	fmt.Println("Connected to", srv.ServerName, " - Listenting for data...")
	decoder := json.NewDecoder(conn)

	for {
		var data MonData
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("Connection closed:", err)
			break
		}

		// Safely update the slice
		mu.Lock()
		updateServerData(data)
		mu.Unlock()
	}
}

// updateServerData updates the serverDataList with the latest data from a server
func updateServerData(newData MonData) {
	found := false
	for i, s := range serverDataList {
		if s.HostName == newData.HostName {
			serverDataList[i] = newData
			found = true
			break
		}
	}

	if !found {
		serverDataList = append(serverDataList, newData)
	}
}

// clearScreen clears the terminal screen (works on Linux/macOS)
func clearScreen() {
	fmt.Print("\033[H\033[2J") // ANSI escape codes
}

// displayLoop continuously prints the server status table
func displayLoop() {
	for {
		time.Sleep(2 * time.Second) // Refresh rate

		mu.Lock()
		printTable()
		mu.Unlock()
	}
}

// printTable prints the server data in a formatted table
func printTable() {
	clearScreen()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	fmt.Fprintln(w, "Server\tCPU Load\tFree RAM\tDisk Usage")
	fmt.Fprintln(w, "-------\t--------\t--------\t---------")

	for _, s := range serverDataList {
		freeMem := fmt.Sprintf("%dMB", s.FreeMem/1024/1024)
		diskUsage := fmt.Sprintf("%d%%", s.DiskUsage)

		if s.FreeMem < 1000 {
			freeMem += " 🔴"
		}
		if s.DiskUsage > 75 {
			diskUsage += " 🔴"
		}

		fmt.Fprintf(w, "%s\t%d%%\t%s\t%s\n", s.HostName, s.CPULoad, freeMem, diskUsage)
	}

	w.Flush()
}

// LoadServers reads servers from a YAML file
func loadServers(filePath string) ([]Server, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Servers, nil
}

func client() {

	fmt.Println("Client started...!")
	servers, err := loadServers("servers.yaml")
	if err != nil {
		fmt.Println("Failed to load servers:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, server := range servers {
		wg.Add(1)
		go monServer(server, &wg)
	}

	// Start the display update loop in a separate goroutine
	go displayLoop()
	wg.Wait()
}
