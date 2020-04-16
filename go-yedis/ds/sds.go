package ds

// sds简单动态字符串数据结构
// Redis在3.2版本后sds结构有五种，分别为sdshdr5、sdshdr8、sdshdr16、sdshdr32和sdshdr64
// Redis3.0 sds只有一种数据结构，为sdshdr
// 为了简化所以采用Redis3.0的sdshdr结构
// 实际上在Go中，使用原生的string就好了，string底层自带了len属性也是O(1)复杂度
// 也可以用rune，rune是int32的别名，对应utf-8的字符数字编码
type Sdshdr struct {

	//buf数组中已经使用字节的数量，相当于当前字符串的长度
	//len的存在是因为Redis的C底层字符数组不支持O(1)的复杂度查询，为O(N),为了优化STRLEN()方法的效率而引入Len属性
	//在Go中实际也是多余的，Go中len(Buf)调用本身就是O(1)的复杂度，为了模仿一下Redis而加这个参数了
	Len uint64

	//buf数组中未使用字节的数量
	Free uint

	//字节数组，用来保存字符串，Redis是字符数组
	//精简一下了，看源码还需要dbAdd->dictAdd->dictSetVal,过于复杂了
	Buf string
}