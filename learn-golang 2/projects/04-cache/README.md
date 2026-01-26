# 项目四：实现一个分布式缓存

## 目标
通过实现分布式缓存，深入理解：
- 缓存淘汰策略（LRU）
- 并发安全
- 一致性哈希
- 防止缓存击穿

## 功能要求

### 阶段1：单机缓存
- [ ] 实现 LRU 缓存淘汰
- [ ] 支持并发安全
- [ ] 支持过期时间
- [ ] 支持内存限制

### 阶段2：HTTP 服务
- [ ] 提供 HTTP API
- [ ] 支持 GET/SET/DELETE
- [ ] 支持批量操作

### 阶段3：分布式
- [ ] 一致性哈希
- [ ] 节点发现
- [ ] 请求转发
- [ ] 防止缓存击穿（singleflight）

## API 设计

```go
// 缓存接口
type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, value []byte, ttl time.Duration)
    Delete(key string)
    Stats() Stats
}

// LRU 缓存
type LRUCache struct {
    maxBytes  int64
    usedBytes int64
    ll        *list.List
    cache     map[string]*list.Element
    mu        sync.RWMutex
}

// 使用示例
cache := NewLRUCache(64 * 1024 * 1024) // 64MB
cache.Set("key", []byte("value"), time.Hour)
value, ok := cache.Get("key")
```

## 实现提示

### LRU 实现
```go
type entry struct {
    key       string
    value     []byte
    expireAt  time.Time
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if ele, ok := c.cache[key]; ok {
        // 检查过期
        e := ele.Value.(*entry)
        if !e.expireAt.IsZero() && time.Now().After(e.expireAt) {
            c.removeElement(ele)
            return nil, false
        }
        // 移到队首
        c.ll.MoveToFront(ele)
        return e.value, true
    }
    return nil, false
}
```

### 一致性哈希
```go
type ConsistentHash struct {
    hash     func([]byte) uint32
    replicas int              // 虚拟节点数
    keys     []int            // 排序的哈希值
    hashMap  map[int]string   // 哈希值 -> 节点
}

func (h *ConsistentHash) Get(key string) string {
    hash := int(h.hash([]byte(key)))
    idx := sort.Search(len(h.keys), func(i int) bool {
        return h.keys[i] >= hash
    })
    return h.hashMap[h.keys[idx%len(h.keys)]]
}
```

### 防止缓存击穿
```go
import "golang.org/x/sync/singleflight"

var g singleflight.Group

func (c *Cache) GetOrLoad(key string, loader func() ([]byte, error)) ([]byte, error) {
    // 先查缓存
    if v, ok := c.Get(key); ok {
        return v, nil
    }
    
    // 使用 singleflight 防止并发请求
    v, err, _ := g.Do(key, func() (interface{}, error) {
        // 再次检查缓存
        if v, ok := c.Get(key); ok {
            return v, nil
        }
        // 加载数据
        data, err := loader()
        if err != nil {
            return nil, err
        }
        c.Set(key, data, time.Hour)
        return data, nil
    })
    
    if err != nil {
        return nil, err
    }
    return v.([]byte), nil
}
```

## 测试场景
1. LRU 淘汰
2. 并发读写
3. 过期清理
4. 一致性哈希分布
5. 缓存击穿防护

## 参考资源
- [groupcache](https://github.com/golang/groupcache)
- [7天用Go从零实现分布式缓存](https://geektutu.com/post/geecache.html)

## 学习收获
- 理解缓存淘汰策略
- 掌握一致性哈希算法
- 学会防止缓存穿透/击穿/雪崩
- 理解分布式系统的基本概念
