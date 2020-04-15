package command

import (
	"Monica/go-yedis/core"
	"bytes"
	"fmt"
	"strconv"
	"time"
)


//info命令
func InfoCommand(c *core.YedisClients, s *core.YedisServer) {

	var result string
	//server信息
	pid := "pid:" + strconv.Itoa(s.Pid)
	addr := "addr:" + s.BindAddr
	port := "port:" + s.Port
	configFile := "config_file:" + s.ConfigFile
	result = appendStr("[server]", pid, addr, port, configFile)

	//db信息
	dbNum := "dbNum:" + strconv.Itoa(s.DbNum)
	hz := "hz:" + strconv.Itoa(s.Hz)
	result = appendStr(result, "[db]", dbNum, hz)

	//持久化信息 persistence
	rdbFileName := "rdb_file_name:" + s.RdbFileName
	aofEnabled := "aof_enabled:" + strconv.Itoa(s.AofEnabled)
	aofFileName := "aof_file_name:" + s.AofFileName
	aofCurrentSize := "aof_current_size:" + strconv.Itoa(s.AofCurrentSize)
	aofSync := "aof_sync:" + s.AofSync
	result = appendStr(result, "[persistence]", rdbFileName, aofEnabled, aofFileName, aofCurrentSize, aofSync)

	//仅用于统计使用的字段，仅取部分 Stats
	statStartTime := "stat_start_time:" + time.Unix(s.StatStartTime / 1000, 0).Format("2006-01-02 15:04:05")
	statNumCommands := "stat_num_commands:" + strconv.Itoa(len(s.Commands))
	statNumConnections := "stat_num_connections:" + strconv.Itoa(int(s.StatNumConnections))
	result = appendStr(result, "[stats]", statStartTime, statNumCommands, statNumConnections)

	//系统硬件信息 system
	systemAllMemorySize := "system_all_memory_size:" + strconv.FormatUint(s.SystemAllMemorySize, 10)
	systemAvailableSize := "systemAvailableSize:" + strconv.FormatUint(s.SystemAvailableSize, 10)
	systemUsedSize := "systemUsedSize:" + strconv.FormatUint(s.SystemUsedSize, 10)
	systemUsedPercent := "systemUsedPercent:" + strconv.FormatFloat(s.SystemUsedPercent,'E',-1,64) + "%"
	systemCpuPercent := "systemCpuPercent:" + strconv.FormatFloat(s.SystemCpuPercent,'E',-1,64) + "%"
	result = appendStr(result, "[system]", systemAllMemorySize, systemAvailableSize, systemUsedSize, systemUsedPercent, systemCpuPercent)

	//keySpace 有效key，失效key和平均ttl
	dbCode := c.Db.ID
	keys := len(c.Db.Data)
	expiredKeys := len(c.Db.Expires)
	avgTTL := c.Db.AvgTTL

	keySpace := fmt.Sprintf("db%d:keys=%d,expires=%d,avg_ttl=%d", dbCode, keys, expiredKeys, avgTTL)
	result = appendStr(result, "[keySpace]", keySpace)
	core.AddReplyStatus(c, result)
}

//多字符串分隔符拼接
func appendStr(args ...string) string {
	var buffer bytes.Buffer
	for _, str := range args {
		buffer.WriteString(str)
		buffer.WriteString(core.INFO_LINE_SEPARATOR)
	}
	return buffer.String()
}