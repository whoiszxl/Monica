package core

//简化结构，使用go map + DictEntry, 将map作为一个数组来使用，key为用户输入key经过hash计算并与数组大小（map没有大小，模拟一个固定大小值来模仿Redis）取余，
//YedisObject中Ptr直接指向一个DictEntry，然后DictEntry作为单链表头继续保存有hash碰撞的值
type DictMap map[int]*YedisObject


//HashTable 哈希表
type DictHt struct {
	Table [4]DictEntry //指针数组，用于存储键值对
	Size int //table数组的大小
	SizeMask int //掩码 = size -1
	Used int //table数组已用的元素个数，包含next单链表的数据
}

//hash表中的元素
type DictEntry struct {

	Key *YedisObject //hash表的键名
	Value interface{} //hash表中的值,Redis中的实际是个联合体，可以同时做数据库的键值对，hash，还有失效键值对三种功能，此处简略，只做hash使用
	Next *DictEntry //hash冲突时，此指针指向冲突的元素，形成单链表
}

//字典数据结构，储存一些特殊操作时候用到的特殊字段，Redis源码名称为Dict，此处因为键名冲突所以修改为DictHash
type DictHash struct {

	//Type DictType //对应的特定的操作函数
	PrivData interface{} //字典依赖的数据
	Ht [2]DictHt //hash表，键值对存储的地方
	ReHashIdx int //rehash标识，默认-1，不为-1代表正在rehash，存储值标识hash表ht[0]操作进行到了哪个索引值
	Iterators int //当前运行的迭代器数
}