package data_api

import (
	"blogx_server/common/res"
	computer "blogx_server/utils/conputer"

	"github.com/gin-gonic/gin"
)

type ComputerInfo struct {
	CpuPercent  float64 `json:"cpuPercent"`
	MemPercent  float64 `json:"memPercent"`
	DiskPercent float64 `json:"diskPercent"`
}

func (DataApi) CpmputerInfoView(c *gin.Context) {
	cpuPercent := computer.GetCpuPercent()
	memPercent := computer.GetMemPercent()
	diskPercent := computer.GetDiskPercent()
	res.SuccessWithData(ComputerInfo{
		CpuPercent:  cpuPercent,
		MemPercent:  memPercent,
		DiskPercent: diskPercent,
	}, c)
}
