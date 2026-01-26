# 项目五：实现一个任务调度器

## 目标
通过实现任务调度器，深入理解：
- 定时任务
- 并发控制
- 任务队列
- 优雅关闭

## 功能要求

### 阶段1：基础调度
- [ ] 支持一次性任务
- [ ] 支持定时任务（Cron 表达式）
- [ ] 支持延迟任务
- [ ] 任务取消

### 阶段2：高级功能
- [ ] 任务优先级
- [ ] 任务重试
- [ ] 任务超时
- [ ] 并发限制

### 阶段3：持久化
- [ ] 任务持久化
- [ ] 任务恢复
- [ ] 分布式锁

## API 设计

```go
// 任务定义
type Job interface {
    Run(ctx context.Context) error
}

type JobFunc func(ctx context.Context) error

func (f JobFunc) Run(ctx context.Context) error {
    return f(ctx)
}

// 调度器
type Scheduler struct {
    // ...
}

func NewScheduler(opts ...Option) *Scheduler

// 添加任务
func (s *Scheduler) AddJob(name string, job Job, schedule Schedule) error
func (s *Scheduler) AddFunc(name string, fn func(ctx context.Context) error, schedule Schedule) error

// 调度类型
func Every(d time.Duration) Schedule           // 每隔一段时间
func Cron(expr string) Schedule                // Cron 表达式
func At(t time.Time) Schedule                  // 指定时间
func After(d time.Duration) Schedule           // 延迟执行

// 使用示例
scheduler := NewScheduler(
    WithWorkers(10),
    WithRetry(3),
)

// 每分钟执行
scheduler.AddFunc("cleanup", cleanup, Every(time.Minute))

// Cron 表达式
scheduler.AddFunc("report", generateReport, Cron("0 9 * * *"))

// 延迟执行
scheduler.AddFunc("notify", sendNotification, After(5*time.Second))

scheduler.Start()
defer scheduler.Stop()
```

## 实现提示

### 时间轮
```go
type TimeWheel struct {
    interval    time.Duration
    slots       int
    currentPos  int
    ticker      *time.Ticker
    buckets     []*list.List
    taskMap     map[string]*list.Element
    mu          sync.Mutex
}

func (tw *TimeWheel) AddTask(delay time.Duration, task *Task) {
    tw.mu.Lock()
    defer tw.mu.Unlock()
    
    pos := (tw.currentPos + int(delay/tw.interval)) % tw.slots
    tw.buckets[pos].PushBack(task)
}
```

### Cron 解析
```go
type CronSchedule struct {
    minute, hour, day, month, weekday []int
}

func (c *CronSchedule) Next(t time.Time) time.Time {
    // 计算下一次执行时间
}
```

### 任务重试
```go
type RetryPolicy struct {
    MaxRetries int
    Delay      time.Duration
    MaxDelay   time.Duration
    Multiplier float64
}

func (s *Scheduler) runWithRetry(ctx context.Context, job Job, policy RetryPolicy) error {
    var err error
    delay := policy.Delay
    
    for i := 0; i <= policy.MaxRetries; i++ {
        err = job.Run(ctx)
        if err == nil {
            return nil
        }
        
        if i < policy.MaxRetries {
            time.Sleep(delay)
            delay = time.Duration(float64(delay) * policy.Multiplier)
            if delay > policy.MaxDelay {
                delay = policy.MaxDelay
            }
        }
    }
    return err
}
```

## 测试场景
1. 定时任务准确性
2. 并发任务执行
3. 任务取消
4. 任务重试
5. 优雅关闭

## 参考资源
- [robfig/cron](https://github.com/robfig/cron)
- [时间轮算法](https://blog.csdn.net/mindfloating/article/details/8033340)

## 学习收获
- 理解定时任务的实现原理
- 掌握时间轮算法
- 学会实现优雅关闭
- 理解任务调度的设计模式
