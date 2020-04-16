package main

import (
	"fmt"
)

func main() {
	s := "foobar阿斯蒂芬"
	fmt.Println(s)
	fmt.Println(&s)
	s = "qweqweqweqweqwe"
	fmt.Println(s)
	fmt.Println(&s)
}
