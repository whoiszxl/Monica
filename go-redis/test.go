package main

import "regexp"
import "fmt"
import "strings"
import "bytes"
import "Monica/go-redis/proto"


func main() {
	test3()
}

func test3() {
	p, e := proto.EncodeCmd("set name www")
	if e != nil {
		fmt.Println(e)
	}

	decoder := proto.NewDecoder(bytes.NewReader([]byte(p)))
	if resp, err := decoder.DecodeMultiBulk(); err == nil {
		for k, s := range resp {
			fmt.Println(k,string(s.Value))
		}
	}
}

func test2() {

	pro := ""

	ret := strings.Split("set name www", " ")

	for k, v := range ret {
		if k == 0 {
            pro = fmt.Sprintf("*%d\r\n", len(ret))
        }
        pro += fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
	}

	fmt.Println("%s", pro)

}

func test1() {
	r := regexp.MustCompile("[^\\s]+")
	parts := r.FindAllString(strings.Trim("set name www", " "), -1)
	argc, argv := len(parts), parts
	//c.Argv = make([]*object.GodisObject, 5)
	fmt.Println(argc)
	j := 0
	for _, v := range argv {
		fmt.Println(v)
		j++
	}
}