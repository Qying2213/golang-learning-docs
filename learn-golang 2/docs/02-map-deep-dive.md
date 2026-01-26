# Map 深入理解

## 1. Map 底层结构

```go
// runtime/map.go
type hmap struct {
    count     int    // 元素个数
    flags     uint8  // 状态标志（是否正在写入等）
    B         uint8  // buckets 数量的对数 (buckets = 2^B)
    noverflow uint16 // 溢出桶的近似数量
    hash0     uint32 // 哈希种子
    
    buckets    unsafe.Pointer // 桶数组指针
    oldbuckets unsafe.Pointer // 扩容时的旧桶
    nevacuate  uintptr        // 扩容进度
    extra      *mapextra      // 溢出桶相关
}

// 每个桶的结构
type bmap struct {
    tophash [8]uint8  // 存储 hash 值的高 8 位
    // 后面紧跟 8 个 key 和 8 个 value（编译时确定）
    // keys     [8]keytype
    // values   [8]valuetype
    // overflow *bmap
}
```

**关键点：**
- 每个桶存储 8 个键值对
- 使用 tophash 快速定位
- 溢出时使用链表连接溢出桶

## 2. Map 的创建

```go
// 方式1: make
m1 := make(map[string]int)
m2 := make(map[string]int, 100)  // 预分配容量

// 方式2: 字面量
m3 := map[string]int{
    "a": 1,
    "b": 2,
}

// 注意: 不能使用 new
m4 := new(map[string]int)  // *map[string]int, 指向 nil
// (*m4)["a"] = 1  // panic!
```

## 3. Map 的扩容

**触发条件：**
1. 负载因子 > 6.5（元素数量 / 桶数量）
2. 溢出桶太多

**扩容方式：**
- 等量扩容：溢出桶太多时，重新整理（不增加桶数量）
- 翻倍扩容：负载因子过高时，桶数量翻倍

**渐进式扩容：**
- 不是一次性迁移，而是每次操作时迁移一部分
- 避免一次性迁移造成的延迟

## 4. 并发安全问题（重点！）

```go
// Map 不是并发安全的！
m := make(map[string]int)

// 并发写会 panic
go func() {
    for {
        m["a"] = 1
    }
}()
go func() {
    for {
        m["b"] = 2
    }
}()
// fatal error: concurrent map writes

// 并发读写也会 panic
go func() {
    for {
        m["a"] = 1
    }
}()
go func() {
    for {
        _ = m["a"]
    }
}()
// fatal error: concurrent map read and map write
```

## 5. 解决并发问题

### 方案1: sync.Mutex

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    val, ok := sm.m[key]
    return val, ok
}

func (sm *SafeMap) Set(key string, val int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m[key] = val
}

func (sm *SafeMap) Delete(key string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.m, key)
}
```

### 方案2: sync.Map（官方并发安全 Map）

```go
var m sync.Map

// 存储
m.Store("key", "value")

// 读取
val, ok := m.Load("key")

// 删除
m.Delete("key")

// 读取或存储（原子操作）
actual, loaded := m.LoadOrStore("key", "value")

// 读取并删除
val, loaded := m.LoadAndDelete("key")

// 遍历
m.Range(func(key, value interface{}) bool {
    fmt.Println(key, value)
    return true  // 返回 false 停止遍历
})
```

### sync.Map 适用场景

```go
// 适合的场景：
// 1. 读多写少
// 2. 多个 goroutine 读写不同的 key

// 不适合的场景：
// 1. 写多读少
// 2. 频繁更新同一个 key
// 这些场景用 RWMutex + map 更好
```

## 6. Map 的遍历

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// 遍历是无序的！每次顺序可能不同
for k, v := range m {
    fmt.Println(k, v)
}

// 如果需要有序遍历，先排序 key
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

## 7. Map 作为函数参数

```go
// Map 是引用类型，传递的是指针
func modifyMap(m map[string]int) {
    m["new"] = 100  // 会影响原 map
}

func main() {
    m := map[string]int{"a": 1}
    modifyMap(m)
    fmt.Println(m)  // map[a:1 new:100]
}
```

## 8. 常见陷阱

```go
// 陷阱1: nil map 可以读，不能写
var m map[string]int
_ = m["a"]     // OK, 返回零值
m["a"] = 1     // panic!

// 陷阱2: 不能对 map 元素取地址
m := map[string]int{"a": 1}
// p := &m["a"]  // 编译错误！
// 因为 map 扩容时地址会变

// 陷阱3: 遍历时删除
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    if k == "a" {
        delete(m, k)  // 这是安全的
    }
}

// 陷阱4: 遍历时添加
// 新添加的元素可能会也可能不会被遍历到
```

## 9. 性能优化

```go
// 1. 预分配容量
m := make(map[string]int, 1000)

// 2. 使用合适的 key 类型
// string 作为 key 需要计算 hash，较慢
// int 作为 key 更快

// 3. 避免频繁的小 map 创建
// 考虑使用 sync.Pool 复用

// 4. 大 map 考虑分片
type ShardedMap struct {
    shards [256]struct {
        sync.RWMutex
        m map[string]int
    }
}

func (sm *ShardedMap) getShard(key string) *struct {
    sync.RWMutex
    m map[string]int
} {
    h := fnv.New32a()
    h.Write([]byte(key))
    return &sm.shards[h.Sum32()%256]
}
```

## 10. 面试题

**Q1: map 可以边遍历边删除吗？**
```go
// 可以删除当前元素，但行为是未定义的
// 官方说法：如果在遍历时删除，删除的元素不会再被遍历到
// 但如果添加，新元素可能会也可能不会被遍历到
```

**Q2: 为什么 map 的 key 必须是可比较的？**
```go
// 因为需要判断 key 是否相等
// 不能作为 key 的类型：slice, map, function
// 可以作为 key：int, string, pointer, struct(所有字段可比较), array
```

**Q3: sync.Map 的实现原理？**
```go
// 使用两个 map：read 和 dirty
// read: 无锁读取，存储经常访问的数据
// dirty: 需要加锁，存储新写入的数据
// 当 read ��� miss 次数过多，dirty 提升为 read
```

## 练习

1. 实现一个带过期时间的缓存 Map
2. 实现一个 LRU Cache
3. 实现一个分片的并发安全 Map
4. 对比 sync.Map 和 RWMutex+map 在不同场景下的性能
