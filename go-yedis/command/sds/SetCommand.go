package sds

import (
	"Monica/go-yedis/core"
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
	if robjKey.ObjectType == core.REDIS_STRING && robjValue.ObjectType == core.REDIS_STRING {
		//Redis如果是int则直接保存到YedisObject的Ptr中,不需要sds封装
		//简化一下，int类型也直接保存到sds中
		//查询是否存在
		isExist := core.GetKeyObj(c.Db.Dict, robjKey)
		if isExist != nil {
			robjKey = isExist
		}

		//判断設置的值是否能转int，能转则设置encoding的方式
		if _, err := strconv.Atoi(robjValue.Ptr.(core.Sdshdr).Buf); err == nil {
			robjValue.Encoding = core.OBJ_ENCODING_INT
		}
		//注意事项：字符串的编码方式在Redis中有三种，首先INT方式已经在上个if判断中添加了，INT编码方式不需要sds对象包装，可以提升效率，它底层实际是个long
		//其次是RAW和EMBSTR, 都是字符串。小于39字节用EMBSTR,大于用RAW，Redis3.2版本则以44字节区分
		//此处省略判断，直接用RAW
		c.Db.Dict[robjKey] = robjValue
	}

	//每次进行增删改的时候自增1，通过这个自增来判断是否需要添加到aof缓存中
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
