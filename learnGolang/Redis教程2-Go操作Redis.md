# Redis 教程2 - Go 操作 Redis

## 安装 go-redis

```bash
go get github.com/redis/go-redis/v9
```

---

## 第一步：连接 Redis

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
    // 创建 Redis 客户端
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",  // Redis 地址
        Password: "",                 // 密码，没有就空
        DB:       0,                  // 数据库编号
    })

    // 测试连接
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("连接成功:", pong)  // PONG
}
```

---

## 第二步：String 操作

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
    rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}

func main() {
    // ========== 基本操作 ==========
    
    // 设置值
    err := rdb.Set(ctx, "name", "秦阳", 0).Err()  // 0 表示不过期
    if err != nil {
        panic(err)
    }
    
    // 获取值
    val, err := rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("name:", val)  // 秦阳
    
    // ========== 带过期时间 ==========
    
    // 设置 1 小时后过期
    rdb.Set(ctx, "token", "abc123", time.Hour)
    
    // 查看剩余时间
    ttl, _ := rdb.TTL(ctx, "token").Result()
    fmt.Println("剩余时间:", ttl)  // 59m59s
    
    // ========== 数字操作 ==========
    
    // 设置计数器
    rdb.Set(ctx, "views", 0, 0)
    
    // 自增
    newVal, _ := rdb.Incr(ctx, "views").Result()
    fmt.Println("views:", newVal)  // 1
    
    // 增加指定值
    newVal, _ = rdb.IncrBy(ctx, "views", 10).Result()
    fmt.Println("views:", newVal)  // 11
    
    // ========== 判断 key 是否存在 ==========
    
    exists, _ := rdb.Exists(ctx, "name").Result()
    fmt.Println("name 存在:", exists == 1)  // true
    
    // ========== 删除 ==========
    
    rdb.Del(ctx, "name")
    
    // 获取不存在的 key
    val, err = rdb.Get(ctx, "name").Result()
    if err == redis.Nil {
        fmt.Println("name 不存在")
    }
}
```

---

## 第三步：Hash 操作

```go
func hashExample() {
    // ========== 设置 ==========
    
    // 设置单个字段
    rdb.HSet(ctx, "user:1", "name", "秦阳")
    rdb.HSet(ctx, "user:1", "age", 22)
    
    // 批量设置
    rdb.HSet(ctx, "user:2", map[string]interface{}{
        "name":  "张三",
        "age":   25,
        "email": "zs@example.com",
    })
    
    // ========== 获取 ==========
    
    // 获取单个字段
    name, _ := rdb.HGet(ctx, "user:1", "name").Result()
    fmt.Println("name:", name)  // 秦阳
    
    // 获取多个字段
    vals, _ := rdb.HMGet(ctx, "user:1", "name", "age").Result()
    fmt.Println("vals:", vals)  // [秦阳 22]
    
    // 获取所有字段
    all, _ := rdb.HGetAll(ctx, "user:1").Result()
    fmt.Println("all:", all)  // map[name:秦阳 age:22]
    
    // ========== 判断字段是否存在 ==========
    
    exists, _ := rdb.HExists(ctx, "user:1", "name").Result()
    fmt.Println("name 存在:", exists)  // true
    
    // ========== 删除字段 ==========
    
    rdb.HDel(ctx, "user:1", "age")
    
    // ========== 数字字段增加 ==========
    
    rdb.HSet(ctx, "user:1", "score", 100)
    newScore, _ := rdb.HIncrBy(ctx, "user:1", "score", 10).Result()
    fmt.Println("score:", newScore)  // 110
}
```

---

## 第四步：List 操作

```go
func listExample() {
    // 清空之前的数据
    rdb.Del(ctx, "messages")
    
    // ========== 插入 ==========
    
    // 从左边插入
    rdb.LPush(ctx, "messages", "msg1", "msg2", "msg3")
    // 列表：msg3 -> msg2 -> msg1
    
    // 从右边插入
    rdb.RPush(ctx, "messages", "msg4")
    // 列表：msg3 -> msg2 -> msg1 -> msg4
    
    // ========== 获取 ==========
    
    // 获取全部
    all, _ := rdb.LRange(ctx, "messages", 0, -1).Result()
    fmt.Println("all:", all)  // [msg3 msg2 msg1 msg4]
    
    // 获取长度
    length, _ := rdb.LLen(ctx, "messages").Result()
    fmt.Println("length:", length)  // 4
    
    // 获取指定位置
    val, _ := rdb.LIndex(ctx, "messages", 0).Result()
    fmt.Println("first:", val)  // msg3
    
    // ========== 弹出 ==========
    
    // 从左边弹出
    val, _ = rdb.LPop(ctx, "messages").Result()
    fmt.Println("lpop:", val)  // msg3
    
    // 从右边弹出
    val, _ = rdb.RPop(ctx, "messages").Result()
    fmt.Println("rpop:", val)  // msg4
    
    // ========== 阻塞弹出（消息队列常用）==========
    
    // 等待 5 秒，有数据就返回
    result, err := rdb.BLPop(ctx, 5*time.Second, "messages").Result()
    if err == redis.Nil {
        fmt.Println("超时，没有数据")
    } else {
        fmt.Println("blpop:", result)  // [messages msg2]
    }
}
```

---

## 第五步：Set 操作

```go
func setExample() {
    // 清空
    rdb.Del(ctx, "tags", "user:1:follow", "user:2:follow")
    
    // ========== 添加 ==========
    
    rdb.SAdd(ctx, "tags", "golang", "redis", "mysql")
    
    // ========== 获取 ==========
    
    // 获取所有成员
    members, _ := rdb.SMembers(ctx, "tags").Result()
    fmt.Println("members:", members)  // [golang redis mysql]
    
    // 成员数量
    count, _ := rdb.SCard(ctx, "tags").Result()
    fmt.Println("count:", count)  // 3
    
    // ========== 判断是否存在 ==========
    
    exists, _ := rdb.SIsMember(ctx, "tags", "golang").Result()
    fmt.Println("golang 存在:", exists)  // true
    
    // ========== 随机获取 ==========
    
    // 随机获取一个
    random, _ := rdb.SRandMember(ctx, "tags").Result()
    fmt.Println("random:", random)
    
    // 随机获取多个
    randoms, _ := rdb.SRandMemberN(ctx, "tags", 2).Result()
    fmt.Println("randoms:", randoms)
    
    // ========== 集合运算 ==========
    
    rdb.SAdd(ctx, "user:1:follow", "张三", "李四", "王五")
    rdb.SAdd(ctx, "user:2:follow", "李四", "王五", "赵六")
    
    // 交集（共同关注）
    inter, _ := rdb.SInter(ctx, "user:1:follow", "user:2:follow").Result()
    fmt.Println("共同关注:", inter)  // [李四 王五]
    
    // 并集
    union, _ := rdb.SUnion(ctx, "user:1:follow", "user:2:follow").Result()
    fmt.Println("并集:", union)
    
    // 差集
    diff, _ := rdb.SDiff(ctx, "user:1:follow", "user:2:follow").Result()
    fmt.Println("差集:", diff)  // [张三]
}
```

---

## 第六步：Sorted Set 操作

```go
func sortedSetExample() {
    // 清空
    rdb.Del(ctx, "leaderboard")
    
    // ========== 添加 ==========
    
    rdb.ZAdd(ctx, "leaderboard",
        redis.Z{Score: 100, Member: "player1"},
        redis.Z{Score: 200, Member: "player2"},
        redis.Z{Score: 150, Member: "player3"},
        redis.Z{Score: 180, Member: "player4"},
    )
    
    // ========== 获取排行榜 ==========
    
    // 从高到低（Top 3）
    result, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
    fmt.Println("Top 3:")
    for i, z := range result {
        fmt.Printf("  第%d名: %s, 分数: %.0f\n", i+1, z.Member, z.Score)
    }
    // 第1名: player2, 分数: 200
    // 第2名: player4, 分数: 180
    // 第3名: player3, 分数: 150
    
    // ========== 获取排名 ==========
    
    // 获取 player1 的排名（从高到低，从0开始）
    rank, _ := rdb.ZRevRank(ctx, "leaderboard", "player1").Result()
    fmt.Println("player1 排名:", rank+1)  // 4
    
    // ========== 获取分数 ==========
    
    score, _ := rdb.ZScore(ctx, "leaderboard", "player1").Result()
    fmt.Println("player1 分数:", score)  // 100
    
    // ========== 增加分数 ==========
    
    newScore, _ := rdb.ZIncrBy(ctx, "leaderboard", 100, "player1").Result()
    fmt.Println("player1 新分数:", newScore)  // 200
    
    // ========== 获取成员数量 ==========
    
    count, _ := rdb.ZCard(ctx, "leaderboard").Result()
    fmt.Println("成员数量:", count)  // 4
}
```

---

## 第七步：实战案例

### 案例1：缓存用户信息

```go
import (
    "encoding/json"
    "time"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 获取用户（先查缓存，没有再查数据库）
func GetUser(userID int) (*User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    // 1. 先查 Redis
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        // 缓存命中
        var user User
        json.Unmarshal([]byte(val), &user)
        fmt.Println("从缓存获取")
        return &user, nil
    }
    
    // 2. 缓存没有，查数据库（这里模拟）
    user := &User{
        ID:       userID,
        Username: "秦阳",
        Email:    "qy@example.com",
    }
    fmt.Println("从数据库获取")
    
    // 3. 写入缓存（1小时过期）
    data, _ := json.Marshal(user)
    rdb.Set(ctx, key, data, time.Hour)
    
    return user, nil
}

func main() {
    // 第一次：从数据库获取，写入缓存
    user, _ := GetUser(1)
    fmt.Printf("用户: %+v\n", user)
    
    // 第二次：从缓存获取
    user, _ = GetUser(1)
    fmt.Printf("用户: %+v\n", user)
}
```

### 案例2：分布式锁

```go
// 加锁
func Lock(key string, value string, expiration time.Duration) bool {
    // SetNX: 只有 key 不存在时才设置成功
    ok, err := rdb.SetNX(ctx, "lock:"+key, value, expiration).Result()
    if err != nil {
        return false
    }
    return ok
}

// 解锁
func Unlock(key string, value string) bool {
    // 用 Lua 脚本保证原子性：只有值匹配才删除
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    result, err := rdb.Eval(ctx, script, []string{"lock:" + key}, value).Result()
    if err != nil {
        return false
    }
    return result.(int64) == 1
}

func main() {
    lockKey := "order:123"
    lockValue := "unique-id-12345"  // 用唯一ID，防止误删别人的锁
    
    // 尝试加锁
    if Lock(lockKey, lockValue, 10*time.Second) {
        fmt.Println("加锁成功")
        
        // 处理业务...
        time.Sleep(2 * time.Second)
        
        // 解锁
        if Unlock(lockKey, lockValue) {
            fmt.Println("解锁成功")
        }
    } else {
        fmt.Println("加锁失败，有人在处理")
    }
}
```

### 案例3：排行榜

```go
// 增加分数
func AddScore(userID string, score float64) error {
    return rdb.ZIncrBy(ctx, "leaderboard", score, userID).Err()
}

// 获取 Top N
func GetTopN(n int64) ([]redis.Z, error) {
    return rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, n-1).Result()
}

// 获取用户排名
func GetRank(userID string) (int64, error) {
    rank, err := rdb.ZRevRank(ctx, "leaderboard", userID).Result()
    if err != nil {
        return 0, err
    }
    return rank + 1, nil  // 从1开始
}

func main() {
    // 添加分数
    AddScore("user:1", 100)
    AddScore("user:2", 200)
    AddScore("user:3", 150)
    AddScore("user:1", 50)  // user:1 再加 50 分
    
    // 获取 Top 3
    top3, _ := GetTopN(3)
    fmt.Println("排行榜 Top 3:")
    for i, z := range top3 {
        fmt.Printf("  第%d名: %s, 分数: %.0f\n", i+1, z.Member, z.Score)
    }
    
    // 获取 user:1 的排名
    rank, _ := GetRank("user:1")
    fmt.Println("user:1 排名:", rank)
}
```

### 案例4：限流

```go
// 滑动窗口限流：每分钟最多 maxCount 次
func RateLimit(userID string, maxCount int64) bool {
    key := fmt.Sprintf("rate:%s", userID)
    now := time.Now().UnixNano()
    windowStart := now - 60*1e9  // 1分钟前
    
    // 使用 Pipeline 减少网络往返
    pipe := rdb.Pipeline()
    
    // 移除1分钟前的记录
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    
    // 统计当前窗口内的请求数
    countCmd := pipe.ZCard(ctx, key)
    
    // 添加当前请求
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
    
    // 设置过期时间
    pipe.Expire(ctx, key, time.Minute)
    
    pipe.Exec(ctx)
    
    count := countCmd.Val()
    return count < maxCount
}

func main() {
    for i := 0; i < 15; i++ {
        if RateLimit("user:1", 10) {
            fmt.Printf("请求 %d: 允许\n", i+1)
        } else {
            fmt.Printf("请求 %d: 限流\n", i+1)
        }
    }
}
```

---

## 第八步：连接池配置

```go
rdb := redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    Password:     "",
    DB:           0,
    
    // 连接池配置
    PoolSize:     100,              // 连接池大小
    MinIdleConns: 10,               // 最小空闲连接
    MaxRetries:   3,                // 最大重试次数
    
    // 超时配置
    DialTimeout:  5 * time.Second,  // 连接超时
    ReadTimeout:  3 * time.Second,  // 读超时
    WriteTimeout: 3 * time.Second,  // 写超时
})
```

---

## 总结

| 操作 | 命令行 | Go 代码 |
|------|--------|---------|
| 设置字符串 | `SET key value` | `rdb.Set(ctx, key, value, 0)` |
| 获取字符串 | `GET key` | `rdb.Get(ctx, key).Result()` |
| 设置哈希 | `HSET key field value` | `rdb.HSet(ctx, key, field, value)` |
| 获取哈希 | `HGET key field` | `rdb.HGet(ctx, key, field).Result()` |
| 列表左插入 | `LPUSH key value` | `rdb.LPush(ctx, key, value)` |
| 列表获取 | `LRANGE key 0 -1` | `rdb.LRange(ctx, key, 0, -1).Result()` |
| 集合添加 | `SADD key member` | `rdb.SAdd(ctx, key, member)` |
| 集合获取 | `SMEMBERS key` | `rdb.SMembers(ctx, key).Result()` |
| 有序集合添加 | `ZADD key score member` | `rdb.ZAdd(ctx, key, redis.Z{...})` |
| 有序集合排行 | `ZREVRANGE key 0 9` | `rdb.ZRevRange(ctx, key, 0, 9).Result()` |

**常用场景：**
- 缓存：`Set/Get` + JSON 序列化
- 分布式锁：`SetNX` + Lua 脚本
- 排行榜：`ZAdd/ZRevRange`
- 限流：`ZAdd` + 滑动窗口
