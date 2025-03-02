package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"sync"

	t "github.com/hanymamdouh82/srvmon/internal/types"
	"gopkg.in/yaml.v3"
)

var (
	serverDataList []t.MonData
	mu             sync.Mutex
)

func monServer(srv t.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", srv.Address)
	if err != nil {
		// fmt.Println("Server ", srv.ServerName, " - is not reachable...!")
		return
	}
	defer conn.Close()

	// fmt.Println("Connected to", srv.ServerName, " - Listenting for data...")
	decoder := json.NewDecoder(conn)

	for {
		var data t.MonData
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
func updateServerData(newData t.MonData) {
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

// LoadServers reads servers from a YAML file
func loadServers(filePath string) ([]t.Server, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config t.Servers
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Servers, nil
}

func Client() {
	hdir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to reference home directory")
	}
	serversfile := path.Join(hdir, ".config/srvmon/servers.yaml")
	servers, err := loadServers(serversfile)
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
	wg.Add(1)
	go displayLoop()
	wg.Wait()
}
