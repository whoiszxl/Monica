package ds

//双向链表节点结构体
type ListNode struct {

	//前置节点
	Prev *ListNode

	//后置节点
	Next *ListNode

	//节点值
	value interface{}
}

//双向列表迭代器
type ListIter struct {

	//当前迭代的节点
	next *ListNode

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
	len int
}