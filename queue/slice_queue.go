package queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

// SliceQueue 切片实现并发队列
type SliceQueue[T any] struct {
	mux   *sync.RWMutex
	data  []T
	zero  T
	head  int
	tail  int
	count int

	enqueue *semaphore.Weighted
	dequeue *semaphore.Weighted
}

func NewSliceQueue[T any](capacity int) *SliceQueue[T] {
	sq := &SliceQueue[T]{
		data:    make([]T, capacity),
		mux:     &sync.RWMutex{},
		enqueue: semaphore.NewWeighted(int64(capacity)),
		dequeue: semaphore.NewWeighted(int64(capacity)),
	}
	_ = sq.dequeue.Acquire(context.Background(), int64(capacity))
	return sq
}

func (s *SliceQueue[T]) Enqueue(ctx context.Context, val T) (err error) {
	err = s.enqueue.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	if ctx.Err() != nil {
		s.enqueue.Release(1)
		return err
	}
	s.data[s.tail] = val
	s.tail++
	s.count++
	if s.tail >= len(s.data) {
		s.tail = 0
	}
	s.dequeue.Release(1)
	return nil
}

func (s *SliceQueue[T]) Dequeue(ctx context.Context) (val T, err error) {
	err = s.dequeue.Acquire(ctx, 1)
	if err != nil {
		return s.zero, err
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	if ctx.Err() != nil {
		s.dequeue.Release(1)
		return s.zero, ctx.Err()
	}
	val = s.data[s.head]
	s.data[s.head] = s.zero
	s.head++
	s.count--
	if s.head >= len(s.data) {
		s.head = 0
	}
	s.enqueue.Release(1)
	return val, nil
}

func (s *SliceQueue[T]) isEmpty() bool {
	return s.count == 0
}

func (s *SliceQueue[T]) IsEmpty() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isEmpty()
}

func (s *SliceQueue[T]) isFull() bool {
	return s.count == len(s.data)
}

func (s *SliceQueue[T]) IsFull() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isFull()
}
