package core

const (
	OBJ_ENCODING_RAW        = 0  /* Raw representation */
	OBJ_ENCODING_INT        = 1  /* Encoded as integer */
	OBJ_ENCODING_HT         = 2  /* Encoded as hash table */
	OBJ_ENCODING_ZIPMAP     = 3  /* Encoded as zipmap */
	OBJ_ENCODING_LINKEDLIST = 4  /* No longer used: old list encoding. */
	OBJ_ENCODING_ZIPLIST    = 5  /* Encoded as ziplist */
	OBJ_ENCODING_INTSET     = 6  /* Encoded as intset */
	OBJ_ENCODING_SKIPLIST   = 7  /* Encoded as skiplist */
	OBJ_ENCODING_EMBSTR     = 8  /* Embedded sds string encoding */
	OBJ_ENCODING_QUICKLIST  = 9  /* Encoded as linked list of ziplists */
	OBJ_ENCODING_STREAM     = 10 /* Encoded as a radix tree of listpacks */

	OBJ_STRING = 0 /* 字符串对象*/
	OBJ_LIST   = 1 /* 列表对象 */
	OBJ_SET    = 2 /* 集合对象 */
	OBJ_ZSET   = 3 /* 有序集合对象 */
	OBJ_HASH   = 4 /* 哈希表对象 */
	OBJ_MODULE = 5 /* 模块对象 */

	ENABLE = 1
	DISABLE = 0

	//TODO info命令使用的换行分隔符，因为\r\n在redis协议加解密的时候会截取掉，所以暂用特殊方案解决
	INFO_LINE_SEPARATOR = "$"
)
