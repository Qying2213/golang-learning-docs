# 单元测试与 Mock

## 1. 基本测试

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

// 运行测试
// go test ./...
// go test -v ./...  // 详细输出
// go test -run TestAdd ./...  // 运行特定测试
```

## 2. 表驱动测试（推荐）

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

## 3. 测试辅助函数

```go
func TestSomething(t *testing.T) {
    // t.Helper() 标记辅助函数，错误时显示调用者位置
    assertEqual := func(t *testing.T, got, want int) {
        t.Helper()
        if got != want {
            t.Errorf("got %d; want %d", got, want)
        }
    }
    
    assertEqual(t, Add(1, 2), 3)
    assertEqual(t, Add(0, 0), 0)
}

// t.Cleanup 注册清理函数
func TestWithCleanup(t *testing.T) {
    f := createTempFile(t)
    t.Cleanup(func() {
        os.Remove(f.Name())
    })
    // 测试代码...
}

// t.Parallel 并行测试
func TestParallel(t *testing.T) {
    t.Parallel()  // 标记可以并行运行
    // 测试代码...
}
```


## 4. Mock 和接口

```go
// 定义接口
type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
}

// 生产实现
type MySQLUserRepo struct {
    db *sql.DB
}

func (r *MySQLUserRepo) GetByID(id int) (*User, error) {
    // 真实数据库查询
}

// Mock 实现
type MockUserRepo struct {
    users map[int]*User
    err   error
}

func (m *MockUserRepo) GetByID(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.users[id], nil
}

func (m *MockUserRepo) Save(user *User) error {
    if m.err != nil {
        return m.err
    }
    m.users[user.ID] = user
    return nil
}

// 测试
func TestUserService_GetUser(t *testing.T) {
    mockRepo := &MockUserRepo{
        users: map[int]*User{
            1: {ID: 1, Name: "Alice"},
        },
    }
    
    service := NewUserService(mockRepo)
    
    user, err := service.GetUser(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("got %s; want Alice", user.Name)
    }
}

// 测试错误情况
func TestUserService_GetUser_Error(t *testing.T) {
    mockRepo := &MockUserRepo{
        err: errors.New("database error"),
    }
    
    service := NewUserService(mockRepo)
    
    _, err := service.GetUser(1)
    if err == nil {
        t.Error("expected error, got nil")
    }
}
```

## 5. HTTP 测试

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    // 创建请求
    req := httptest.NewRequest("GET", "/users/1", nil)
    
    // 创建响应记录器
    w := httptest.NewRecorder()
    
    // 调用处理器
    handler := NewUserHandler()
    handler.ServeHTTP(w, req)
    
    // 检查响应
    if w.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
    }
    
    // 检查响应体
    expected := `{"id":1,"name":"Alice"}`
    if w.Body.String() != expected {
        t.Errorf("body = %s; want %s", w.Body.String(), expected)
    }
}

// 测试整个服务器
func TestServer(t *testing.T) {
    server := httptest.NewServer(NewRouter())
    defer server.Close()
    
    resp, err := http.Get(server.URL + "/users/1")
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
    }
}
```

## 6. Benchmark 测试

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}

// 带设置的 benchmark
func BenchmarkProcess(b *testing.B) {
    data := prepareData()  // 准备数据
    
    b.ResetTimer()  // 重置计时器
    
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// 并行 benchmark
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Add(1, 2)
        }
    })
}

// 运行 benchmark
// go test -bench=. ./...
// go test -bench=BenchmarkAdd -benchmem ./...
```

## 7. 测试覆盖率

```bash
# 生成覆盖率报告
go test -cover ./...

# 生成详细报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # 在浏览器中查看

# 查看函数级别覆盖率
go tool cover -func=coverage.out
```

## 8. 测试最佳实践

```go
// 1. 测试文件命名: xxx_test.go
// 2. 测试函数命名: TestXxx
// 3. 使用 t.Run 组织子测试
// 4. 使用表驱动测试
// 5. 测试边界条件和错误情况

// 6. 使用 testdata 目录存放测试数据
// testdata/input.json
// testdata/expected.json

func TestProcessFile(t *testing.T) {
    input, _ := os.ReadFile("testdata/input.json")
    expected, _ := os.ReadFile("testdata/expected.json")
    
    result := Process(input)
    
    if !bytes.Equal(result, expected) {
        t.Errorf("result mismatch")
    }
}

// 7. 使用 golden files
func TestOutput(t *testing.T) {
    result := Generate()
    
    golden := filepath.Join("testdata", t.Name()+".golden")
    
    if *update {  // -update 标志
        os.WriteFile(golden, result, 0644)
    }
    
    expected, _ := os.ReadFile(golden)
    if !bytes.Equal(result, expected) {
        t.Errorf("result mismatch")
    }
}
```

## 9. 常用测试库

```go
// testify - 断言和 mock
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/mock"

func TestWithTestify(t *testing.T) {
    assert.Equal(t, 5, Add(2, 3))
    assert.NotNil(t, result)
    assert.NoError(t, err)
}

// gomock - 自动生成 mock
// go install github.com/golang/mock/mockgen@latest
// mockgen -source=repository.go -destination=mock_repository.go

// httptest - HTTP 测试（标准库）
// 已在上面介绍
```

## 练习

1. 为你的项目添加表驱动测试
2. 实现一个 Mock 数据库，用于测试 Service 层
3. 编写 HTTP 处理器的测试
4. 添加 benchmark 测试，优化性能瓶颈
