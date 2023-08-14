package list

import "errors"

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

func (l *LinkedList[T]) Get(index int) (T, error) {
	if !l.checkIndex(index) {
		return nil, errors.New("out of range")
	}
	current := l.head
	cnt := 0
	for current != nil {
		if cnt == index {
			return current.val, nil
		}
		cnt++
		current = current.next
	}
	return nil, errors.New("out of range")
}

func (l *LinkedList[T]) checkIndex(index int) bool {
	return index >= 0 && index < l.length
}

func (l *LinkedList[T]) Append(values ...T) {
	for _, v := range values {
		l.tail.next = &node[T]{prev: l.tail, val: v}
		l.tail = l.tail.next
	}
}

func (l *LinkedList[T]) Delete(index int) (T, error) {
	if !l.checkIndex(index) {
		return nil, errors.New("out of range")
	}
}
