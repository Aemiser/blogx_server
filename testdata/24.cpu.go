package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}
func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}
func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}
func main() {
	fmt.Printf("CPU 使用率: %.2f%%\n", GetCpuPercent())
	fmt.Printf("MEM 使用率: %.2f%%\n", GetMemPercent())
	fmt.Printf("DISK 使用率: %.2f%%\n", GetDiskPercent())
}
