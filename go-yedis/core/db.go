package core

type YedisDb struct {
	Dict dict
	Expires dict
	ID int32
}

//使用Go原生数据结构map作为redis中dict结构体
type dict map[string]*YedisObject