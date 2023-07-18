package pool

import (
	"context"
	"errors"
	"fmt"
	"geektime_work/option"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	stateCreated int32 = 1
	stateRunning int32 = 2
	stateClosing int32 = 3
	stateStopped int32 = 4
	stateLocked  int32 = 5

	errTaskPoolIsNotRunning = errors.New("pool: TaskPool未运行")
	errTaskPoolIsClosing    = errors.New("pool: TaskPool关闭中")
	errTaskPoolIsStopped    = errors.New("pool: TaskPool已停止")
	errTaskPoolIsStarted    = errors.New("pool: TaskPool已运行")
	errTaskIsInvalid        = errors.New("pool: Task非法")
	errTaskRunningPanic     = errors.New("pool: Task运行时异常")

	errInvalidArgument = errors.New("pool: 参数非法")
	defaultMaxIdleTime = 10 * time.Second

	panicBuffLen = 2048
)

type OnDemandBlockTaskPool struct {
	state int32

	coreGo  int32
	initGo  int32
	maxGo   int32
	totalGo int32

	runningTaskCnt int32
	queue          chan Task

	timeoutGroup     *group
	id               int32
	maxIdleTime      time.Duration
	queueBacklogRate float64

	interruptCtx    context.Context
	interruptCancel context.CancelFunc

	mux *sync.RWMutex
}

type group struct {
	mp  map[int32]struct{}
	cnt int32
	mux *sync.RWMutex
}

func (g *group) isIn(id int32) bool {
	g.mux.RLock()
	defer g.mux.RUnlock()
	_, ok := g.mp[id]
	return ok
}

func (g *group) add(id int32) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.mp[id] = struct{}{}
	g.cnt++
}

func (g *group) delete(id int32) {
	g.mux.Lock()
	defer g.mux.Unlock()
	delete(g.mp, id)
	g.cnt--
}

func (g *group) size() int32 {
	g.mux.RLock()
	defer g.mux.RUnlock()
	return g.cnt
}

type TaskFunc func(ctx context.Context) error

func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

type taskWrapper struct {
	t Task
}

func (tw *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("%w：%s", errTaskRunningPanic, fmt.Sprintf("[PANIC]:\t%+v\n%s\n", r, buf))
		}
	}()
	return tw.t.Run(ctx)
}

func NewOnDemandBlockTaskPool(initGo int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) (*OnDemandBlockTaskPool, error) {
	if initGo < 1 {
		return nil, fmt.Errorf("%w: initGo应该大于0", errInvalidArgument)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("%w: queueSize应该大于等于0", errInvalidArgument)
	}
	pool := &OnDemandBlockTaskPool{
		initGo:      int32(initGo),
		coreGo:      int32(initGo),
		maxGo:       int32(initGo),
		maxIdleTime: defaultMaxIdleTime,
		queue:       make(chan Task, queueSize),
		mux:         &sync.RWMutex{},
		timeoutGroup: &group{
			mp:  make(map[int32]struct{}),
			mux: &sync.RWMutex{},
		},
	}
	atomic.StoreInt32(&pool.state, stateCreated)
	ctx := context.Background()
	pool.interruptCtx, pool.interruptCancel = context.WithCancel(ctx)

	for _, opt := range opts {
		opt(pool)
	}

	//规避只设置核心数出错
	if pool.maxGo == pool.initGo && pool.coreGo != pool.initGo {
		pool.maxGo = pool.coreGo
	}

	if !(pool.initGo <= pool.coreGo && pool.coreGo <= pool.maxGo) {
		return nil, errInvalidArgument
	}

	if pool.queueBacklogRate < float64(0) || pool.queueBacklogRate > float64(1) {
		return nil, errInvalidArgument
	}
	return pool, nil
}

func WithCoreGo(coreGo int32) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.coreGo = coreGo
	}
}

func WithMaxGo(maxGo int32) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.maxGo = maxGo
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.maxIdleTime = maxIdleTime
	}
}

func WithQueueBacklogRate(queueBacklogRate float64) option.Option[OnDemandBlockTaskPool] {
	return func(pool *OnDemandBlockTaskPool) {
		pool.queueBacklogRate = queueBacklogRate
	}
}

func (o *OnDemandBlockTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return errTaskIsInvalid
	}
	task = &taskWrapper{t: task}
	//当阻塞时一直去尝试提交，直到超时
	for {
		if atomic.LoadInt32(&o.state) == stateStopped {
			return errTaskPoolIsStopped
		}
		if atomic.LoadInt32(&o.state) == stateClosing {
			return errTaskPoolIsClosing
		}
		ok, err := o.trySubmit(ctx, task, stateCreated)
		if ok || err != nil {
			return err
		}

		ok, err = o.trySubmit(ctx, task, stateRunning)
		if ok || err != nil {
			return err
		}
	}
}

func (o *OnDemandBlockTaskPool) trySubmit(ctx context.Context, task Task, state int32) (bool, error) {
	//上锁，任务进来后不可被关闭阻塞
	if atomic.CompareAndSwapInt32(&o.state, state, stateLocked) {
		defer atomic.CompareAndSwapInt32(&o.state, stateLocked, state)
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case o.queue <- task:
			//如可创建线程就建新线程处理任务，否则加载到队列等待由已有线程去处理
			if state == stateRunning && o.allowToCreateGoroutine() {
				o.increaseTotalGo(1)
				id := atomic.AddInt32(&o.id, 1)
				go o.goroutine(id)
			}
			return true, nil
		default:
		}
	}
	return false, nil
}

/**
工作线程
*/
func (o *OnDemandBlockTaskPool) numOfGo() int32 {
	o.mux.RLock()
	defer o.mux.RUnlock()
	return o.totalGo
}

func (o *OnDemandBlockTaskPool) allowToCreateGoroutine() bool {
	o.mux.RLock()
	defer o.mux.RUnlock()
	rate := float64(len(o.queue)) / float64(cap(o.queue))
	return o.totalGo < o.maxGo && (rate != 0 && rate >= o.queueBacklogRate)
}

func (o *OnDemandBlockTaskPool) increaseTotalGo(n int32) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.totalGo += n
}

func (o *OnDemandBlockTaskPool) decreaseTotalGo(n int32) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.totalGo -= n
}

/**
处理任务
*/
func (o *OnDemandBlockTaskPool) goroutine(id int32) {

	//实例化一个已超时的timer
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	for {
		select {
		case <-o.interruptCtx.Done():
			//接收到任务池关闭信号，所有线程中断关闭
			o.decreaseTotalGo(1)
			return
		case <-timer.C:
			o.mux.Lock()
			o.totalGo--
			o.timeoutGroup.delete(id)
			o.mux.Unlock()
			return
		case task, ok := <-o.queue:
			if o.timeoutGroup.isIn(id) {
				o.timeoutGroup.delete(id)
				//保证协程从超时组里捞出来后，不会因为未再次开始计时，而又立马超时
				if !timer.Stop() {
					<-timer.C
				}
			}
			//拿不到任务，且任务队列已关闭，中断线程
			if !ok {
				o.decreaseTotalGo(1)
				if o.numOfGo() <= 0 {
					if atomic.CompareAndSwapInt32(&o.state, stateClosing, stateStopped) {
						o.interruptCancel()
					}
				}
				return
			}

			atomic.AddInt32(&o.runningTaskCnt, 1)
			_ = task.Run(o.interruptCtx)
			atomic.AddInt32(&o.runningTaskCnt, -1)

			o.mux.Lock()
			//当前线程数超过核心数且待处理任务量为空或小于线程数时直接退出
			noTaskToExecute := len(o.queue) == 0 || int32(len(o.queue)) < o.totalGo
			if o.totalGo > o.coreGo && noTaskToExecute {
				o.totalGo--
				o.mux.Unlock()
				return
			}

			//当前线程数超过 初始数+超时组线程，加入到超时组
			if o.totalGo > o.initGo+o.timeoutGroup.size() {
				timer.Reset(o.maxIdleTime)
				o.timeoutGroup.add(id)
			}
			o.mux.Unlock()

		}
	}
}

func (o *OnDemandBlockTaskPool) Start() error {
	for {
		if atomic.LoadInt32(&o.state) == stateRunning {
			return errTaskPoolIsStarted
		}

		if atomic.LoadInt32(&o.state) == stateClosing {
			return errTaskPoolIsClosing
		}

		if atomic.LoadInt32(&o.state) == stateStopped {
			return errTaskPoolIsStopped
		}

		if atomic.CompareAndSwapInt32(&o.state, stateCreated, stateLocked) {
			n := o.numOfThatCanBeCreate()
			o.increaseTotalGo(n)
			for i := int32(0); i < n; i++ {
				go o.goroutine(atomic.AddInt32(&o.id, 1))
			}
			atomic.CompareAndSwapInt32(&o.state, stateLocked, stateRunning)
			return nil
		}
	}
}

func (o *OnDemandBlockTaskPool) numOfThatCanBeCreate() int32 {
	n := o.initGo
	needGo := int32(len(o.queue)) - o.initGo
	allowGo := o.maxGo - o.initGo
	if needGo > 0 {
		if needGo < allowGo {
			n += needGo
		} else {
			n += allowGo
		}
	}
	return n
}

func (o *OnDemandBlockTaskPool) Shutdown() (<-chan struct{}, error) {
	//拒绝新提交的任务，执行完已提交的任务
	for {
		if atomic.LoadInt32(&o.state) == stateClosing {
			return nil, errTaskPoolIsClosing
		}

		if atomic.LoadInt32(&o.state) == stateStopped {
			return nil, errTaskPoolIsStopped
		}

		if atomic.LoadInt32(&o.state) == stateCreated {
			return nil, errTaskPoolIsNotRunning
		}

		if atomic.CompareAndSwapInt32(&o.state, stateRunning, stateClosing) {
			close(o.queue)
			return o.interruptCtx.Done(), nil
		}
	}
}

func (o *OnDemandBlockTaskPool) ShutdownNow() ([]Task, error) {
	for {
		if atomic.LoadInt32(&o.state) == stateStopped {
			return nil, errTaskPoolIsStopped
		}
		if atomic.LoadInt32(&o.state) == stateClosing {
			return nil, errTaskPoolIsClosing
		}
		if atomic.LoadInt32(&o.state) == stateCreated {
			return nil, errTaskPoolIsNotRunning
		}
		if atomic.CompareAndSwapInt32(&o.state, stateRunning, stateStopped) {

			close(o.queue)
			o.interruptCancel()

			tasks := make([]Task, 0, len(o.queue))
			for task := range o.queue {
				tasks = append(tasks, task)
			}
			return tasks, nil
		}
	}
}

func (o *OnDemandBlockTaskPool) States(ctx context.Context, interval time.Duration) (<-chan State, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if o.interruptCtx.Err() != nil {
		return nil, o.interruptCtx.Err()
	}
	stateChan := make(chan State)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				o.sendState(stateChan, time.Now().Unix())
				close(stateChan)
				return
			case <-o.interruptCtx.Done():
				o.sendState(stateChan, time.Now().Unix())
				close(stateChan)
				return
			case timeStamp := <-ticker.C:
				o.sendState(stateChan, timeStamp.Unix())
			}
		}
	}()
	return stateChan, nil
}

func (o *OnDemandBlockTaskPool) sendState(ch chan<- State, timeStamp int64) {
	select {
	case ch <- o.getState(timeStamp):
	default:
	}
}

func (o *OnDemandBlockTaskPool) getState(timeStamp int64) State {
	return State{
		PoolState:      atomic.LoadInt32(&o.state),
		GoCnt:          o.numOfGo(),
		WaitingTaskCnt: len(o.queue),
		QueueSize:      cap(o.queue),
		RunningTaskCnt: atomic.LoadInt32(&o.runningTaskCnt),
		TimesStamp:     timeStamp,
	}
}

func (o *OnDemandBlockTaskPool) internalState() int32 {
	for {
		state := atomic.LoadInt32(&o.state)
		if state != stateLocked {
			return state
		}
	}
}
