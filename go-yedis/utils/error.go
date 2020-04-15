package utils

import (
	"fmt"
	"os"
)

//error校验
func ErrorVerify(message string, err error, isExit bool) {
	if err != nil {
		fmt.Println(message, err)
		if isExit {
			os.Exit(1)
		}
	}
}