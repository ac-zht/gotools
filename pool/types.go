package pool

import (
	"context"
	"time"
)

type TaskPool interface {
	Submit(ctx context.Context, task Task) error

	Start() error

	Shutdown() (<-chan struct{}, error)

	ShutdownNow() ([]Task, error)

	States(ctx context.Context, interval time.Duration) (<-chan State, error)
}

type Task interface {
	Run(ctx context.Context) error
}

type State struct {
	PoolState      int32
	GoCnt          int32
	WaitingTaskCnt int
	QueueSize      int
	RunningTaskCnt int32
	TimesStamp     int64
}
