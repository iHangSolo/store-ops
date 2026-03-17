package resource

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type ResourceInfo struct {
	CPUPercent        float64    `json:"cpu_percent"`
	MemoryPercent     float64    `json:"memory_percent"`
	MemoryAvailableGB float64    `json:"memory_available_gb"`
	Disks             []DiskInfo `json:"disks"`
	RecordedAt        time.Time  `json:"recorded_at"`
}

type DiskInfo struct {
	Name        string  `json:"name"`
	Percent     float64 `json:"percent"`
	TotalGB     float64 `json:"total_gb"`
	AvailableGB float64 `json:"available_gb"`
}

// Collect 收集资源信息
func Collect() (*ResourceInfo, error) {
	info := &ResourceInfo{
		RecordedAt: time.Now(),
	}

	// CPU
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		info.CPUPercent = cpuPercent[0]
	}

	// 内存
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryPercent = memInfo.UsedPercent
		info.MemoryAvailableGB = float64(memInfo.Available) / 1024 / 1024 / 1024
	}

	// 磁盘
	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, p := range partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			info.Disks = append(info.Disks, DiskInfo{
				Name:        p.Mountpoint,
				Percent:     usage.UsedPercent,
				TotalGB:     float64(usage.Total) / 1024 / 1024 / 1024,
				AvailableGB: float64(usage.Free) / 1024 / 1024 / 1024,
			})
		}
	}

	return info, nil
}

// ToJSON 转换为JSON
func (r *ResourceInfo) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

// GetCPUCoreCount 获取CPU核心数
func GetCPUCoreCount() int {
	return runtime.NumCPU()
}