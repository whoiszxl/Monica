package utils

import "fmt"

//error校验
func ErrorVerify(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
	}
}