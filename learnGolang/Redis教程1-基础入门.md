# Redis 教程1 - 基础入门

## Redis 是什么？

**Redis = Remote Dictionary Server = 远程字典服务**

一句话：**内存数据库，速度极快**

```
MySQL：数据存在磁盘，查询要读磁盘 → 慢（毫秒级）
Redis：数据存在内存，直接读内存 → 快（微秒级）

Redis 比 MySQL 快 100-1000 倍！
```

---

## Redis 能做什么？

| 场景 | 说明 | 例子 |
|------|------|------|
| **缓存** | 把热点数据放内存 | 用户信息、商品详情 |
| **Session** | 存储登录状态 | 用户登录后的 Token |
| **排行榜** | 有序集合天然支持 | 游戏积分榜、热搜榜 |
| **计数器** | 原子操作，高并发安全 | 点赞数、阅读量 |
| **分布式锁** | 防止并发冲突 | 秒杀、抢购 |
| **消息队列** | 简单的队列功能 | 异步任务 |

---

## 第一步：安装 Redis

### Mac
```bash
brew install redis

# 启动
brew services start redis

# 或者前台启动
redis-server
```

### Ubuntu
```bash
sudo apt update
sudo apt install redis-server

# 启动
sudo systemctl start redis
sudo systemctl enable redis
```

### Docker（推荐）
```bash
docker run -d --name redis -p 6379:6379 redis
```

### 验证安装
```bash
redis-cli ping
# 返回 PONG 就成功了
```

---

## 第二步：Redis 命令行基础

### 连接 Redis
```bash
redis-cli
```

### 最基本的操作
```bash
# 设置值
SET name "秦阳"

# 获取值
GET name
# "秦阳"

# 删除
DEL name

# 查看是否存在
EXISTS name
# 0 (不存在)
```

---

## 第三步：5 种数据类型

### 1. String（字符串）

**最基础的类型，可以存字符串、数字、JSON**

```bash
# 设置
SET name "秦阳"
SET age 22

# 获取
GET name    # "秦阳"
GET age     # "22"

# 设置过期时间（秒）
SET token "abc123" EX 3600    # 3600秒后过期
SETEX token 3600 "abc123"     # 同上

# 查看剩余时间
TTL token   # 返回剩余秒数，-1 表示永不过期，-2 表示已过期

# 数字操作
SET count 0
INCR count      # 1（自增1）
INCRBY count 10 # 11（增加10）
DECR count      # 10（自减1）
DECRBY count 5  # 5（减少5）

# 追加字符串
APPEND name "同学"
GET name    # "秦阳同学"

# 获取长度
STRLEN name # 12（UTF-8 中文3字节）
```

**练习1：**
```bash
# 实现一个阅读量计数器
SET article:1:views 0
INCR article:1:views    # 每次访问 +1
GET article:1:views     # 查看阅读量
```

---

### 2. Hash（哈希）

**存储对象，类似 Go 的 map[string]string**

```bash
# 设置单个字段
HSET user:1 name "秦阳"
HSET user:1 age 22
HSET user:1 email "qy@example.com"

# 批量设置
HMSET user:2 name "张三" age 25 email "zs@example.com"

# 获取单个字段
HGET user:1 name    # "秦阳"

# 获取多个字段
HMGET user:1 name age   # "秦阳" "22"

# 获取所有字段和值
HGETALL user:1
# name "秦阳"
# age "22"
# email "qy@example.com"

# 获取所有字段名
HKEYS user:1    # name age email

# 获取所有值
HVALS user:1    # "秦阳" "22" "qy@example.com"

# 判断字段是否存在
HEXISTS user:1 name     # 1（存在）
HEXISTS user:1 phone    # 0（不存在）

# 删除字段
HDEL user:1 email

# 字段数量
HLEN user:1     # 2

# 数字字段增加
HINCRBY user:1 age 1    # age 变成 23
```

**练习2：**
```bash
# 存储购物车：用户1 的购物车
HSET cart:1 "iPhone" 2      # iPhone 2件
HSET cart:1 "MacBook" 1     # MacBook 1件
HINCRBY cart:1 "iPhone" 1   # iPhone 再加1件
HGETALL cart:1              # 查看购物车
```

---

### 3. List（列表）

**有序列表，可以从两端操作，像队列或栈**

```bash
# 从左边插入（头部）
LPUSH messages "msg1"
LPUSH messages "msg2"
LPUSH messages "msg3"
# 列表现在是：msg3 -> msg2 -> msg1

# 从右边插入（尾部）
RPUSH messages "msg4"
# 列表现在是：msg3 -> msg2 -> msg1 -> msg4

# 查看列表（0 到 -1 表示全部）
LRANGE messages 0 -1
# msg3 msg2 msg1 msg4

# 查看指定范围
LRANGE messages 0 2     # 前3个：msg3 msg2 msg1

# 从左边弹出
LPOP messages   # 返回 msg3，列表变成 msg2 msg1 msg4

# 从右边弹出
RPOP messages   # 返回 msg4，列表变成 msg2 msg1

# 获取长度
LLEN messages   # 2

# 获取指定位置的元素
LINDEX messages 0   # msg2（第一个）
LINDEX messages -1  # msg1（最后一个）

# 保留指定范围（删除其他）
LTRIM messages 0 99     # 只保留前100个

# 阻塞弹出（消息队列常用）
BLPOP messages 10   # 等待10秒，有数据就返回，没有就超时
```

**练习3：**
```bash
# 实现最新消息列表（只保留最新10条）
LPUSH news "新闻1"
LPUSH news "新闻2"
LPUSH news "新闻3"
LTRIM news 0 9          # 只保留最新10条
LRANGE news 0 -1        # 查看所有新闻
```

---

### 4. Set（集合）

**无序、不重复的集合**

```bash
# 添加元素
SADD tags "golang" "redis" "mysql"
SADD tags "golang"      # 重复添加无效

# 查看所有元素
SMEMBERS tags   # golang redis mysql（顺序随机）

# 元素数量
SCARD tags      # 3

# 判断是否存在
SISMEMBER tags "golang"     # 1（存在）
SISMEMBER tags "java"       # 0（不存在）

# 删除元素
SREM tags "mysql"

# 随机获取元素
SRANDMEMBER tags        # 随机返回一个
SRANDMEMBER tags 2      # 随机返回2个

# 随机弹出（获取并删除）
SPOP tags               # 随机弹出一个

# 集合运算
SADD set1 "a" "b" "c"
SADD set2 "b" "c" "d"

SINTER set1 set2        # 交集：b c
SUNION set1 set2        # 并集：a b c d
SDIFF set1 set2         # 差集（set1 有但 set2 没有）：a
```

**练习4：**
```bash
# 实现共同关注
SADD user:1:follow "张三" "李四" "王五"
SADD user:2:follow "李四" "王五" "赵六"
SINTER user:1:follow user:2:follow  # 共同关注：李四 王五

# 实现抽奖
SADD lottery "用户1" "用户2" "用户3" "用户4" "用户5"
SRANDMEMBER lottery 3   # 随机抽3个中奖者
```

---

### 5. Sorted Set（有序集合）

**带分数的有序集合，按分数排序**

```bash
# 添加（分数 成员）
ZADD leaderboard 100 "player1"
ZADD leaderboard 200 "player2"
ZADD leaderboard 150 "player3"
ZADD leaderboard 180 "player4"

# 查看所有（按分数从低到高）
ZRANGE leaderboard 0 -1
# player1 player3 player4 player2

# 查看所有（按分数从高到低）
ZREVRANGE leaderboard 0 -1
# player2 player4 player3 player1

# 带分数查看
ZREVRANGE leaderboard 0 -1 WITHSCORES
# player2 200
# player4 180
# player3 150
# player1 100

# 查看 Top 3
ZREVRANGE leaderboard 0 2 WITHSCORES

# 获取某个成员的分数
ZSCORE leaderboard "player1"    # 100

# 获取某个成员的排名（从高到低，从0开始）
ZREVRANK leaderboard "player1"  # 3（第4名）
ZREVRANK leaderboard "player2"  # 0（第1名）

# 增加分数
ZINCRBY leaderboard 50 "player1"    # player1 变成 150 分

# 获取分数范围内的成员
ZRANGEBYSCORE leaderboard 100 200   # 分数在100-200之间的

# 删除成员
ZREM leaderboard "player1"

# 成员数量
ZCARD leaderboard   # 3
```

**练习5：**
```bash
# 实现热搜榜
ZADD hotsearch 1000 "话题1"
ZADD hotsearch 2000 "话题2"
ZADD hotsearch 1500 "话题3"

# 每次搜索增加热度
ZINCRBY hotsearch 1 "话题1"

# 获取热搜 Top 10
ZREVRANGE hotsearch 0 9 WITHSCORES
```

---

## 第四步：通用命令

```bash
# 查看所有 key
KEYS *

# 模糊匹配
KEYS user:*     # 所有以 user: 开头的 key

# 查看 key 类型
TYPE name       # string
TYPE user:1     # hash

# 删除 key
DEL name
DEL user:1 user:2   # 删除多个

# 判断 key 是否存在
EXISTS name     # 1 存在，0 不存在

# 设置过期时间
EXPIRE name 60      # 60秒后过期
EXPIREAT name 1735689600    # 指定时间戳过期

# 查看剩余时间
TTL name        # 返回秒数，-1 永不过期，-2 已过期/不存在

# 移除过期时间
PERSIST name

# 重命名
RENAME oldkey newkey

# 选择数据库（默认16个，0-15）
SELECT 1        # 切换到数据库1

# 清空当前数据库
FLUSHDB

# 清空所有数据库
FLUSHALL
```

---

## 第五步：实战练习

### 练习1：用户登录 Session

```bash
# 用户登录，生成 token，存储用户信息
SET session:abc123 '{"userId":1,"username":"秦阳"}' EX 86400

# 验证 token
GET session:abc123

# 用户登出
DEL session:abc123
```

### 练习2：文章点赞

```bash
# 用户1 给文章1 点赞
SADD article:1:likes "user:1"

# 用户2 给文章1 点赞
SADD article:1:likes "user:2"

# 查看点赞数
SCARD article:1:likes   # 2

# 用户1 是否点赞过
SISMEMBER article:1:likes "user:1"  # 1

# 用户1 取消点赞
SREM article:1:likes "user:1"
```

### 练习3：限流（每分钟最多10次请求）

```bash
# 用户请求时
INCR rate:user:1        # 计数+1
EXPIRE rate:user:1 60   # 60秒后重置

# 检查是否超限
GET rate:user:1         # 如果 > 10 就拒绝
```

---

## 总结

| 类型 | 特点 | 常用场景 |
|------|------|---------|
| String | 最基础 | 缓存、计数器、Session |
| Hash | 存对象 | 用户信息、购物车 |
| List | 有序列表 | 消息队列、最新列表 |
| Set | 无序不重复 | 标签、共同好友、抽奖 |
| Sorted Set | 带分数排序 | 排行榜、热搜 |

| 命令 | 作用 |
|------|------|
| SET/GET | 字符串操作 |
| HSET/HGET | 哈希操作 |
| LPUSH/RPOP | 列表操作 |
| SADD/SMEMBERS | 集合操作 |
| ZADD/ZRANGE | 有序集合操作 |
| EXPIRE/TTL | 过期时间 |
| DEL/EXISTS | 删除/判断存在 |

---

## 下一篇

[Redis教程2-Go操作Redis](./Redis教程2-Go操作Redis.md) - 用 Go 代码操作 Redis
