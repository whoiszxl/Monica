package main

import "fmt"

func main() {

	//times 33 算法
	key := Key("wwwwwww")

	fmt.Println(key)
}

func Key(s string) int64 {
	var hash int64  = 5381

	for _, c := range s {
		hash = ((hash << 5) + hash) + int64(c)
	}

	return hash
}
