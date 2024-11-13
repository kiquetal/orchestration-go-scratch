package worker

import "github.com/c9s/goprocinfo/linux"

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
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
