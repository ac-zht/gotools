package rlock

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

var (
	//go:embed scripts/lock.lua
	luaLock string

	//go:embed scripts/unlock.lua
	luaUnlock string

	//go:embed scripts/refresh.lua
	luaRefresh string

	ErrFailedToPreemptLock = errors.New("rlock：抢锁失败")

	ErrLockNotHold = errors.New("rlock：未持有锁")
)

type Client struct {
	client redis.Cmdable
	g      singleflight.Group
	valuer func() string
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{
		client: client,
		valuer: func() string {
			return uuid.NewString()
		},
	}
}

type Lock struct {
	client           redis.Cmdable
	key              string
	val              string
	expiration       time.Duration
	unlock           chan struct{}
	signalUnlockOnce sync.Once
}

func newLock(client redis.Cmdable, key, val string, expiration time.Duration) *Lock {
	return &Lock{
		client:     client,
		key:        key,
		val:        val,
		expiration: expiration,
		unlock:     make(chan struct{}, 1),
	}
}

func (c *Client) SingleflightLock(ctx context.Context, key string, expiration time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	for {
		flag := false
		result := c.g.DoChan(key, func() (interface{}, error) {
			flag = true
			return c.Lock(ctx, key, expiration, retry, timeout)
		})
		select {
		case res := <-result:
			if flag {
				c.g.Forget(key)
				if res.Err != nil {
					return nil, res.Err
				}
				return res.Val.(*Lock), nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	val := c.valuer()
	defer retry.Reset()
	for {
		lctx, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if res == "OK" {
			return newLock(c.client, key, val, expiration), nil
		}
		interval, ok := retry.Next()
		if !ok {
			if err != nil {
				//应用层面上错误
				err = fmt.Errorf("最后一次重试错误：%w", err)
			} else {
				//业务层面上错误
				err = fmt.Errorf("锁被人持有：%w", ErrFailedToPreemptLock)
			}
			return nil, fmt.Errorf("rlock：重试机会耗尽，%w", err)
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := c.valuer()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrFailedToPreemptLock
	}
	return newLock(c.client, key, val, expiration), nil
}

func (l *Lock) Unlock(context context.Context) error {
	res, err := l.client.Eval(context, luaUnlock, []string{l.key}, l.val).Int64()
	defer func() {
		l.signalUnlockOnce.Do(func() {
			l.unlock <- struct{}{}
			close(l.unlock)
		})
	}()
	if err == redis.Nil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	ch := make(chan struct{}, 1)
	defer func() {
		ticker.Stop()
		close(ch)
	}()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				//上一次调用超时未重新处理，同时计时器也触发了
				select {
				case ch <- struct{}{}:
				default:
				}
				continue
			}
			if err != nil {
				return err
			}
		case <-ch:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				select {
				case ch <- struct{}{}:
				default:
				}
				continue
			}
			if err != nil {
				return err
			}
		case <-l.unlock:
			return nil
		}
	}
}
