# HTTP 教程（二）：RESTful API

---

## 1. 什么是 API

API = Application Programming Interface（应用程序接口）

**简单说：程序之间"说话"的方式。**

```
前端 App                        后端服务器
   |                               |
   |  ---- 请求用户列表 ---->       |
   |       GET /api/users          |
   |                               |
   |  <---- 返回 JSON ----         |
   |       [{"id":1,"name":"秦阳"}] |
```

---

## 2. 什么是 REST

REST = Representational State Transfer（表述性状态转移）

**别管这个名字，记住核心思想**：

- 把后端数据看成"资源"
- 用 URL 定位资源
- 用 HTTP 方法操作资源

### 2.1 RESTful 的核心

```
资源 = 名词（users, orders, products）
操作 = HTTP 方法（GET, POST, PUT, DELETE）
```

**不是 RESTful**：
```
GET  /getUsers
POST /createUser
POST /deleteUser?id=1
GET  /getUserById?id=1
```

**RESTful**：
```
GET    /users        获取用户列表
GET    /users/1      获取用户1
POST   /users        创建用户
PUT    /users/1      更新用户1
DELETE /users/1      删除用户1
```

---

## 3. RESTful API 设计规范

### 3.1 URL 设计

```
# 资源用名词，复数形式
GET /users              # 用户列表
GET /users/123          # 单个用户
GET /users/123/orders   # 用户的订单列表

# 避免动词
❌ GET /getUsers
❌ POST /createUser
✅ GET /users
✅ POST /users
```

### 3.2 HTTP 方法对应 CRUD

| 方法 | 操作 | URL 示例 | 说明 |
|------|------|----------|------|
| GET | Read | GET /users | 获取列表 |
| GET | Read | GET /users/1 | 获取单个 |
| POST | Create | POST /users | 创建 |
| PUT | Update | PUT /users/1 | 全量更新 |
| PATCH | Update | PATCH /users/1 | 部分更新 |
| DELETE | Delete | DELETE /users/1 | 删除 |

### 3.3 状态码使用

```
# 成功
200 OK           - GET 成功
201 Created      - POST 创建成功
204 No Content   - DELETE 成功

# 客户端错误
400 Bad Request  - 参数错误
401 Unauthorized - 未登录
403 Forbidden    - 没权限
404 Not Found    - 资源不存在
422 Unprocessable Entity - 参数验证失败

# 服务器错误
500 Internal Server Error - 服务器出错
```

### 3.4 请求和响应格式

**请求**：
```http
POST /api/users HTTP/1.1
Host: example.com
Content-Type: application/json

{
    "name": "秦阳",
    "age": 22
}
```

**响应**：
```http
HTTP/1.1 201 Created
Content-Type: application/json

{
    "id": 1,
    "name": "秦阳",
    "age": 22,
    "created_at": "2024-12-25T10:00:00Z"
}
```

---

## 4. 完整的 RESTful API 示例

### 4.1 用户管理 API

```
# 用户列表
GET /api/users
响应: [{"id":1,"name":"秦阳"}, {"id":2,"name":"张三"}]

# 获取单个用户
GET /api/users/1
响应: {"id":1,"name":"秦阳","age":22}

# 创建用户
POST /api/users
请求: {"name":"李四","age":25}
响应: {"id":3,"name":"李四","age":25}

# 更新用户
PUT /api/users/1
请求: {"name":"秦阳","age":23}
响应: {"id":1,"name":"秦阳","age":23}

# 删除用户
DELETE /api/users/1
响应: 204 No Content
```

### 4.2 嵌套资源

```
# 用户的订单列表
GET /api/users/1/orders

# 用户的某个订单
GET /api/users/1/orders/100

# 创建用户的订单
POST /api/users/1/orders
```

### 4.3 查询参数

```
# 分页
GET /api/users?page=1&limit=10

# 排序
GET /api/users?sort=created_at&order=desc

# 过滤
GET /api/users?status=active&age_min=18

# 搜索
GET /api/users?q=秦阳
```

---

## 5. 统一响应格式

### 5.1 成功响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "秦阳"
    }
}
```

### 5.2 列表响应（带分页）

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {"id": 1, "name": "秦阳"},
            {"id": 2, "name": "张三"}
        ],
        "total": 100,
        "page": 1,
        "limit": 10
    }
}
```

### 5.3 错误响应

```json
{
    "code": 40001,
    "message": "参数错误",
    "data": null
}
```

---

## 6. Go 实现 RESTful API

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
)

// 用户结构
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// 模拟数据库
var users = []User{
    {ID: 1, Name: "秦阳", Age: 22},
    {ID: 2, Name: "张三", Age: 25},
}
var nextID = 3

// 统一响应
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func success(w http.ResponseWriter, data interface{}) {
    sendJSON(w, http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func created(w http.ResponseWriter, data interface{}) {
    sendJSON(w, http.StatusCreated, Response{
        Code:    0,
        Message: "created",
        Data:    data,
    })
}

func notFound(w http.ResponseWriter) {
    sendJSON(w, http.StatusNotFound, Response{
        Code:    404,
        Message: "not found",
        Data:    nil,
    })
}

// 处理 /api/users
func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // 获取用户列表
        success(w, users)
        
    case "POST":
        // 创建用户
        var user User
        json.NewDecoder(r.Body).Decode(&user)
        user.ID = nextID
        nextID++
        users = append(users, user)
        created(w, user)
        
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// 处理 /api/users/{id}
func userHandler(w http.ResponseWriter, r *http.Request) {
    // 解析 ID
    path := strings.TrimPrefix(r.URL.Path, "/api/users/")
    id, err := strconv.Atoi(path)
    if err != nil {
        notFound(w)
        return
    }
    
    // 查找用户
    var index = -1
    for i, u := range users {
        if u.ID == id {
            index = i
            break
        }
    }
    
    switch r.Method {
    case "GET":
        if index == -1 {
            notFound(w)
            return
        }
        success(w, users[index])
        
    case "PUT":
        if index == -1 {
            notFound(w)
            return
        }
        var user User
        json.NewDecoder(r.Body).Decode(&user)
        user.ID = id
        users[index] = user
        success(w, user)
        
    case "DELETE":
        if index == -1 {
            notFound(w)
            return
        }
        users = append(users[:index], users[index+1:]...)
        w.WriteHeader(http.StatusNoContent)
        
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func main() {
    http.HandleFunc("/api/users", usersHandler)
    http.HandleFunc("/api/users/", userHandler)
    
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## 7. 测试 API

用 curl 测试：

```bash
# 获取用户列表
curl http://localhost:8080/api/users

# 获取单个用户
curl http://localhost:8080/api/users/1

# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","age":30}'

# 更新用户
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"秦阳","age":23}'

# 删除用户
curl -X DELETE http://localhost:8080/api/users/1
```

---

## 下一篇

[HTTP教程-3-实战进阶](./HTTP教程-3-实战进阶.md)
