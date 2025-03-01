package main

import "time"

type MonData struct {
	Timestamp time.Time `json:"timestamp"`
	HostName  string    `json:"hostname"`
	FreeMem   uint64    `json:"freemem"`
	DiskUsage int       `json:"diskusage"`
	CPULoad   int       `json:"cpuload"`
}

// ServerConfig represents the structure of the YAML file
type ServerConfig struct {
	Servers []Server `yaml:"servers"`
}

// Server represents an individual server
type Server struct {
	ServerName string `yaml:"servername"`
	Address    string `yaml:"address"`
}
