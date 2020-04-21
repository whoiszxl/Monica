package utils

import "github.com/shirou/gopsutil/mem"

//获取当前机器已用内存大小，单位为b
func GetUsedMemory() uint64 {
	memInfo, err := mem.VirtualMemory()
	ErrorVerify("获取机器内存信息失败", err, false)
	return memInfo.Used
}