package worker

import (
	"github.com/c9s/goprocinfo/linux"
	"log"
)

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
	TaskCount int
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) MemFreeKb() uint64 {
	return s.MemStats.MemFree
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemFree
}

func (s *Stats) MemUsedPercent() float64 {
	return float64(s.MemUsedKb()) / float64(s.MemTotalKb()) * 100
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsedPercent() float64 {
	return float64(s.DiskUsed()) / float64(s.DiskTotal()) * 100
}

func (s *Stats) CpuUsage() float64 {
	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle

	if total == 0 {
		return 0.00
	}
	return float64(total) - float64(idle)/float64(total)
}

func GetMemoryInfo() *linux.MemInfo {
	memstats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("Failed to read memory from /proc/meminfo: %v", err)
		return &linux.MemInfo{}
	}
	return memstats
}

func GetDiskInfo() *linux.Disk {
	diskstats, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("Failed to read disk from /: %v", err)
		return &linux.Disk{}
	}
	return diskstats
}

func GetCpuInfo() *linux.CPUStat {
	cpustats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Failed to read cpu from /proc/stat: %v", err)
		return &linux.CPUStat{}
	}
	return &cpustats.CPUStatAll
}

func GetLoadAvg() *linux.LoadAvg {
	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("Failed to read loadavg from /proc/loadavg: %v", err)
		return &linux.LoadAvg{}
	}
	return loadavg
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuInfo(),
		LoadStats: GetLoadAvg(),
	}
}
