package utils

import "time"

//获取当前时间的毫秒级时间戳
func CurrentTimeMillis() int{
	return int(time.Now().UnixNano() / 1e6)
}

//获取当前时间的秒级时间戳
func CurrentTimeSecond() int{
	return int(time.Now().UnixNano() / 1e9)
}