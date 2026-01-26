# 项目二：实现一个通用连接池

## 目标
通过实现连接池，深入理解：
- 资源管理
- 并发控制
- 超时处理
- 健康检查

## 功能要求

### 阶段1：基础连接池
- [ ] 支持创建连接
- [ ] 支持获取连接（Get）
- [ ] 支持归还连接（Put）
- [ ] 支持关闭连接池

### 阶段2：高级功能
- [ ] 最大连接数限制
- [ ] 最小空闲连接数
- [ ] 连接超时
- [ ] 获取连接超时
- [ ] 空闲连接超时回收

### 阶段3：健康检查
- [ ] 连接健康检查
- [ ] 自动重连
- [ ] 连接预热

## API 设计

```go
// 连接接口
type Conn interface {
    Close() error
    IsAlive() bool
}

// 连接工厂
type Factory func() (Conn, error)

// 连接池配置
type Config struct {
    InitialSize  int           // 初始连接数
    MaxSize      int           // 最大连接数
    MaxIdleSize  int           // 最大空闲连接数
    IdleTimeout  time.Duration // 空闲超时
    WaitTimeout  time.Duration // 获取连接超时
}

// 连接池
type Pool struct {
    // ...
}

func NewPool(factory Factory, config Config) (*Pool, error)
func (p *Pool) Get(ctx context.Context) (Conn, error)
func (p *Pool) Put(conn Conn)
func (p *Pool) Close()
func (p *Pool) Stats() Stats

// 使用示例
pool, _ := NewPool(func() (Conn, error) {
    return net.Dial("tcp", "localhost:6379")
}, Config{
    InitialSize: 5,
    MaxSize:     20,
    IdleTimeout: 5 * time.Minute,
})

conn, _ := pool.Get(ctx)
defer pool.Put(conn)
// 使用连接...
```

## 实现提示

### 核心结构
```go
type Pool struct {
    mu          sync.Mutex
    factory     Factory
    config      Config
    
    idle        chan Conn     // 空闲连接
    active      int           // 活跃连接数
    closed      bool
    
    cond        *sync.Cond    // 等待可用连接
}
```

### 获取连接逻辑
```go
func (p *Pool) Get(ctx context.Context) (Conn, error) {
    // 1. 尝试从空闲池获取
    // 2. 如果没有空闲，且未达到最大连接数，创建新连接
    // 3. 如果达到最大连接数，等待或超时
    // 4. 检查连接健康状态
}
```

## 测试场景
1. 并发获取连接
2. 连接超时
3. 连接池耗尽
4. 连接健康检查失败
5. 优雅关闭

## 参考资源
- [database/sql 连接池实现](https://golang.org/src/database/sql/sql.go)
- [go-redis 连接池](https://github.com/go-redis/redis)

## 学习收获
- 理解连接池的设计模式
- 掌握并发资源管理
- 学会处理超时和取消
- 理解健康检查的重要性
