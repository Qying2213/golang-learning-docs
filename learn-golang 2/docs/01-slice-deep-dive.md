# Slice 深入理解

## 1. Slice 底层结构

```go
// runtime/slice.go
type slice struct {
    array unsafe.Pointer  // 指向底层数组的指针
    len   int             // 当前长度
    cap   int             // 容量
}
```

Slice 本质是一个**描述符**，包含三个字段：
- `array`: 指向底层数组的指针
- `len`: 切片的长度（可访问的元素数量）
- `cap`: 切片的容量（从切片起始位置到底层数组末尾的长度）

## 2. 创建 Slice 的方式

```go
// 方式1: 字面量
s1 := []int{1, 2, 3}

// 方式2: make
s2 := make([]int, 5)      // len=5, cap=5
s3 := make([]int, 5, 10)  // len=5, cap=10

// 方式3: 从数组切片
arr := [5]int{1, 2, 3, 4, 5}
s4 := arr[1:3]  // len=2, cap=4

// 方式4: 从切片切片
s5 := s1[0:2]
```

## 3. 扩容机制（重点！）

当 append 时容量不足，会触发扩容：

```go
// Go 1.18+ 的扩容策略
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
    newcap := oldCap
    doublecap := newcap + newcap
    
    if newLen > doublecap {
        newcap = newLen
    } else {
        const threshold = 256
        if oldCap < threshold {
            newcap = doublecap  // 小于256，直接翻倍
        } else {
            // 大于256，每次增长 (oldcap + 3*threshold) / 4
            // 约等于增长 25% + 192
            for 0 < newcap && newcap < newLen {
                newcap += (newcap + 3*threshold) / 4
            }
        }
    }
    // ... 内存对齐处理
}
```

**扩容规则总结：**
- 如果需要的容量 > 2倍旧容量，直接使用需要的容量
- 如果旧容量 < 256，新容量 = 旧容量 * 2
- 如果旧容量 >= 256，新容量 = 旧容量 * 1.25 + 192（大约）
- 最后会进行内存对齐，实际容量可能更大

## 4. 共享底层数组的坑（必须掌握！）

```go
func main() {
    // 坑1: 切片共享底层数组
    original := []int{1, 2, 3, 4, 5}
    slice1 := original[1:3]  // [2, 3]
    slice1[0] = 100          // 修改 slice1
    
    fmt.Println(original)    // [1, 100, 3, 4, 5] 原切片也被修改了！
    
    // 坑2: append 可能影响原切片
    original2 := []int{1, 2, 3, 4, 5}
    slice2 := original2[1:3]  // len=2, cap=4
    slice2 = append(slice2, 999)
    
    fmt.Println(original2)    // [1, 2, 3, 999, 5] 原切片被修改！
    
    // 坑3: append 扩容后不再共享
    original3 := []int{1, 2, 3}
    slice3 := original3[0:2]  // len=2, cap=3
    slice3 = append(slice3, 4, 5, 6)  // 触发扩容
    slice3[0] = 100
    
    fmt.Println(original3)    // [1, 2, 3] 原切片不受影响
}
```

## 5. 安全复制切片

```go
// 方法1: copy 函数
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)

// 方法2: append 到空切片
dst2 := append([]int(nil), src...)

// 方法3: 完整切片表达式（限制容量）
original := []int{1, 2, 3, 4, 5}
// [low:high:max] max 限制新切片的容量
safe := original[1:3:3]  // len=2, cap=2
// 此时 append 会触发扩容，不会影响原切片
```

## 6. 切片作为函数参数

```go
// 切片是引用类型，但传递的是切片头的副本
func modifySlice(s []int) {
    s[0] = 100      // 会影响原切片
    s = append(s, 4) // 不会影响原切片（除非扩容）
}

// 如果需要修改切片本身，传指针
func appendToSlice(s *[]int, val int) {
    *s = append(*s, val)
}

// 或者返回新切片
func appendAndReturn(s []int, val int) []int {
    return append(s, val)
}
```

## 7. 性能优化技巧

```go
// 1. 预分配容量，避免多次扩容
data := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    data = append(data, i)
}

// 2. 重用切片
data = data[:0]  // 清空但保留容量

// 3. 避免内存泄漏
func getFirst(s []int) []int {
    // 不好：返回的切片仍引用原大数组
    return s[:1]
    
    // 好：复制出来
    result := make([]int, 1)
    copy(result, s[:1])
    return result
}
```

## 8. 常见面试题

**Q1: nil 切片和空切片的区别？**
```go
var s1 []int          // nil 切片, s1 == nil, len=0, cap=0
s2 := []int{}         // 空切片, s2 != nil, len=0, cap=0
s3 := make([]int, 0)  // 空切片, s3 != nil, len=0, cap=0

// JSON 序列化区别
json.Marshal(s1)  // null
json.Marshal(s2)  // []
```

**Q2: 下面代码输出什么？**
```go
s := []int{1, 2, 3}
for i, v := range s {
    s = append(s, v)
    fmt.Print(i, " ")
}
// 输出: 0 1 2
// range 在开始时就确定了迭代次数
```

## 练习

1. 实现一个 `Reverse` 函数，原地反转切片
2. 实现一个 `Remove` 函数，删除指定索引的元素
3. 分析以下代码的内存分配次数：
```go
var s []int
for i := 0; i < 1000; i++ {
    s = append(s, i)
}
```



