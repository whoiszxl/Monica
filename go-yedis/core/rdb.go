package core

import (
	"Monica/go-yedis/utils"
	"os"
)

//bgasve
func RdbSaveBackground(fileName string) {

}

//将rdb文件中的数据加载到内存中
func RdbLoad(fileName string) int {
	return REDIS_OK
}

//标记程序正在载入中
func StartLoading(s *YedisServer, f *os.File) {
	//标记服务器正在载入中，并且记录开始载入时间，纳秒
	s.Loading = 1
	s.LoadingStartTime = utils.CurrentTimeNano()

	//记录载入文件的大小
	info, e := f.Stat()
	if e != nil {
		return
	}

	s.LoadingTotalBytes = int(info.Size())
}
