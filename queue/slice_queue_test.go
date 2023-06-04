package queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	weight := semaphore.NewWeighted(5)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 5; i++ {
		go func() {
			err := weight.Acquire(context.Background(), 1)
			t.Log(err)
			wg.Done()
		}()
		go func() {
			time.Sleep(time.Second)
			weight.Release(1)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestSliceQueue_In(t *testing.T) {
	testCases := []struct {
		name string
		ctx  context.Context
		in   int
		q    *SliceQueue[int]

		wantErr  error
		wantData []int
	}{
		{
			name: "超时",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			}(),
			in: 10,
			q: func() *SliceQueue[int] {
				q := NewSliceQueue[int](2)
				_ = q.In(context.Background(), 11)
				_ = q.In(context.Background(), 12)
				return q
			}(),
			wantErr:  context.DeadlineExceeded,
			wantData: []int{11, 12},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.q.In(tc.ctx, tc.in)
			assert.Equal(t, tc.wantData, tc.q.data)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestSliceQueue_InOut(t *testing.T) {
	sq := NewSliceQueue[int](5)
	closed := false
	for i := 0; i < 3; i++ {
		go func() {
			for {
				if closed {
					return
				}
				in := rand.Int()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_ = sq.In(ctx, in)
				cancel()
			}
		}()
	}

	for j := 0; j < 3; j++ {
		go func() {
			for {
				if closed {
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, _ = sq.Out(ctx)
				cancel()
			}
		}()
	}

	time.Sleep(time.Second * 1)
	closed = true
}
