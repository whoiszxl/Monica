package config

import (
    "strings"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "regexp"
)

var jsonData map[string]interface{}

func initJSON() {
    bytes, err := ioutil.ReadFile("./zedis.json")
    if err != nil {
        fmt.Println("read zedis.json error: ", err.Error())
    }

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

    if err := json.Unmarshal(bytes, &jsonData); err != nil {
        fmt.Println("invalid config: ", err.Error())
    }
}