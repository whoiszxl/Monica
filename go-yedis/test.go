package main

import (
	"Monica/go-yedis/ds"
	"fmt"
)

func main() {

	//1. 创建LinkedList
	list := ds.ListCreate()
	fmt.Println("创建LinkedList", list)

	//2. 将一个任意类型的对象添加到链表表头
	ds.ListAddNodeHead(list, "xieanqi")
	//3. //将一个任意类型的对象添加到链表表尾
	ds.ListAddNodeTail(list, "zhangguorong")

	key := ds.ListSearchKey(list, "xieanqi")

	ds.ListInsertNode(list, key, "whoiszxl", 1)

	ds.ListDelNode(list, key)

	index := ds.ListIndex(list, 1)
	fmt.Println("index", index)

	//迭代打印
	iter := ds.ListGetIterator(list, ds.AL_START_HEAD)
	for {
		node := ds.ListNext(iter)
		if node == nil {
			break
		}else {
			fmt.Println(node)
		}

	}
}
