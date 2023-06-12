package queue

import (
	"context"
	"sync"
	"time"
)

type delayQueue[T Delayable] struct {
	mux           *sync.RWMutex
	p             *priorityQueue[T]
	enqueueSignal *Cond
	dequeueSignal *Cond
	zero          T
}

func NewCond(l sync.Locker) *Cond {
	return &Cond{
		signal: make(chan struct{}),
		L:      l,
	}
}

func NewDelayQueue[T Delayable](capacity int) *delayQueue[T] {
	d := &delayQueue[T]{
		mux: &sync.RWMutex{},
	}
	d.enqueueSignal = NewCond(d.mux)
	d.dequeueSignal = NewCond(d.mux)
	d.p = NewPriorityQueue[T](capacity, func(src, dst T) int {
		x := src.Delay()
		y := dst.Delay()
		if x < y {
			return 1
		} else if x == y {
			return 0
		}
		return -1
	})
	return d
}

func (d *delayQueue[T]) Enqueue(ctx context.Context, val T) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		d.mux.Lock()

		select {
		case <-ctx.Done():
			d.mux.Unlock()
			return ctx.Err()
		default:
		}

		err := d.p.Enqueue(val)
		switch err {
		case nil:
			d.enqueueSignal.broadcast()
			return err
		case ErrOutOfCapacity:
			ch := d.dequeueSignal.Signal()
			select {
			case <-ch:
			case <-ctx.Done():
				return ctx.Err()
			}
		default:
			d.mux.Unlock()
			return err
		}
	}
}

func (d *delayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	var timer *time.Timer
	for {
		select {
		case <-ctx.Done():
			return d.zero, ctx.Err()
		default:
		}

		d.mux.Lock()

		select {
		case <-ctx.Done():
			d.mux.Unlock()
			return d.zero, ctx.Err()
		default:
		}

		top, err := d.p.peek()
		switch err {
		case nil:
			delay := top.Delay()
			if delay <= 0 {
				val, err := d.p.Dequeue()
				d.dequeueSignal.broadcast()
				return val, err
			}
			if timer == nil {
				timer = time.NewTimer(delay)
			} else {
				timer.Reset(delay)
			}
			ch := d.dequeueSignal.Signal()
			select {
			case <-timer.C:
			case <-ch:
			case <-ctx.Done():
				return d.zero, ctx.Err()
			}
		case ErrEmptyQueue:
			ch := d.dequeueSignal.Signal()
			select {
			case <-ch:
			case <-ctx.Done():
				return d.zero, ctx.Err()
			}
		default:
			d.mux.Unlock()
			return top, err
		}
	}
}

type Delayable interface {
	Delay() time.Duration
}

type delayElem struct {
	endTime time.Time
}

func (e *delayElem) Delay() time.Duration {
	return e.endTime.Sub(time.Now())
}

type Cond struct {
	signal chan struct{}
	L      sync.Locker
}

func (c *Cond) broadcast() {
	signal := make(chan struct{})
	old := c.signal
	c.signal = signal
	close(old)
	c.L.Unlock()
}

func (c *Cond) Signal() <-chan struct{} {
	res := c.signal
	c.L.Unlock()
	return res
}
