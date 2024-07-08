package pool

import (
	"context"
	"errors"
	"github.com/ac-zht/gotools/option"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func TestOnDemandBlockTaskPool_States(t *testing.T) {
	t.Parallel()

	t.Run("调用State方法时使用已取消的context应该返回错误", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 3)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = pool.States(ctx, time.Microsecond)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("调用ShutdownNow方法后再调用States方法应该返回错误", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 3)
		assert.NoError(t, err)

		err = pool.Start()
		assert.NoError(t, err)

		_, err = pool.ShutdownNow()
		assert.NoError(t, err)

		_, err = pool.States(context.Background(), time.Millisecond)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("调用Shutdown方法后再调用States方法应该返回错误", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 3)
		assert.NoError(t, err)

		err = pool.Start()
		assert.NoError(t, err)

		done, err := pool.Shutdown()
		assert.NoError(t, err)

		<-done
		_, err = pool.States(context.Background(), time.Millisecond)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("调用States方法返回的chan应该能够正常读取数据", func(t *testing.T) {
		t.Parallel()
		pool, err := NewOnDemandBlockTaskPool(1, 3)
		assert.NoError(t, err)

		ch, err := pool.States(context.Background(), time.Millisecond)
		assert.NoError(t, err)
		assert.NotZero(t, <-ch)
	})

	t.Run("当调用States方法时传入的context超时返回的chan应该被关闭", func(t *testing.T) {
		t.Parallel()

		initGo, queueSize := 1, 3
		pool, syncChan := testNewRunningStateTaskPoolWithQueueFullFilled(t, initGo, queueSize)

		ctx, cancel := context.WithCancel(context.Background())
		ch, err := pool.States(ctx, time.Millisecond)
		assert.NoError(t, err)

		go func() {
			<-time.After(3 * time.Millisecond)
			cancel()
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}

		close(syncChan)
		_, err = pool.Shutdown()
		assert.NoError(t, err)
	})

	t.Run("调用Shutdown方法应该关闭States方法返回的chan", func(t *testing.T) {
		t.Parallel()
		pool := testNewRunningStateTaskPool(t, 1, 3)
		ch, err := pool.States(context.Background(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(5 * time.Millisecond)
			_, err := pool.Shutdown()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})

	t.Run("调用ShutdownNow方法应该关闭States方法返回的chan", func(t *testing.T) {
		t.Parallel()

		pool := testNewRunningStateTaskPool(t, 1, 3)

		ch, err := pool.States(context.Background(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(5 * time.Millisecond)
			_, err := pool.ShutdownNow()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})
}

func TestOnDemandBlockTaskPool_In_Created_State(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, -1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(1, 0)
		assert.NoError(t, err)
		assert.NotNil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(-1, 1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(0, 1)
		assert.ErrorIs(t, err, errInvalidArgument)
		assert.Nil(t, pool)

		pool, err = NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})

	t.Run("With Options", func(t *testing.T) {
		t.Parallel()

		//未定义核心和最大Go时，都设置为初始Go
		initGo := 10
		pool, err := NewOnDemandBlockTaskPool(initGo, 10)
		assert.NoError(t, err)
		assert.Equal(t, int32(initGo), pool.initGo)
		assert.Equal(t, int32(initGo), pool.coreGo)
		assert.Equal(t, int32(initGo), pool.maxGo)
		assert.Equal(t, defaultMaxIdleTime, pool.maxIdleTime)

		coreGo, maxGo, maxIdleTime := int32(20), int32(30), 10*time.Second
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
		assert.NoError(t, err)

		assert.Equal(t, int32(initGo), pool.initGo)
		assert.Equal(t, coreGo, pool.coreGo)
		assert.Equal(t, maxGo, pool.maxGo)
		assert.Equal(t, maxIdleTime, pool.maxIdleTime)
		//init10，core20，max10->20
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo))
		assert.NoError(t, err)
		assert.Equal(t, pool.coreGo, pool.maxGo)

		//init30，core20，max30->20
		initGo, coreGo = 30, 20
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithMaxGo(maxGo))
		assert.NoError(t, err)
		assert.Equal(t, pool.maxGo, pool.coreGo)

		//max<init=core
		initGo, maxGo = 30, 10
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		//max<core<init
		initGo, coreGo, maxGo = 30, 20, 10
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		//core<max<init
		initGo, coreGo, maxGo = 30, 10, 20
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		//core<init<max
		initGo, coreGo, maxGo = 20, 10, 30
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		//max<init<core
		initGo, coreGo, maxGo = 20, 30, 10
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		//init<max<core
		initGo, coreGo, maxGo = 10, 30, 20
		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithCoreGo(coreGo), WithMaxGo(maxGo))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(-0.1))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)

		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(1.0))
		assert.NotNil(t, pool)
		assert.NoError(t, err)

		pool, err = NewOnDemandBlockTaskPool(initGo, 10, WithQueueBacklogRate(1.1))
		assert.Nil(t, pool)
		assert.ErrorIs(t, err, errInvalidArgument)
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			assert.ErrorIs(t, pool.Submit(context.Background(), nil), errTaskIsInvalid)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 3)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitValidTask(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()

			pool, _ := NewOnDemandBlockTaskPool(1, 1)
			assert.Equal(t, stateCreated, pool.internalState())
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, stateCreated, pool.internalState())
		})
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, stateCreated, pool.internalState())

		done, err := pool.Shutdown()
		assert.Nil(t, done)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})

	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool, err := NewOnDemandBlockTaskPool(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, stateCreated, pool.internalState())

		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsNotRunning)
		assert.Equal(t, stateCreated, pool.internalState())
	})
}

func TestOnDemandBlockTaskPool_In_Running_State(t *testing.T) {
	t.Parallel()

	t.Run("Start —— 使TaskPool状态由Created变为Running", func(t *testing.T) {
		t.Parallel()

		pool, _ := NewOnDemandBlockTaskPool(1, 1)
		errChan := make(chan error)
		go func() {
			time.Sleep(1 * time.Millisecond)
			errChan <- pool.Start()
		}()

		assert.Equal(t, stateCreated, pool.internalState())
		testSubmitBlockingAndTimeout(t, pool)

		assert.NoError(t, <-errChan)
		assert.Equal(t, stateRunning, pool.internalState())

		assert.ErrorIs(t, pool.Start(), errTaskPoolIsStarted)
		assert.Equal(t, stateRunning, pool.internalState())
	})

	t.Run("Start —— 在TaskPool启动前队列中已有任务，启动后不再Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("WithCoreGo,WithMaxIdleTime，所需要协程数 <= 允许创建的协程数", func(t *testing.T) {
			t.Parallel()

			initGo, coreGo, maxIdleTime := 1, 3, 3*time.Millisecond
			queueSize := coreGo

			needGo, allowGo := queueSize-initGo, coreGo-initGo
			assert.LessOrEqual(t, needGo, allowGo)

			pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithCoreGo(int32(coreGo)), WithMaxIdleTime(maxIdleTime))
			assert.NoError(t, err)

			assert.Equal(t, int32(0), pool.numOfGo())

			done := make(chan struct{}, coreGo)
			wait := make(chan struct{}, coreGo)

			for i := 0; i < coreGo; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}

			assert.Equal(t, int32(0), pool.numOfGo())

			assert.NoError(t, pool.Start())
			for i := 0; i < coreGo; i++ {
				<-wait
			}
			assert.Equal(t, int32(coreGo), pool.numOfGo())
			close(done)
		})

		t.Run("WithMaxGo, 所需要协程数 > 允许创建的协程数", func(t *testing.T) {
			t.Parallel()

			initGo, maxGo := 3, 5
			queueSize := maxGo + 1

			needGo, allowGo := queueSize-initGo, maxGo-initGo
			assert.Greater(t, needGo, allowGo)

			pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithMaxGo(int32(maxGo)))
			assert.NoError(t, err)

			assert.Equal(t, int32(0), pool.numOfGo())

			done := make(chan struct{}, queueSize)
			wait := make(chan struct{}, queueSize)

			for i := 0; i < queueSize; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}

			assert.Equal(t, int32(0), pool.numOfGo())
			assert.NoError(t, pool.Start())

			for i := 0; i < maxGo; i++ {
				<-wait
			}
			assert.Equal(t, int32(maxGo), pool.numOfGo())
			close(done)
		})
	})

	t.Run("Start —— 与Submit并发调用，WithCoreGo,WithMaxIdleTime,WithMaxGo，所需要协程数 < 允许创建的协程数", func(t *testing.T) {
		t.Parallel()

		initGo, coreGo, maxGo, maxIdleTime := 2, 4, 6, 3*time.Millisecond
		queueSize := coreGo

		needGo, allowGo := queueSize-initGo, maxGo-initGo
		assert.Less(t, needGo, allowGo)

		pool, err := NewOnDemandBlockTaskPool(initGo, queueSize, WithCoreGo(int32(coreGo)), WithMaxGo(int32(maxGo)), WithMaxIdleTime(maxIdleTime))
		assert.NoError(t, err)

		wait := make(chan struct{}, queueSize)
		done := make(chan struct{}, queueSize)

		errChan := make(chan error)
		go func() {
			time.Sleep(10 * time.Millisecond)
			//初始化线程，消费队列里任务
			errChan <- pool.Start()
		}()

		for i := 0; i < maxGo; i++ {
			//任务提交4个，阻塞2个等任务池Start
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				wait <- struct{}{}
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.NoError(t, <-errChan)

		for i := 0; i < queueSize; i++ {
			<-wait
		}

		assert.Equal(t, int32(maxGo), pool.numOfGo())

		close(done)
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		t.Run("提交非法Task", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 1)
			assert.ErrorIs(t, pool.Submit(context.Background(), nil), errTaskIsInvalid)
			assert.Equal(t, stateRunning, pool.internalState())
		})

		t.Run("正常提交Task", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 1)
			testSubmitValidTask(t, pool)
			assert.Equal(t, stateRunning, pool.internalState())
		})

		t.Run("阻塞提交并导致超时", func(t *testing.T) {
			t.Parallel()

			pool := testNewRunningStateTaskPool(t, 1, 1)
			testSubmitBlockingAndTimeout(t, pool)
			assert.Equal(t, stateRunning, pool.internalState())
		})
	})

	t.Run("工作协程", func(t *testing.T) {
		t.Parallel()

		t.Run("保持在初始数不变", func(t *testing.T) {
			initGo, queueSize := 1, 3
			pool := testNewRunningStateTaskPool(t, initGo, queueSize)

			n := queueSize
			done1 := make(chan struct{}, n)
			wait := make(chan struct{}, n)

			for i := 0; i < n; i++ {
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done1
					return nil
				}))
				assert.NoError(t, err)
			}

			//执行中
			for i := 0; i < initGo; i++ {
				<-wait
			}

			assert.Equal(t, int32(initGo), pool.numOfGo())

			//执行完
			for i := 0; i < initGo; i++ {
				done1 <- struct{}{}
			}

			//队列积压任务执行
			for i := 0; i < n-initGo; i++ {
				<-wait
				assert.Equal(t, int32(initGo), pool.numOfGo())
				done1 <- struct{}{}
			}
		})

		t.Run("从初始数达到核心数", func(t *testing.T) {
			t.Parallel()

			t.Run("核心数比初始数多1个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(3), 3*time.Millisecond, 0.1
				queueSize := int(coreGo)
				testExtendGoFromInitGoToCoreGo(t, initGo, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("核心数比初始数多n个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime, queueBacklogRate := int32(2), int32(5), 3*time.Millisecond, 0.1
				queueSize := int(coreGo)
				testExtendGoFromInitGoToCoreGo(t, initGo, queueSize, coreGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("在(初始数,核心数]区间的协程运行完任务后，在等待退出期间再次抢到任务", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxIdleTime := int32(1), int32(6), 100*time.Millisecond
				queueSize := int(coreGo)

				pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))

				assert.Equal(t, initGo, pool.numOfGo())
				t.Log("1")
				done := make(chan struct{}, queueSize)
				wait := make(chan struct{}, queueSize)

				for i := 0; i < queueSize; i++ {
					i := i
					err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
						wait <- struct{}{}
						<-done
						t.Log("task done", i)
						return nil
					}))
					assert.NoError(t, err)
				}
				t.Log("2")
				for i := 0; i < queueSize; i++ {
					t.Log("wait", i)
					<-wait
				}

				assert.Equal(t, coreGo, pool.numOfGo())

				close(done)
				t.Log("3")
				err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-done
					t.Log("task done [X]")
					return nil
				}))
				assert.NoError(t, err)
				t.Log("4")
				for pool.numOfGo() > initGo {
					t.Log("loop", "numOfGo", pool.numOfGo(), "timeoutGroup", pool.timeoutGroup.size())
					time.Sleep(maxIdleTime)
				}
				assert.Equal(t, initGo, pool.numOfGo())
			})
		})

		t.Run("从核心数到达最大数", func(t *testing.T) {
			t.Parallel()

			t.Run("最大数比核心数多1个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(2), int32(4), int32(5), 3*time.Millisecond, 0.1
				queueSize := int(maxGo)
				testExtendGoFromInitGoToCoreGoAndMaxGo(t, initGo, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})

			t.Run("最大数比核心数多n个", func(t *testing.T) {
				t.Parallel()

				initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(2), int32(4), int32(7), 3*time.Millisecond, 0.1
				queueSize := int(maxGo)
				testExtendGoFromInitGoToCoreGoAndMaxGo(t, initGo, queueSize, coreGo, maxGo, maxIdleTime, WithQueueBacklogRate(queueBacklogRate))
			})
		})
	})
}

func TestOnDemandBlockTaskPool_In_Closing_State(t *testing.T) {
	t.Parallel()

	t.Run("Shutdown —— 使TaskPool状态由Running变为Closing", func(t *testing.T) {
		t.Parallel()

		initGo, queueSize := 2, 4
		pool := testNewRunningStateTaskPool(t, initGo, queueSize)

		n := initGo + queueSize + 1
		eg := new(errgroup.Group)
		wait := make(chan struct{})
		done := make(chan struct{})
		for i := 0; i < n; i++ {
			eg.Go(func() error {
				return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					wait <- struct{}{}
					<-done
					return nil
				}))
			})
		}
		for i := 0; i < initGo; i++ {
			<-wait
		}
		shutdown, err := pool.Shutdown()
		assert.NoError(t, err)
		assert.ErrorIs(t, eg.Wait(), errTaskPoolIsClosing)

		shutdown2, err2 := pool.Shutdown()
		assert.Nil(t, shutdown2)
		assert.ErrorIs(t, err2, errTaskPoolIsClosing)

		assert.Equal(t, stateClosing, pool.internalState())

		close(wait)
		close(done)
		<-shutdown
		assert.Equal(t, stateStopped, pool.internalState())

		shutdown3, err3 := pool.Shutdown()
		assert.Nil(t, shutdown3)
		assert.ErrorIs(t, err3, errTaskPoolIsStopped)
	})

	t.Run("Shutdown —— 协程数按需扩展至maxGo，调用Shutdown成功后，所有协程运行完任务后可以自动退出", func(t *testing.T) {
		t.Parallel()

		initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, initGo, coreGo)
		assert.LessOrEqual(t, coreGo, maxGo)

		wait := make(chan struct{})
		done := make(chan struct{})

		for i := int32(0); i < maxGo; i++ {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				wait <- struct{}{}
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		shutdown, err := pool.Shutdown()
		assert.NoError(t, err)

		for i := int32(0); i < maxGo; i++ {
			<-wait
			t.Log("wait", i)
		}
		assert.Equal(t, maxGo, pool.numOfGo())

		close(wait)
		close(done)
		<-shutdown
		for pool.numOfGo() > 0 {
		}
		assert.Equal(t, int32(0), pool.numOfGo())
	})
	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, 2, 4)
		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			assert.ErrorIs(t, pool.Start(), errTaskPoolIsClosing)
		}
		close(wait)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})
	t.Run("ShutdownNow", func(t *testing.T) {
		t.Parallel()

		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, 2, 4)
		done, err := pool.Shutdown()
		assert.NoError(t, err)

		select {
		case <-done:
		default:
			tasks, err := pool.ShutdownNow()
			assert.ErrorIs(t, err, errTaskPoolIsClosing)
			assert.Nil(t, tasks)
		}
		close(wait)
		<-done
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

type ShutdownNowResult struct {
	tasks []Task
	err   error
}

func TestOnDemandBlockTaskPool_In_Stopped_State(t *testing.T) {
	t.Parallel()

	t.Run("ShutdownNow —— 使TasksPool状态由Running变为Stopped", func(t *testing.T) {
		t.Parallel()

		initGo, queueSize := 2, 4
		pool, wait := testNewRunningStateTaskPoolWithQueueFullFilled(t, initGo, queueSize)

		eg := new(errgroup.Group)
		for i := 0; i < queueSize; i++ {
			eg.Go(func() error {
				return pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
					<-wait
					return nil
				}))
			})
		}
		result := make(chan ShutdownNowResult, 1)
		go func() {
			tasks, err := pool.ShutdownNow()
			result <- ShutdownNowResult{tasks: tasks, err: err}
			close(wait)
		}()

		assert.ErrorIs(t, eg.Wait(), errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())

		r := <-result
		assert.NoError(t, r.err)
		assert.Equal(t, queueSize, len(r.tasks))

		tasks, err := pool.ShutdownNow()
		assert.Nil(t, tasks)
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})
	t.Run("ShutdownNow —— 工作协程数扩展至maxGo后，调用ShutdownNow成功，所有协程能够接收到信号", func(t *testing.T) {
		t.Parallel()

		initGo, coreGo, maxGo, maxIdleTime, queueBacklogRate := int32(1), int32(3), int32(5), 10*time.Millisecond, 0.1
		queueSize := int(maxGo)
		pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime), WithQueueBacklogRate(queueBacklogRate))

		assert.LessOrEqual(t, initGo, coreGo)
		assert.LessOrEqual(t, coreGo, maxGo)

		done := make(chan struct{})

		for i := 0; i < queueSize; i++ {
			err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		tasks, err := pool.ShutdownNow()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(tasks), 0)

		close(done)
		for pool.numOfGo() > 0 {
		}
		assert.Equal(t, int32(0), pool.numOfGo())
	})

	t.Run("Start", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		assert.ErrorIs(t, pool.Start(), errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Submit", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			return nil
		}))
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Parallel()

		pool := testNewStoppedStateTaskPool(t, 1, 1)
		_, err := pool.Shutdown()
		assert.ErrorIs(t, err, errTaskPoolIsStopped)
		assert.Equal(t, stateStopped, pool.internalState())
	})
}

func testExtendGoFromInitGoToCoreGo(t *testing.T, initGo int32, queueSize int, coreGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {

	opts = append(opts, WithCoreGo(coreGo), WithMaxIdleTime(maxIdleTime))
	pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, opts...)
	assert.Equal(t, initGo, pool.numOfGo())
	assert.LessOrEqual(t, initGo, coreGo)

	done := make(chan struct{})
	wait := make(chan struct{}, coreGo)

	t.Log("XX")
	//将任务稳定在initGo
	for i := int32(0); i < initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		t.Log("submit", i)
	}

	t.Log("YY")
	for i := int32(0); i < initGo; i++ {
		<-wait
	}

	assert.GreaterOrEqual(t, pool.numOfGo(), initGo)

	t.Log("ZZ")
	//加任务
	for i := int32(1); i <= coreGo-initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait coreGo", coreGo, i, pool.numOfGo())
	}

	t.Log("UU")

	assert.Equal(t, pool.numOfGo(), coreGo)
	close(done)

	//等到最大空闲时间后稳定在initGo
	for pool.numOfGo() > initGo {
	}

	assert.Equal(t, initGo, pool.numOfGo())
}

func testExtendGoFromInitGoToCoreGoAndMaxGo(t *testing.T, initGo int32, queueSize int, coreGo, maxGo int32, maxIdleTime time.Duration, opts ...option.Option[OnDemandBlockTaskPool]) {
	opts = append(opts, WithCoreGo(coreGo), WithMaxGo(maxGo), WithMaxIdleTime(maxIdleTime))
	pool := testNewRunningStateTaskPool(t, int(initGo), queueSize, opts...)

	assert.Equal(t, initGo, pool.numOfGo())

	assert.LessOrEqual(t, initGo, coreGo)
	assert.LessOrEqual(t, coreGo, maxGo)

	done := make(chan struct{})
	wait := make(chan struct{}, maxGo)

	t.Log("00")
	for i := int32(0); i < initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		t.Log("Submit", i)
	}
	t.Log("AA")
	for i := int32(0); i < initGo; i++ {
		<-wait
	}
	assert.GreaterOrEqual(t, pool.numOfGo(), initGo)

	t.Log("BB")

	for i := int32(1); i <= coreGo-initGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait coreGo", coreGo, i, pool.numOfGo())
	}

	t.Log("CC")

	assert.GreaterOrEqual(t, pool.numOfGo(), coreGo)

	for i := int32(1); i <= maxGo-coreGo; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			wait <- struct{}{}
			<-done
			return nil
		}))
		assert.NoError(t, err)
		<-wait
		t.Log("after wait maxGo", maxGo, i, pool.numOfGo())
	}

	t.Log("DD")
	assert.Equal(t, maxGo, pool.numOfGo())
	close(done)

	for pool.numOfGo() > initGo {
	}

	assert.Equal(t, initGo, pool.numOfGo())
}

func testSubmitBlockingAndTimeout(t *testing.T, pool *OnDemandBlockTaskPool) {
	done := make(chan struct{})
	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		<-done
		return nil
	}))
	assert.NoError(t, err)

	n := cap(pool.queue) + 2
	eg := new(errgroup.Group)

	for i := 0; i < n; i++ {
		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			return pool.Submit(ctx, TaskFunc(func(ctx context.Context) error {
				<-done
				return nil
			}))
		})
	}
	assert.ErrorIs(t, eg.Wait(), context.DeadlineExceeded)
	close(done)
}

func testSubmitValidTask(t *testing.T, pool *OnDemandBlockTaskPool) {
	err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		return nil
	}))
	assert.NoError(t, err)
	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		panic("task panic")
	}))
	assert.NoError(t, err)
	err = pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
		return errors.New("fake error")
	}))
	assert.NoError(t, err)
}

func testNewRunningStateTaskPool(t *testing.T, initGo int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) *OnDemandBlockTaskPool {
	pool, _ := NewOnDemandBlockTaskPool(initGo, queueSize, opts...)
	assert.Equal(t, stateCreated, pool.internalState())
	assert.NoError(t, pool.Start())
	assert.Equal(t, stateRunning, pool.internalState())
	return pool
}

func testNewStoppedStateTaskPool(t *testing.T, initGo int, queueSize int) *OnDemandBlockTaskPool {
	pool := testNewRunningStateTaskPool(t, initGo, queueSize)
	tasks, err := pool.ShutdownNow()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, stateStopped, pool.internalState())
	return pool
}

func testNewRunningStateTaskPoolWithQueueFullFilled(t *testing.T, initGo int, queueSize int) (*OnDemandBlockTaskPool, chan struct{}) {
	pool := testNewRunningStateTaskPool(t, initGo, queueSize)
	wait := make(chan struct{})
	for i := 0; i < initGo+queueSize; i++ {
		err := pool.Submit(context.Background(), TaskFunc(func(ctx context.Context) error {
			<-wait
			return nil
		}))
		assert.NoError(t, err)
	}
	return pool, wait
}

func TestGroup(t *testing.T) {
	t.Parallel()

	var n, i int32
	n = 10
	g := &group{mp: make(map[int32]struct{}), mux: &sync.RWMutex{}}
	for i = 0; i < n; i++ {
		assert.False(t, g.isIn(i))
		g.add(i)
		assert.True(t, g.isIn(i))
		assert.Equal(t, i+1, g.size())
	}

	assert.Equal(t, n, g.size())

	for i = 0; i < n; i++ {
		g.delete(i)
		assert.Equal(t, n-i-1, g.size())
	}

	assert.Equal(t, int32(0), g.size())
	assert.False(t, g.isIn(n+1))

	id := int32(100)
	g.add(id)
	assert.Equal(t, int32(1), g.size())
	assert.True(t, g.isIn(id))
	g.delete(id)
	assert.Equal(t, int32(0), g.size())
	assert.False(t, g.isIn(id))
}
