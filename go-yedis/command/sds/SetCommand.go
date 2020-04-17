package sds

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strconv"
)

// set命令
// Redis源码中key用的是robj,这里精简一下直接采用go原生string
func SetCommand(c *core.YedisClients, s *core.YedisServer) {
	//校验有效性
	if !checkParam(c) {
		return
	}

	//获取键值对
	robjKey := c.Argv[1]
	robjValue := c.Argv[2]

	//判断是否是字符串，是则设置到Db的Data中
	if stringKey, ok1 := robjKey.Ptr.(string); ok1 {
		//Redis如果是int则直接保存到YedisObject的Ptr中,不需要sds封装
		//简化一下，int类型也直接保存到sds中
		if stringValue, ok3 := robjValue.Ptr.(string); ok3 {

			//判断是否能转int，能转则设置encoding的方式
			if _, err := strconv.Atoi(stringValue); err == nil {
				//创建一个sdshdr保存到字典中
				robjSds := ds.Sdshdr{Len: uint64(len(stringValue)), Free: 0, Buf: stringValue}
				c.Db.Data[stringKey] = core.CreateObject(core.OBJ_STRING, core.OBJ_ENCODING_INT, robjSds)
			}else {
				//创建一个sdshdr保存到字典中
				robjSds := ds.Sdshdr{Len: uint64(len(stringValue)), Free: 0, Buf: stringValue}
				//注意事项：字符串的编码方式在Redis中有三种，首先INT方式已经在上个if判断中添加了，INT编码方式不需要sds对象包装，可以提升效率，它底层实际是个long
				//其次是RAW和EMBSTR, 都是字符串。小于39字节用EMBSTR,大于用RAW，Redis3.2版本则以44字节区分
				//此处省略判断，直接用RAW
				c.Db.Data[stringKey] = core.CreateObject(core.OBJ_STRING, core.OBJ_ENCODING_RAW, robjSds)
			}
		}
	}

	s.Dirty++
	core.AddReplyStatus(c, "OK")
}

//校验参数是否有效
func checkParam(c *core.YedisClients) bool {
	if c.Argc < 3 {
		core.AddReplyError(c, "(error) ERR wrong number of arguments for 'set' command")
		return false
	}
	if c.Argc > 3 {
		core.AddReplyError(c, "(error) ERR syntax error")
		return false
	}
	return true
}