package core

import (
	"Monica/go-yedis/encrypt"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)


//将命令追加到aof文件中
func FeedAppendOnlyFile(server *YedisServer, client *YedisClients) {

	buf := ""

	//刚开始执行的时候需要SELECT数据库
	buf = "*2\r\n$6\r\nselect\r\n$"+ strconv.Itoa(len(strconv.Itoa(int(client.Db.ID)))) +"\r\n" + strconv.Itoa(int(client.Db.ID)) + "\r\n"
	server.AofSelectedDb = client.Db.ID

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
		cmd += v.Ptr.(Sdshdr).Buf + " "
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

//将aof文件加载到内存中
func LoadAppendOnlyFile(s *YedisServer) int {
	//创建一个假客户端来专门执行AOF文件中的语句
	fakeClient := createFakeClient(s)

	//打开AOF文件，并检查是否有效
	f, err := os.OpenFile(s.AofFileName, os.O_WRONLY|syscall.O_CREAT, 0644)
	if err != nil {
		s.AofCurrentSize = 0
		fmt.Println("open aof file fail")
		return REDIS_ERR
	}
	//暂时关闭aof，防止执行multi时exec命令被传播到了现在在读取的aof文件中
	s.AofState = REDIS_AOF_OFF

	//设置服务器的状态为正在载入中,Redis的这个设计是将此个aof和rdb共用的函数放在rdb.c文件里，Redis的不分包的代码属实难受
	StartLoading(s, f)

	//Redis这个时候还要处理PUBSUB模块，先省了省了... 代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/aof.c#L887
	//以下部分先简化一下，不清楚的可以看Redis代码

	aofData := ReadAof(s.AofFileName)

	for _, v := range aofData {
		fakeClient.QueryBuf = string(v)
		err := fakeClient.ProcessInputBuffer()
		if err != nil {
			log.Println("LoadAppendOnlyFile fail", err)
		}
		s.ProcessCommand(fakeClient)
	}
	s.AofState = REDIS_AOF_ON
	return REDIS_OK
}

//创建一个伪客户端
//Redis代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/aof.c#L771
func createFakeClient(s *YedisServer) *YedisClients {
	fakeClient := new(YedisClients)
	SelectDb(fakeClient, s, 0)
	return fakeClient
}