package core

type YedisDb struct {
	Data Dict
	Expires Dict
	ID int32
}

//使用Go原生数据结构map作为redis中dict结构体
type Dict map[string]*YedisObject