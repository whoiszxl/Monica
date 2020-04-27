package core


// 从表头向表尾进行迭代
const AL_START_HEAD = 0

// 从表尾到表头进行迭代
const AL_START_TAIL = 1

//双向链表节点结构体
type ListNode struct {

	//前置节点
	Prev *ListNode

	//后置节点
	Next *ListNode

	//节点值
	Value interface{}
}

//双向列表迭代器
type ListIter struct {

	//当前迭代的节点
	Next *ListNode

	//迭代方向
	Direction int
}


//双向链表结构体
type LinkedList struct {
	//表头节点
	Head *ListNode

	//表尾节点
	Tail *ListNode

	//双向链表长度
	Len int
}

//创建一个双向链表
func ListCreate() *LinkedList {

	//创建对象
	list := new(LinkedList)

	//分配属性
	list.Head, list.Tail = nil, nil
	list.Len = 0
	return list
}

//将一个任意类型的对象添加到链表表头
func ListAddNodeHead(list *LinkedList, value interface{}) *LinkedList {
	//创建节点并将值传入节点
	node := new(ListNode)
	node.Value = value

	if list.Len == 0 {
		//链表为空,设置节点的前后置节点都为空，链表的头尾都是node本身
		node.Prev, node.Next = nil, nil
		list.Head, list.Tail = node, node
	}else {
		//链表不为空，则添加到表头，表头节点没有上一个，则赋值为nil，表头节点的下个节点则为原链表的头节点
		//将原链表头结点的前置节点设置为当前节点,并将链表头结点设置为当前节点
		node.Prev = nil
		node.Next = list.Head
		list.Head.Prev = node
		list.Head = node
	}

	//链表长度累加并返回链表
	list.Len++
	return list
}

//将一个任意类型的对象添加到链表表尾
func ListAddNodeTail(list *LinkedList, value interface{}) *LinkedList {
	//创建节点并将值传入节点
	node := new(ListNode)
	node.Value = value

	if list.Len == 0 {
		//链表为空,设置节点的前后置节点都为空，链表的头尾都是node本身
		node.Prev, node.Next = nil, nil
		list.Head, list.Tail = node, node
	}else {
		//链表不为空，则添加到表尾，当前节点和链表建立关联，将当前节点的上一个节点指向原链表的尾部节点，并将当前节点的下一个节点指向nil
		//将原链表的尾部节点的下一个节点指向当前节点，然后将当前链表的尾部指向当前节点
		node.Prev = list.Tail
		node.Next = nil
		list.Tail.Next = node
		list.Tail = node
	}

	//链表长度累加并返回链表
	list.Len++
	return list
}


//将value对象插入到oldNode的前或后
//after: 1->插入到之后 0->插入到之前
func ListInsertNode(list *LinkedList, oldNode *ListNode, value interface{}, after int) *LinkedList {
	//创建节点并将值传入节点
	node := new(ListNode)
	node.Value = value

	if after == 1 {
		//将节点插入到oldNode节点之后
		node.Prev = oldNode
		node.Next = oldNode.Next

		//如果旧节点是尾部，则传入的新节点需要成为新尾部
		if list.Tail == oldNode {
			list.Tail = node
		}
	}else {
		//将节点插入到oldNode节点之前
		node.Next = oldNode
		node.Prev = oldNode.Prev
		//如果旧节点是头部，则传入的新节点需要成为新头部
		if list.Head == oldNode {
			list.Head = node
		}
	}

	//更新后新节点的前后置两个节点还没指向到新节点，需要重新指向
	if node.Prev != nil {
		node.Prev.Next = node
	}
	if node.Next != nil {
		node.Next.Prev = node
	}

	//链表长度累加并返回链表
	list.Len++
	return list
}

//删除指定节点
func ListDelNode(list *LinkedList, node *ListNode) {
	//如果要删除的节点的前置节点存在，则将前置节点的下一个节点指针指向要删除节点的下一个，直接跳过
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}else {
		//要删除节点的前置节点不存在，则其是头部节点，直接头部指向要删除节点的下一个节点
		list.Head = node.Next
	}

	//如果要删除的节点后置节点存在
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}else {
		list.Tail = node.Prev
	}

	//Go不需要和C一样释放内存 node == nil

	//长度-1
	list.Len--
}

//给链表创建一个迭代器
//direction为迭代方向： AL_START_HEAD->从头到尾 AL_START_TAIL->从尾到头
func ListGetIterator(list *LinkedList, direction int) *ListIter {
	iter := new(ListIter)

	if direction == AL_START_HEAD {
		iter.Next = list.Head
	}else if direction == AL_START_TAIL {
		iter.Next = list.Tail
	}

	//记录迭代方向
	iter.Direction = direction
	return iter
}

//设置迭代器的方向为AL_START_HEAD
func ListRewind(list *LinkedList, li *ListIter) {
	li.Next = list.Head
	li.Direction = AL_START_HEAD
}


//设置迭代器的方向为AL_START_TAIL
func ListRewindTail(list *LinkedList, li *ListIter) {
	li.Next = list.Tail
	li.Direction = AL_START_TAIL
}

//获取迭代器当前指向的节点，并将指针移动一位
func ListNext(iter *ListIter) *ListNode {
	current := iter.Next
	if current != nil {
		//根据方向指向下一个节点
		if iter.Direction == AL_START_HEAD {
			iter.Next = current.Next
		}else if iter.Direction == AL_START_TAIL {
			iter.Next = current.Prev
		}
	}
	return current
}


//查找链表中和key匹配的节点
func ListSearchKey(list *LinkedList, key interface{}) *ListNode {

	iter := ListGetIterator(list, AL_START_HEAD)

	for {
		node := ListNext(iter)
		if node == nil {
			break
		}
		if key == node.Value {
			return node
		}
	}
	return nil
}

//返回索引上的值
func ListIndex(list *LinkedList, index int) *ListNode {
	var n *ListNode
	if index < 0 {
		index = (-index)-1
		n = list.Tail
		for index != 0 && n != nil{
			n = n.Prev
			index--
		}
	}else {
		n = list.Head
		for index != 0 && n != nil{
			n = n.Next
			index--
		}
	}
	return n
}


// 处于阻塞状态的list
type ReadyList struct {
	Db *YedisDb
	Key *YedisObject
}