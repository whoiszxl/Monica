package core

import (
	"Monica/go-yedis/encrypt"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


//将命令追加到aof文件中
func FeedAppendOnlyFile(server *YedisServer, client *YedisClients) {

	buf := ""

	//刚开始执行的时候需要SELECT数据库
	if client.Db.ID != server.AofSelectedDb {
		buf = fmt.Sprintf("*2\r\n$6\r\nSELECT\r\n$%lu\r\n%s\r\n", len(strconv.Itoa(int(client.Db.ID))), client.Db.ID)
		server.AofSelectedDb = client.Db.ID
	}

	//TODO 需要将EXPIRE 、 PEXPIRE 和 EXPIREAT 命令转换成PEXPIREAT
	//TODO 需要将SETEX 和 PSETEX 命令转换成 SET 和 PEXPIREAT

	//Redis没有QueryBuf这个字段，所以每次都需要把命令重新转回协议格式，Yedis直接将协议格式的命令缓存在client的QueryBuf字段中，便可以直接追加了
	//将命令和参数转换为Redis协议格式，再添加到需要写入的buf中
	//commandStr := catAppendOnlyGenericCommand(client.Argc, client.Argv)

	//判断是否需要aof，开启了则将命令写入server的aofBuff缓冲区
	if server.AofState == REDIS_AOF_ON {
		server.AofBuf = server.AofBuf + buf + client.QueryBuf
	}

}

//将执行命令转换为Redis协议格式
func catAppendOnlyGenericCommand(argc int,argv []*YedisObject) string {

	//将argv中的命令和参数转换为输入时的字符串
	var cmd string
	for _, v := range argv {
		cmd += v.Ptr.(string) + " "
	}
	cmd = strings.TrimRight(cmd, ",")
	encodeCmd, e := encrypt.EncodeCmd(cmd)
	if e != nil {
		return ""
	}
	return string(encodeCmd)
}

//AppendToFile 写文件
func AppendToFile(f *os.File, content string) (*os.File, error) {
	n, _ := f.Seek(0, os.SEEK_END)
	_, err := f.WriteAt([]byte(content), n)
	return f, err
}

//读取aof文件
func ReadAof(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("aof file open failed" + err.Error())
	}
	defer f.Close()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("aof file read failed" + err.Error())
	}
	ret := bytes.Split(content, []byte{'*'})
	var pros = make([]string, len(ret)-1)
	for k, v := range ret[1:] {
		v := append(v[:0], append([]byte{'*'}, v[0:]...)...)
		pros[k] = string(v)
	}
	return pros
}

//
func RewriteAppendOnlyFileBackground() {

}