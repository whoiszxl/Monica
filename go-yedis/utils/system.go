package utils

import "time"

//获取当前时间的毫秒级时间戳
func CurrentTimeMillis() int{
	return int(time.Now().UnixNano() / 1e6)
}

//获取当前时间的微秒时间戳
func CurrentTimeMicrosecond() int {
	return int(time.Now().UnixNano() / 1e3)
}

//获取当前时间的秒级时间戳
func CurrentTimeSecond() int{
	return int(time.Now().UnixNano() / 1e9)
}

//获取当前时间的秒和毫秒,格式为： sec:1587452770 ms:770
func CurrentSecondAndMillis() (int, int) {
	time := int(time.Now().UnixNano() / 1e6)
	sec := time / 1e3
	ms := time % 1000
	return sec , ms
}