package data

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	t "github.com/hanymamdouh82/srvmon/internal/types"
	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/process"
)

func GetData(conf *t.ServerConfig) (data t.MonData) {
	data.Timestamp = time.Now()

	hostName(&data)
	mem(&data)
	cpuLoad(&data)
	diskUsage(conf, &data)
	processes(conf, &data)

	// utils.LogJSON(data)
	return data
}

func mem(d *t.MonData) {
	d.FreeMem = memory.FreeMemory()
}

func hostName(d *t.MonData) {
	hostname, _ := os.Hostname()
	d.HostName = hostname
}

func cpuLoad(d *t.MonData) {
	// Get CPU load (1 second average)
	cpuPercentages, _ := cpu.Percent(time.Second, false)
	d.CPULoad = int(cpuPercentages[0]) // Convert float to int for simplicity
}

func diskUsage(conf *t.ServerConfig, d *t.MonData) {
	// Get Disk Usage (%)
	d.DiskUsage = map[string]int{}

	// add root if no disk is added
	if len(conf.Disks) == 0 {
		conf.Disks = append(conf.Disks, "/")
	}

	for _, v := range conf.Disks {
		// to avoid wrong or unmounted partitions
		diskStats, err := disk.Usage(v)
		if err != nil {
			return
		}
		d.DiskUsage[v] = int(diskStats.UsedPercent)
	}
}

func processes(conf *t.ServerConfig, d *t.MonData) {
	d.Processes = []t.ProcessInfo{}

	for _, proc := range conf.Processes {
		pids, err := FindProcessByName(proc)
		if err != nil {
			return
		}

		for _, pid := range pids {
			p, err := process.NewProcess(pid)
			if err != nil {
				return
			}

			// Get process name
			name, _ := p.Name()

			// Get process status
			status, _ := p.Status()

			// Get memory usage
			memInfo, _ := p.MemoryInfo()

			// Get CPU usage
			cpuPercent, _ := p.CPUPercent()

			// Fetch systemd logs if available
			logs := GetSystemdLogs(p.Pid, conf.LogDuration)

			pr := t.ProcessInfo{
				ProcessID:   p.Pid,
				Name:        name,
				Status:      status,
				MemoryUsage: memInfo.RSS,
				CPUUsage:    cpuPercent,
				ProcessLogs: logs,
			}

			d.Processes = append(d.Processes, pr)
		}
	}
}

// GetSystemdLogs fetches logs from systemd using journalctl
func GetSystemdLogs(pid int32, since int) string {
	cmd := exec.Command("journalctl", "--no-pager", fmt.Sprintf("_PID=%d", pid), "--since", fmt.Sprintf("%vmin ago", since))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return string(output)
}

// FindProcessByName returns all PIDs of processes matching the given name
func FindProcessByName(name string) ([]int32, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var pids []int32
	for _, p := range processes {
		pName, err := p.Name()
		if err == nil && strings.Contains(strings.ToLower(pName), strings.ToLower(name)) {
			pids = append(pids, p.Pid)
		}
	}

	if len(pids) == 0 {
		return nil, fmt.Errorf("no process found with name: %s", name)
	}

	return pids, nil
}
