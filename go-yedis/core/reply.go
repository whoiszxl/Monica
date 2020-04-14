package core

import "Monica/go-yedis/encrypt"

//添加回复
func addReply(c *YedisClients, o *YedisObject) {
	c.Reply = o.Ptr.(string)
}

//字符串回复
func addReplyStatus(c *YedisClients, s string) {
	r := encrypt.NewString([]byte(s))
	addReplyString(c, r)
}

//错误回复
func addReplyError(c *YedisClients, s string) {
	r := encrypt.NewError([]byte(s))
	addReplyString(c, r)
}
func addReplyString(c *YedisClients, r *encrypt.Resp) {
	if ret, err := encrypt.EncodeToBytes(r); err == nil {
		c.Reply = string(ret)
	}
}