package main

import "regexp"
import "fmt"
import "strings"
import "bytes"
import "Monica/go-redis/proto"
import "io/ioutil"
import "os"


func main() {
	test4()
}

func test4() {
	f, err := os.Open("F:\\go_dev\\src\\Monica\\go-redis\\zedis.aof")
	if err != nil {
		fmt.Println("aof file open failed" + err.Error())
	}
	defer f.Close()
	content, err := ioutil.ReadFile("F:\\go_dev\\src\\Monica\\go-redis\\zedis.aof")
	if err != nil {
		fmt.Println("aof file read failed" + err.Error())
	}
	ret := bytes.Split(content, []byte{'*'})
	var pros = make([]string, len(ret)-1)
	for k, v := range ret[1:] {
		v := append(v[:0], append([]byte{'*'}, v[0:]...)...)
		pros[k] = string(v)
		fmt.Println(pros[k])
	}
	
}

func test3() {
	p, e := proto.EncodeCmd("set name www")
	if e != nil {
		fmt.Println(e)
	}

	decoder := proto.NewDecoder(bytes.NewReader([]byte(p)))
	fmt.Println(decoder.DecodeMultiBulk())
	// if resp, err := decoder.DecodeMultiBulk(); err == nil {
	// 	for k, s := range resp {
	// 		fmt.Println(k,string(s.Value))
	// 	}
	// }
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