package main

import (
	"fmt"
	"reflect"
)

func main() {
	s := "zxl"
	i := 1
	c := 1000000000000000
	fmt.Println(reflect.TypeOf(s).String())
	fmt.Println(reflect.TypeOf(i).String())
	fmt.Println(reflect.TypeOf(c).Name())

}
