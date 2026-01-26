# Go Map 完全教程

> 从零开始理解 map，适合实习生

---

## 一、什么是 Map

**Map 就是键值对的集合，类似于字典。底层是哈希表。**

```go
// key -> value
// "name" -> "秦阳"
// "age" -> 22
```

---

## 二、创建 Map

### 2.1 用 make 创建

```go
// make(map[键类型]值类型)
m := make(map[string]int)

m["apple"] = 5
m["banana"] = 3

fmt.Println(m)  // map[apple:5 banana:3]
```

### 2.2 字面量创建

```go
m := map[string]int{
    "apple":  5,
    "banana": 3,
    "orange": 8,
}
```

### 2.3 注意：不能用 nil map

```go
var m map[string]int  // m 是 nil

m["apple"] = 5  // ❌ panic! 不能往 nil map 写入

// 必须先初始化
m = make(map[string]int)
m["apple"] = 5  // ✅ OK
```

---

## 三、基本操作

### 3.1 增加/修改

```go
m := make(map[string]int)

// 增加
m["apple"] = 5

// 修改（key 已存在就是修改）
m["apple"] = 10
```

### 3.2 查询

```go
m := map[string]int{"apple": 5, "banana": 3}

// 方式1：直接取值
v := m["apple"]
fmt.Println(v)  // 5

// 如果 key 不存在，返回零值
v2 := m["orange"]
fmt.Println(v2)  // 0（int 的零值）

// 方式2：判断 key 是否存在（推荐）
v, ok := m["apple"]
if ok {
    fmt.Println("存在:", v)
} else {
    fmt.Println("不存在")
}

// 简写
if v, ok := m["apple"]; ok {
    fmt.Println("存在:", v)
}
```

### 3.3 删除

```go
m := map[string]int{"apple": 5, "banana": 3}

delete(m, "apple")

fmt.Println(m)  // map[banana:3]

// 删除不存在的 key 不会报错
delete(m, "orange")  // 什么都不会发生
```

### 3.4 获取长度

```go
m := map[string]int{"apple": 5, "banana": 3}

fmt.Println(len(m))  // 2
```

---

## 四、遍历 Map

```go
m := map[string]int{
    "apple":  5,
    "banana": 3,
    "orange": 8,
}

// 遍历 key 和 value
for k, v := range m {
    fmt.Println(k, v)
}

// 只遍历 key
for k := range m {
    fmt.Println(k)
}

// 只遍历 value
for _, v := range m {
    fmt.Println(v)
}
```

**⚠️ 重要：Map 的遍历顺序是随机的！**

```go
// 每次运行顺序可能不同
for k, v := range m {
    fmt.Println(k, v)
}
// 可能输出：
// apple 5
// banana 3
// orange 8

// 也可能输出：
// banana 3
// orange 8
// apple 5
```

---

## 五、Map 的 Key 类型

**Key 必须是可比较的类型：**

```go
// ✅ 可以作为 key
map[int]string       // int
map[string]int       // string
map[float64]int      // float64
map[bool]string      // bool
map[[3]int]string    // 数组

// ❌ 不能作为 key
map[[]int]string     // 切片不行
map[map[string]int]string  // map 不行
map[func()]string    // 函数不行
```

---

## 六、Map 的 Value 类型

**Value 可以是任何类型：**

```go
// 值是切片
m1 := map[string][]int{
    "scores": {90, 85, 88},
}

// 值是结构体
type Person struct {
    Name string
    Age  int
}
m2 := map[string]Person{
    "p1": {"秦阳", 22},
}

// 值是 map（嵌套）
m3 := map[string]map[string]int{
    "class1": {"math": 90, "english": 85},
}
```

---

## 七、Map 是引用类型

```go
m1 := map[string]int{"apple": 5}
m2 := m1  // m2 和 m1 指向同一个 map

m2["apple"] = 100

fmt.Println(m1["apple"])  // 100（m1 也变了！）
```

**函数传参也是一样：**

```go
func modify(m map[string]int) {
    m["apple"] = 100
}

func main() {
    m := map[string]int{"apple": 5}
    modify(m)
    fmt.Println(m["apple"])  // 100（被修改了）
}
```

---

## 八、Map 不是并发安全的（重要！）

```go
m := make(map[string]int)

// 多个 goroutine 同时读写会 panic
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
```

**解决方案：加锁**

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Get(key string) int {
    sm.mu.RLock()         // 读锁
    defer sm.mu.RUnlock()
    return sm.m[key]
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()          // 写锁
    defer sm.mu.Unlock()
    sm.m[key] = value
}
```

---

## 九、常见用法

### 9.1 统计字符出现次数

```go
func CountChars(s string) map[rune]int {
    result := make(map[rune]int)
    for _, c := range s {
        result[c]++
    }
    return result
}

func main() {
    fmt.Println(CountChars("hello"))
    // map[e:1 h:1 l:2 o:1]
}
```

### 9.2 去重

```go
func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            result = append(result, num)
            seen[num] = true
        }
    }
    return result
}

func main() {
    fmt.Println(RemoveDuplicates([]int{1, 2, 2, 3, 3, 3}))
    // [1 2 3]
}
```

### 9.3 两数之和（经典面试题）

```go
// 找出数组中和为 target 的两个数的下标
func TwoSum(nums []int, target int) []int {
    m := make(map[int]int)  // 值 -> 下标
    
    for i, num := range nums {
        // 需要的另一个数
        need := target - num
        
        // 看看 map 里有没有
        if j, ok := m[need]; ok {
            return []int{j, i}
        }
        
        // 记录当前数和下标
        m[num] = i
    }
    return nil
}

func main() {
    fmt.Println(TwoSum([]int{2, 7, 11, 15}, 9))  // [0 1]
}
```

### 9.4 用 map 实现集合（Set）

```go
// Go 没有内置 Set，用 map 模拟
set := make(map[string]bool)

// 添加
set["apple"] = true
set["banana"] = true

// 判断存在
if set["apple"] {
    fmt.Println("apple 存在")
}

// 删除
delete(set, "apple")
```

---

## 十、Map 的坑

### 10.1 nil map 可以读，不能写

```go
var m map[string]int  // nil

// 读：返回零值，不报错
v := m["apple"]
fmt.Println(v)  // 0

// 写：panic!
m["apple"] = 5  // panic: assignment to entry in nil map
```

### 10.2 不能对 map 元素取地址

```go
m := map[string]int{"apple": 5}

p := &m["apple"]  // ❌ 编译错误！

// 因为 map 扩容时，元素地址会变
```

### 10.3 遍历时删除是安全的

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

for k := range m {
    if k == "a" {
        delete(m, k)  // ✅ 安全
    }
}
```

### 10.4 遍历时添加是不确定的

```go
m := map[string]int{"a": 1}

for k := range m {
    m["b"] = 2  // 新加的 "b" 可能会被遍历到，也可能不会
    fmt.Println(k)
}
```

---

## 十一、练习题

### 练习1：统计单词出现次数 ⭐

```go
// 统计字符串中每个单词出现的次数
func WordCount(s string) map[string]int {
    // 你来实现
    // 提示：用 strings.Fields(s) 分割单词
}

func main() {
    s := "hello world hello go go go"
    fmt.Println(WordCount(s))
    // map[go:3 hello:2 world:1]
}
```

---

### 练习2：判断两个字符串是否是字母异位词 ⭐⭐

```go
// 字母异位词：两个字符串包含相同的字母，只是顺序不同
// 例如："listen" 和 "silent" 是字母异位词

func IsAnagram(s1, s2 string) bool {
    // 你来实现
    // 提示：统计两个字符串的字母出现次数，比较是否相同
}

func main() {
    fmt.Println(IsAnagram("listen", "silent"))  // true
    fmt.Println(IsAnagram("hello", "world"))    // false
}
```

---

### 练习3：合并两个 Map ⭐

```go
// 合并两个 map，如果 key 相同，用 m2 的值覆盖
func MergeMaps(m1, m2 map[string]int) map[string]int {
    // 你来实现
}

func main() {
    m1 := map[string]int{"a": 1, "b": 2}
    m2 := map[string]int{"b": 3, "c": 4}
    fmt.Println(MergeMaps(m1, m2))
    // map[a:1 b:3 c:4]
}
```

---

### 练习4：反转 Map ⭐⭐

```go
// 把 map 的 key 和 value 互换
func ReverseMap(m map[string]int) map[int]string {
    // 你来实现
}

func main() {
    m := map[string]int{"a": 1, "b": 2, "c": 3}
    fmt.Println(ReverseMap(m))
    // map[1:a 2:b 3:c]
}
```

---

## 十二、Map 底层原理

### 12.1 整体结构

```
Map 底层结构：

┌─────────────────────────────────────────────────────────────┐
│                          hmap                               │
├─────────────────────────────────────────────────────────────┤
│  count     = 5          // 元素个数                          │
│  B         = 2          // 桶数量 = 2^B = 4 个桶             │
│  hash0     = 12345      // 哈希种子（随机）                   │
│  buckets   ──────────────────┐                              │
│  oldbuckets = nil            │  // 扩容时用                  │
└──────────────────────────────│──────────────────────────────┘
                               │
                               ▼
              ┌────────┬────────┬────────┬────────┐
              │ 桶0    │ 桶1    │ 桶2    │ 桶3    │
              │ bmap   │ bmap   │ bmap   │ bmap   │
              └────────┴────────┴────────┴────────┘
```

### 12.2 hmap 结构

```go
// runtime/map.go
type hmap struct {
    count     int            // 元素个数，len(m) 返回这个值
    flags     uint8          // 状态标志（是否正在写入等）
    B         uint8          // 桶数量的对数，桶数量 = 2^B
    noverflow uint16         // 溢出桶的大概数量
    hash0     uint32         // 哈希种子，每个 map 不同
    
    buckets    unsafe.Pointer // 桶数组的指针
    oldbuckets unsafe.Pointer // 扩容时，指向旧桶
    nevacuate  uintptr        // 扩容进度
    extra      *mapextra      // 溢出桶相关
}
```

**关键字段解释：**

| 字段 | 作用 |
|------|------|
| count | 元素个数 |
| B | 桶数量 = 2^B，比如 B=3 就是 8 个桶 |
| hash0 | 哈希种子，让每个 map 的哈希结果不同 |
| buckets | 指向桶数组 |
| oldbuckets | 扩容时指向旧桶 |

### 12.3 bmap 结构（桶）

```go
// 每个桶的结构
type bmap struct {
    tophash [8]uint8  // 存储 8 个 key 的哈希值高 8 位
    
    // 下面这些字段在编译时才确定，源码里看不到
    // keys     [8]keytype    // 8 个 key
    // values   [8]valuetype  // 8 个 value
    // overflow *bmap         // 溢出桶指针
}
```

**一个桶可以存 8 个键值对：**

```
┌─────────────────────────────────────────────────────────────┐
│                          bmap（桶）                          │
├─────────────────────────────────────────────────────────────┤
│  tophash  [t0][t1][t2][t3][t4][t5][t6][t7]  // 8个哈希高8位  │
├─────────────────────────────────────────────────────────────┤
│  keys     [k0][k1][k2][k3][k4][k5][k6][k7]  // 8个key        │
├─────────────────────────────────────────────────────────────┤
│  values   [v0][v1][v2][v3][v4][v5][v6][v7]  // 8个value      │
├─────────────────────────────────────────────────────────────┤
│  overflow  ──────► 指向溢出桶（如果8个位置满了）              │
└─────────────────────────────────────────────────────────────┘
```

### 12.4 为什么 tophash 存高 8 位？

**快速过滤，减少比较次数。**

```
查找 key 的过程：

1. 计算 key 的哈希值
   hash("apple") = 0x1234567890ABCDEF
                          ↑↑
                        高8位 = 0x12

2. 用哈希值低位确定桶
   桶索引 = hash & (2^B - 1)

3. 在桶里用 tophash 快速过滤
   遍历 tophash[0..7]，找到 == 0x12 的位置
   
4. tophash 匹配了，再比较完整的 key
   如果 key 也相等，找到了！
```

**好处：** 比较 1 字节的 tophash 比比较完整 key 快得多。

### 12.5 哈希冲突怎么办？

**用链表法：桶满了就挂溢出桶。**

```
桶0 装满了 8 个键值对
    │
    ▼
┌────────┐     ┌────────┐     ┌────────┐
│ 桶0    │────►│ 溢出桶1 │────►│ 溢出桶2 │
│ 8个KV  │     │ 8个KV  │     │ 8个KV  │
└────────┘     └────────┘     └────────┘
```

### 12.6 扩容机制

**触发条件：**
1. 负载因子 > 6.5（元素数量 / 桶数量）
2. 溢出桶太多

**扩容方式：**
- **翻倍扩容：** 负载因子过高时，桶数量 × 2
- **等量扩容：** 溢出桶太多时，重新整理（不增加桶数量）

**渐进式扩容：**
- 不是一次性迁移所有数据
- 每次读写操作时，迁移一部分
- 避免一次性迁移造成卡顿

```
扩容过程：

旧桶（4个）                    新桶（8个）
┌──┬──┬──┬──┐                ┌──┬──┬──┬──┬──┬──┬──┬──┐
│0 │1 │2 │3 │  ──迁移──►     │0 │1 │2 │3 │4 │5 │6 │7 │
└──┴──┴──┴──┘                └──┴──┴──┴──┴──┴──┴──┴──┘

每次操作迁移 1-2 个旧桶的数据
```

### 12.7 为什么遍历顺序随机？

**Go 故意的！**

```go
// 遍历开始时，随机选一个起始桶和起始位置
// 这样每次遍历顺序都不同
```

**为什么这样设计？**
- 防止程序员依赖遍历顺序
- 因为扩容后顺序会变，依赖顺序的代码会出 bug

### 12.8 为什么不是并发安全的？

**没有加锁，多个 goroutine 同时写会破坏数据结构。**

```go
// 写入时会检查 flags
if h.flags & hashWriting != 0 {
    throw("concurrent map writes")
}
```

Go 在运行时检测到并发写，直接 panic，而不是返回错误数据。

### 12.9 底层原理总结

```
┌─────────────────────────────────────────────────────────────┐
│                      Map 底层要点                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 底层是哈希表，由 hmap + bmap（桶）组成                    │
│                                                             │
│  2. 每个桶存 8 个键值对                                      │
│                                                             │
│  3. 用 tophash（哈希高8位）快速过滤                          │
│                                                             │
│  4. 哈希冲突用链表法（溢出桶）                               │
│                                                             │
│  5. 扩容是渐进式的，每次操作迁移一部分                        │
│                                                             │
│  6. 遍历顺序故意随机，防止依赖顺序                           │
│                                                             │
│  7. 不是并发安全的，并发写会 panic                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 十三、总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Map 要点                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. 创建：make(map[K]V) 或 map[K]V{...}                                    │
│                                                                             │
│   2. 操作：m[k] = v（增/改）、v := m[k]（查）、delete(m, k)（删）            │
│                                                                             │
│   3. 判断存在：v, ok := m[k]                                                 │
│                                                                             │
│   4. 遍历顺序是随机的                                                        │
│                                                                             │
│   5. nil map 可读不可写                                                      │
│                                                                             │
│   6. 不是并发安全的，多 goroutine 要加锁                                     │
│                                                                             │
│   7. 是引用类型，赋值和传参不会复制                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```
