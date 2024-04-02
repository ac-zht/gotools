package queue

import (
	"errors"
)

var (
	ErrOutOfCapacity = errors.New("超出最大容量限制")
	ErrEmptyQueue    = errors.New("队列为空")
)

type Comparer[T any] func(x, y T) int

type priorityQueue[T any] struct {
	//队列容量，小于0表示无界队列，大于0表示有界队列
	capacity int
	compare  Comparer[T]
	zero     T
	//存储切片，留下标0作为哨兵节点，便于节点计算
	data []T
}

func (p *priorityQueue[T]) Enqueue(val T) error {
	if p.isFull() {
		return ErrOutOfCapacity
	}
	p.data = append(p.data, val)
	node, parent := len(p.data)-1, (len(p.data)-1)/2
	//上滤操作：与父节点比较符合要求交换位置
	for node > 1 && p.compare(p.data[node], p.data[parent]) > 0 {
		p.data[node], p.data[parent] = p.data[parent], p.data[node]
		node = parent
		parent = node / 2
	}
	return nil
}

func (p *priorityQueue[T]) Dequeue() (val T, err error) {
	if p.isEmpty() {
		return p.zero, ErrEmptyQueue
	}
	val = p.data[1]
	last := len(p.data) - 1
	p.data[1] = p.data[last]
	p.data = p.data[:last]
	i, j := 1, 1
	for {
		if left := 2 * i; left < len(p.data) && p.compare(p.data[left], p.data[i]) > 0 {
			i = left
		}
		if right := (2 * i) + 1; right < len(p.data) && p.compare(p.data[right], p.data[i]) > 0 {
			i = right
		}
		if j == i {
			break
		}
		p.data[j], p.data[i] = p.data[i], p.data[j]
		j = i
	}
	return val, nil
}

func (p *priorityQueue[T]) peek() (val T, err error) {
	if p.isEmpty() {
		return p.zero, ErrEmptyQueue
	}
	return p.data[1], nil
}

func (p *priorityQueue[T]) isFull() bool {
	return p.capacity > 0 && len(p.data)-1 == p.capacity
}

func (p *priorityQueue[T]) isEmpty() bool {
	return len(p.data) < 2
}

func (p *priorityQueue[T]) Len() int {
	return len(p.data) - 1
}

func NewPriorityQueue[T any](capacity int, compare Comparer[T]) *priorityQueue[T] {
	sliceCap := capacity + 1
	if capacity <= 0 {
		capacity = 0
		sliceCap = 64
	}
	return &priorityQueue[T]{
		capacity: capacity,
		compare:  compare,
		data:     make([]T, 1, sliceCap),
	}
}
