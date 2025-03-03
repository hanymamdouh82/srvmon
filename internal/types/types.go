package types

import "time"

type RunFlags struct {
	ConfFilePath  string
	LogToTerminal bool
	RunAsServer   bool
}

type ServerConfig struct {
	Port        int      `yaml:"port"`
	Interval    int      `yaml:"interval"`
	Disks       []string `yaml:"disks"`
	Processes   []string `yaml:"processes"`
	LogDuration int      `yaml:"logduration"`
}

type ProcessInfo struct {
	Name        string   `json:"Name"`
	Status      []string `json:"status"`
	MemoryUsage uint64   `json:"memoryUsage"`
	CPUUsage    float64  `json:"cpuUsage"`
	ProcessID   int32    `json:"processId"`
	ProcessLogs string   `json:"processLogs"`
}

type MonData struct {
	Timestamp   time.Time      `json:"timestamp"`
	HostName    string         `json:"hostname"`
	FreeMem     uint64         `json:"freemem"`
	TotalMemory uint64         `json:"totalMemory"`
	DiskUsage   map[string]int `json:"diskusage"`
	CPULoad     int            `json:"cpuload"`
	Processes   []ProcessInfo  `json:"processes"`
}

// ServerConfig represents the structure of the YAML file
type Servers struct {
	Servers []Server `yaml:"servers"`
}

// Server represents an individual server
type Server struct {
	ServerName string `yaml:"servername"`
	Address    string `yaml:"address"`
}
