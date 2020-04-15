package utils

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

func ReadConfig(configPath string) (NetConfig, DbConfig, AofConfig) {
	//加载基础配置文件
	cfg, err := goconfig.LoadConfigFile(configPath)
	if err != nil {
		panic("读取配置文件错误")
	}
	//获取网络相关配置
	netBind, err := cfg.GetValue("net", "bind")
	errVerify(err)
	netPort, err := cfg.GetValue("net", "port")
	errVerify(err)

	dbDatabases, err := cfg.Int("db", "databases")
	errVerify(err)
	dbDbfilename, err := cfg.GetValue("db", "dbfilename")
	errVerify(err)
	dbSavetime, err := cfg.Int("db", "savetime")
	errVerify(err)
	dbSavenumber, err := cfg.Int("db", "savenumber")
	errVerify(err)
	dbRequirepass, err := cfg.GetValue("db", "requirepass")
	errVerify(err)
	dbHz, err := cfg.Int("db", "hz")
	errVerify(err)

	aofAppendonly, err := cfg.GetValue("aof", "appendonly")
	errVerify(err)
	aofAppendfilename, err := cfg.GetValue("aof", "appendfilename")
	errVerify(err)
	aofAppendfsync, err := cfg.GetValue("aof", "appendfsync")
	errVerify(err)

	netConfig := NetConfig{netBind, netPort, netBind + ":" + string(netPort)}
	dbConfig := DbConfig{dbDbfilename, dbDatabases, dbSavetime, dbSavenumber, dbRequirepass, dbHz}
	aofConfig := AofConfig{aofAppendonly, aofAppendfilename, aofAppendfsync}

	return netConfig, dbConfig, aofConfig
}

//校验error
func errVerify(err error) {
	if err != nil {
		log.Println("[error]read yedis.conf", err)
		os.Exit(1)
	}
}

//网络配置结构体
type NetConfig struct {
	NetBind string
	NetPort string
	NetHost string
}

//数据库配置结构体
type DbConfig struct {
	DbDbfilename  string
	DbDatabases   int
	DbSavetime    int
	DbSavenumber  int
	DbRequirepass string
	Hz int
}

//AOF配置结构体
type AofConfig struct {
	AofAppendonly     string
	AofAppendfilename string
	AofAppendfsync    string
}
