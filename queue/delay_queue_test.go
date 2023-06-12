package queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestDelayQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name    string
		timeout time.Duration
		d       *delayQueue[Delayable]
		want    Delayable
		wantErr error
	}{
		{
			name:    "context timeout",
			timeout: time.Second * 2,
			d:       NewDelayQueue[Delayable](5),
			want:    nil,
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			val, err := tc.d.Dequeue(ctx)
			cancel()
			assert.Equal(t, tc.want, val)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestDelayQueue_EnqueueDequeue(t *testing.T) {
	testCases := []struct {
		name    string
		d       *delayQueue[Delayable]
		in      Delayable
		out     Delayable
		wantErr error
	}{
		{
			name:    "success",
			d:       NewDelayQueue[Delayable](5),
			in:      &delayElem{endTime: time.Now().Add(time.Second * 2)},
			out:     &delayElem{endTime: time.Now().Add(time.Second * 2)},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				val Delayable
				err error
			)
			err = tc.d.Enqueue(context.Background(), tc.in)
			assert.Equal(t, tc.wantErr, err)
			val, err = tc.d.Dequeue(context.Background())
			assert.Equal(t, tc.out, val)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestDelayQueue_ConcurrentEnqueueDequeue(t *testing.T) {
	dp := NewDelayQueue[Delayable](6)
	var outs [10]int64
	for i := 0; i < 10; i++ {
		go func() {
			num := rand.Intn(4) + 1
			t := time.Now().Add(time.Second * time.Duration(int64(num)))
			_ = dp.Enqueue(context.Background(), &delayElem{endTime: t})
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			val, err := dp.Dequeue(context.Background())
			if err != nil {
				return
			}
			outs[i] = (*(val.(*delayElem))).endTime.Unix()
		}(i)
	}
	time.Sleep(time.Second * 7)
	success := true
	for _, val := range outs {
		if val <= 0 {
			success = false
			break
		}
	}
	assert.Equal(t, true, success)
}
