package error

import (
	"log"
	"os"
)

func CheckError(err error) {

	if err != nil {
		log.Println("err ", err.Error())
		os.Exit(1)
	}
}