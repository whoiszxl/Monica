package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().UnixNano() / 1000000)
}
