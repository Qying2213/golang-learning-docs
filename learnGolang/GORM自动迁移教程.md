# GORM 自动迁移教程（PostgreSQL）

## 什么是自动迁移？

**自动迁移 = 根据 Go 结构体自动创建/修改数据库表**

你改了 struct，GORM 帮你改数据库表，不用手写 SQL。

```
Go struct 改了字段
       ↓
AutoMigrate()
       ↓
数据库表自动更新
```

---

## 第一步：连接 PostgreSQL

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // PostgreSQL 连接字符串
    dsn := "host=localhost user=postgres password=123456 dbname=testdb port=5432 sslmode=disable"
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("连接失败:", err)
    }
    
    fmt.Println("连接成功")
}
```

**DSN 参数说明：**
| 参数 | 说明 |
|------|------|
| host | 数据库地址 |
| user | 用户名 |
| password | 密码 |
| dbname | 数据库名 |
| port | 端口（默认 5432） |
| sslmode | SSL 模式（本地开发用 disable） |

---

## 第二步：定义模型

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Username  string         `gorm:"column:username;type:varchar(50);not null;unique"`
    Password  string         `gorm:"column:password;type:varchar(100);not null"`
    Email     string         `gorm:"column:email;type:varchar(100)"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## 第三步：执行自动迁移

```go
func main() {
    // 连接数据库
    dsn := "host=localhost user=postgres password=123456 dbname=testdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 自动迁移（创建表）
    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal("迁移失败:", err)
    }
    
    fmt.Println("迁移成功")
}
```

**执行后，PostgreSQL 会自动创建 users 表：**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

---

## 第四步：添加新字段

### 4.1 在 struct 里加字段

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Username  string         `gorm:"column:username;type:varchar(50);not null;unique"`
    Password  string         `gorm:"column:password;type:varchar(100);not null"`
    Email     string         `gorm:"column:email;type:varchar(100)"`
    
    // ========== 新增字段 ==========
    Phone     string         `gorm:"column:phone;type:varchar(20)"`           // 新增：手机号
    Age       int            `gorm:"column:age;type:int;default:0"`           // 新增：年龄
    IsActive  bool           `gorm:"column:is_active;type:bool;default:true"` // 新增：是否激活
    
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 4.2 重新运行 AutoMigrate

```go
err = db.AutoMigrate(&User{})
```

**GORM 会自动执行：**
```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN age INT DEFAULT 0;
ALTER TABLE users ADD COLUMN is_active BOOL DEFAULT true;
```

---

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// 用户模型
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"column:username;type:varchar(50);not null;unique"`
    Password  string    `gorm:"column:password;type:varchar(100);not null"`
    Email     string    `gorm:"column:email;type:varchar(100)"`
    Phone     string    `gorm:"column:phone;type:varchar(20)"`
    Age       int       `gorm:"column:age;type:int;default:0"`
    IsActive  bool      `gorm:"column:is_active;type:bool;default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func main() {
    // 1. 连接数据库
    dsn := "host=localhost user=postgres password=123456 dbname=testdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("连接失败:", err)
    }
    fmt.Println("✅ 数据库连接成功")

    // 2. 自动迁移
    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal("迁移失败:", err)
    }
    fmt.Println("✅ 自动迁移成功")

    // 3. 测试插入数据
    user := User{
        Username: "qinyang",
        Password: "123456",
        Email:    "qinyang@example.com",
        Phone:    "13800138000",
        Age:      22,
        IsActive: true,
    }
    
    result := db.Create(&user)
    if result.Error != nil {
        log.Fatal("插入失败:", result.Error)
    }
    fmt.Printf("✅ 插入成功，ID: %d\n", user.ID)
}
```

---

## 常用 GORM 标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `primaryKey` | 主键 | `gorm:"primaryKey"` |
| `column` | 指定列名 | `gorm:"column:user_name"` |
| `type` | 指定类型 | `gorm:"type:varchar(100)"` |
| `not null` | 非空 | `gorm:"not null"` |
| `unique` | 唯一 | `gorm:"unique"` |
| `default` | 默认值 | `gorm:"default:0"` |
| `index` | 索引 | `gorm:"index"` |
| `uniqueIndex` | 唯一索引 | `gorm:"uniqueIndex"` |

**组合使用：**
```go
Username string `gorm:"column:username;type:varchar(50);not null;unique;index"`
```

---

## PostgreSQL 常用类型对应

| Go 类型 | PostgreSQL 类型 | GORM 写法 |
|---------|----------------|-----------|
| string | VARCHAR | `type:varchar(100)` |
| string | TEXT | `type:text` |
| int | INTEGER | `type:int` |
| int64 | BIGINT | `type:bigint` |
| float64 | DOUBLE PRECISION | `type:double precision` |
| bool | BOOLEAN | `type:bool` |
| time.Time | TIMESTAMP | `type:timestamp` |
| []byte | BYTEA | `type:bytea` |

---

## 自动迁移的限制

**AutoMigrate 能做的：**
- ✅ 创建表
- ✅ 添加新字段
- ✅ 创建索引

**AutoMigrate 不能做的：**
- ❌ 删除字段（安全考虑）
- ❌ 修改字段类型
- ❌ 删除索引

**如果需要删除字段或修改类型，要手动执行 SQL：**
```go
// 删除字段
db.Exec("ALTER TABLE users DROP COLUMN phone")

// 修改字段类型
db.Exec("ALTER TABLE users ALTER COLUMN age TYPE bigint")
```

---

## 多个模型一起迁移

```go
// 一次迁移多个表
err = db.AutoMigrate(
    &User{},
    &Product{},
    &Order{},
)
```

---

## 查看生成的 SQL（调试用）

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),  // 打印 SQL
})
```

运行后会在控制台看到：
```
[INFO] CREATE TABLE "users" ...
[INFO] ALTER TABLE "users" ADD COLUMN "phone" ...
```

---

## 练习

### 练习1：给 User 添加一个 Avatar 字段

```go
type User struct {
    // ... 原有字段
    
    // TODO: 添加 Avatar 字段
    // 要求：列名 avatar，类型 varchar(200)，可为空
}
```

<details>
<summary>点击查看答案</summary>

```go
Avatar string `gorm:"column:avatar;type:varchar(200)"`
```

</details>

### 练习2：添加一个带默认值的 Role 字段

```go
type User struct {
    // ... 原有字段
    
    // TODO: 添加 Role 字段
    // 要求：列名 role，类型 varchar(20)，默认值 "user"
}
```

<details>
<summary>点击查看答案</summary>

```go
Role string `gorm:"column:role;type:varchar(20);default:'user'"`
```

</details>

---

## 总结

| 操作 | 方法 |
|------|------|
| 连接 PostgreSQL | `gorm.Open(postgres.Open(dsn), &gorm.Config{})` |
| 自动迁移 | `db.AutoMigrate(&User{})` |
| 添加字段 | 在 struct 加字段，重新 AutoMigrate |
| 删除字段 | 手动执行 SQL |
| 修改类型 | 手动执行 SQL |

**流程：**
1. 改 struct
2. 运行程序（会自动 AutoMigrate）
3. 数据库表自动更新
