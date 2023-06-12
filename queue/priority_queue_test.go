package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityQueue_EnqueueDequeue(t *testing.T) {
	minQueue := NewPriorityQueue[int](4, func(src, dst int) int {
		if src < dst {
			return 1
		} else if src == dst {
			return 0
		}
		return -1
	})
	_ = minQueue.Enqueue(4)
	_ = minQueue.Enqueue(3)
	_ = minQueue.Enqueue(2)
	_ = minQueue.Enqueue(1)

	val, err := minQueue.Dequeue()
	assert.Equal(t, 1, val)
	val, err = minQueue.Dequeue()
	assert.Equal(t, 2, val)
	val, err = minQueue.Dequeue()
	assert.Equal(t, 3, val)
	val, err = minQueue.Dequeue()
	assert.Equal(t, 4, val)
	assert.NoError(t, err)

	maxQueue := NewPriorityQueue[int](4, func(src, dst int) int {
		if src > dst {
			return 1
		} else if src == dst {
			return 0
		}
		return -1
	})
	_ = maxQueue.Enqueue(1)
	_ = maxQueue.Enqueue(2)
	_ = maxQueue.Enqueue(3)
	_ = maxQueue.Enqueue(4)
	val, err = maxQueue.Dequeue()
	assert.Equal(t, 4, val)
	val, err = maxQueue.Dequeue()
	assert.Equal(t, 3, val)
	val, err = maxQueue.Dequeue()
	assert.Equal(t, 2, val)
	val, err = maxQueue.Dequeue()
	assert.Equal(t, 1, val)
	assert.NoError(t, err)
}
