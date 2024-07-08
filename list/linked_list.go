package list

import (
	"github.com/ac-zht/gotools"
)

type node[T any] struct {
	prev *node[T]
	next *node[T]
	val  T
}

type LinkedList[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

func NewLinkedList[T any]() *LinkedList[T] {
	head := &node[T]{}
	tail := &node[T]{prev: head, next: head}
	head.prev, head.next = tail, tail
	return &LinkedList[T]{
		head: head,
		tail: tail,
	}
}

func NewLinkedListOf[T any](src []T) *LinkedList[T] {
	list := NewLinkedList[T]()
	if err := list.Append(src...); err != nil {
		panic(err)
	}
	return list
}

func (l *LinkedList[T]) findNode(index int) *node[T] {
	var current *node[T]
	if index < l.length/2 {
		current = l.head
		for i := -1; i < index; i++ {
			current = current.next
		}
	} else {
		current = l.tail
		for i := l.Len(); i > index; i-- {
			current = current.prev
		}
	}
	return current
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	if !l.checkIndex(index) {
		var zeroValue T
		return zeroValue, gotools.NewErrIndexOutOfRange(l.length, index)
	}
	current := l.findNode(index)
	return current.val, nil
}

func (l *LinkedList[T]) checkIndex(index int) bool {
	return index >= 0 && index < l.length
}

// Add 指定下标后插入新元素
func (l *LinkedList[T]) Add(index int, val T) error {
	if !l.checkIndex(index) {
		return gotools.NewErrIndexOutOfRange(l.length, index)
	}
	current := l.findNode(index)
	next := current.next
	current.next = &node[T]{val: val}
	current.next.next = next
	l.length++
	return nil
}

func (l *LinkedList[T]) Append(values ...T) error {
	for _, v := range values {
		node := &node[T]{prev: l.tail.prev, next: l.tail, val: v}
		node.prev.next, node.next.prev = node, node
		l.length++
	}
	return nil
}

func (l *LinkedList[T]) Delete(index int) (T, error) {
	if !l.checkIndex(index) {
		var zeroValue T
		return zeroValue, gotools.NewErrIndexOutOfRange(l.length, index)
	}
	current := l.findNode(index)
	prev := current.prev
	next := current.next
	prev.next, next.prev = next, prev
	current.prev, current.next = nil, nil
	l.length--
	return current.val, nil
}

func (l *LinkedList[T]) Set(index int, val T) error {
	if !l.checkIndex(index) {
		return gotools.NewErrIndexOutOfRange(l.length, index)
	}
	current := l.findNode(index)
	current.val = val
	return nil
}

func (l *LinkedList[T]) Range(fn func(i int, val T) error) error {
	for i, current := 0, l.head.next; i < l.length; i++ {
		err := fn(i, current.val)
		if err != nil {
			return err
		}
		current = current.next
	}
	return nil
}

func (l *LinkedList[T]) Len() int {
	return l.length
}

func (l *LinkedList[T]) Cap() int {
	return l.Len()
}

func (l *LinkedList[T]) AsSlice() []T {
	current := l.head.next
	slice := make([]T, l.length)
	for i := 0; i < l.length; i++ {
		slice[i] = current.val
		current = current.next
	}
	return slice
}
