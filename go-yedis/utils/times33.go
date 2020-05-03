package utils

//times33算法获取字符串数字hash值
func Times33Encoding(str string) int64 {
	var hash int64  = 5381
	for _, c := range str {
		hash = ((hash << 5) + hash) + int64(c)
	}
	return hash
}
