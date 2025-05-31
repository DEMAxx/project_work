package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := new(ListItem)
	item.Value = v
	item.Next = l.front

	if l.len == 0 {
		l.back = item
	} else {
		l.front.Prev = item
	}

	l.front = item
	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := new(ListItem)
	item.Value = v

	if l.len == 0 {
		l.front = item
	} else {
		l.back.Next = item
	}

	item.Prev = l.back
	l.back = item
	l.len++

	return item
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		panic("ListItem is nil")
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.front = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.back = item.Prev
	}

	item.Next = nil
	item.Prev = nil
	l.len--
}

func (l *list) MoveToFront(item *ListItem) {
	exNext := item.Next
	exPrev := item.Prev

	if exPrev == nil {
		return
	}

	if exNext == nil {
		exPrev.Next = nil
		l.back = exPrev
	} else {
		exPrev.Next = item.Next
		exNext.Prev = item.Prev
	}

	exFront := l.Front()
	exFront.Prev = item
	item.Next = exFront
	item.Prev = nil
	l.front = item
}

func NewList() List {
	return new(list)
}
