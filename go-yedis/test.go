package main

import (
	"fmt"
	"reflect"
)

func main() {
	s := "zxl"
	i := 1
	c := 1000000000000000
	fmt.Println(reflect.TypeOf(s))
	fmt.Println(reflect.TypeOf(i))
	fmt.Println(reflect.TypeOf(c))

}
