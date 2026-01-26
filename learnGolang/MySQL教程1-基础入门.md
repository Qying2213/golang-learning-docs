# MySQL 教程1 - 基础入门

## 1. MySQL 简介

### 1.1 什么是 MySQL？

MySQL 是一个开源的**关系型数据库管理系统（RDBMS）**，由瑞典 MySQL AB 公司开发，现属于 Oracle 公司。它是目前最流行的数据库之一，广泛应用于 Web 开发。

**关系型数据库的特点：**
- 数据以**表格**形式存储
- 表与表之间可以建立**关系**
- 使用 **SQL（结构化查询语言）** 进行操作
- 支持**事务**，保证数据一致性

### 1.2 核心概念

```
数据库服务器 (MySQL Server)
    └── 数据库 (Database) - 相当于一个文件夹
            └── 表 (Table) - 相当于 Excel 表格
                    ├── 列/字段 (Column/Field) - 表格的列
                    └── 行/记录 (Row/Record) - 表格的每一行数据
```

**举个例子：**
假设你要做一个学校管理系统：
- 数据库：`school_db`
- 表：`students`（学生表）、`teachers`（教师表）、`courses`（课程表）
- 字段：学生表有 `id`、`name`、`age`、`class` 等字段

---

## 2. 安装 MySQL

### 2.1 macOS 安装

```bash
# 使用 Homebrew 安装
brew install mysql

# 启动 MySQL 服务
brew services start mysql

# 设置 root 密码（首次安装）
mysql_secure_installation
```

### 2.2 连接 MySQL

```bash
# 连接到本地 MySQL
mysql -u root -p

# 连接到远程 MySQL
mysql -h 主机地址 -P 端口 -u 用户名 -p
```

连接成功后会看到：
```
Welcome to the MySQL monitor.  Commands end with ; or \g.
mysql>
```

**重要提示：** SQL 语句以分号 `;` 结尾！

---

## 3. 数据库操作

### 3.1 查看所有数据库

```sql
SHOW DATABASES;
```

输出示例：
```
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
```

### 3.2 创建数据库

```sql
-- 基本语法
CREATE DATABASE 数据库名;

-- 示例：创建学校数据库
CREATE DATABASE school_db;

-- 推荐写法：指定字符集（支持中文）
CREATE DATABASE school_db 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- 如果数据库已存在不报错
CREATE DATABASE IF NOT EXISTS school_db;
```

### 3.3 选择/使用数据库

```sql
USE school_db;
```

执行后会提示：`Database changed`

### 3.4 查看当前使用的数据库

```sql
SELECT DATABASE();
```

### 3.5 删除数据库

```sql
-- 删除数据库（谨慎操作！）
DROP DATABASE school_db;

-- 安全写法
DROP DATABASE IF EXISTS school_db;
```

⚠️ **警告：** 删除数据库会删除其中所有表和数据，无法恢复！

---

## 4. 数据类型

在创建表之前，必须了解 MySQL 的数据类型。

### 4.1 数值类型

| 类型 | 大小 | 范围 | 用途 |
|------|------|------|------|
| `TINYINT` | 1 字节 | -128 ~ 127 | 很小的整数 |
| `INT` | 4 字节 | -21亿 ~ 21亿 | 常用整数 |
| `BIGINT` | 8 字节 | 很大 | 大整数（如订单号） |
| `FLOAT` | 4 字节 | - | 单精度浮点数 |
| `DOUBLE` | 8 字节 | - | 双精度浮点数 |
| `DECIMAL(M,D)` | - | - | 精确小数（如金额） |

**示例：**
```sql
age INT                    -- 年龄
price DECIMAL(10, 2)       -- 价格，最多10位，小数2位，如 12345678.99
```

### 4.2 字符串类型

| 类型 | 大小 | 用途 |
|------|------|------|
| `CHAR(N)` | 固定 N 字符 | 固定长度字符串（如性别 M/F） |
| `VARCHAR(N)` | 可变，最大 N 字符 | 可变长度字符串（如姓名） |
| `TEXT` | 最大 65535 字符 | 长文本（如文章内容） |
| `LONGTEXT` | 最大 4GB | 超长文本 |

**CHAR vs VARCHAR：**
```sql
-- CHAR(10) 存储 "abc" 会占用 10 个字符空间
-- VARCHAR(10) 存储 "abc" 只占用 3 个字符空间 + 1字节长度标记

name VARCHAR(50)           -- 姓名，最长50字符
gender CHAR(1)             -- 性别，固定1字符
```

### 4.3 日期时间类型

| 类型 | 格式 | 用途 |
|------|------|------|
| `DATE` | YYYY-MM-DD | 日期 |
| `TIME` | HH:MM:SS | 时间 |
| `DATETIME` | YYYY-MM-DD HH:MM:SS | 日期时间 |
| `TIMESTAMP` | YYYY-MM-DD HH:MM:SS | 时间戳（自动更新） |

**示例：**
```sql
birthday DATE              -- 生日 2000-01-15
created_at DATETIME        -- 创建时间 2024-01-15 10:30:00
updated_at TIMESTAMP       -- 更新时间（自动记录）
```

### 4.4 其他常用类型

```sql
BOOLEAN                    -- 布尔值（实际是 TINYINT(1)）
ENUM('值1', '值2', '值3')   -- 枚举，只能选其中一个
JSON                       -- JSON 数据（MySQL 5.7+）
```

---

## 5. 表操作

### 5.1 创建表

**基本语法：**
```sql
CREATE TABLE 表名 (
    列名1 数据类型 [约束],
    列名2 数据类型 [约束],
    ...
);
```

**示例：创建学生表**
```sql
CREATE TABLE students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    age INT,
    gender CHAR(1),
    email VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**解释：**
- `PRIMARY KEY`：主键，唯一标识每条记录
- `AUTO_INCREMENT`：自动递增，插入数据时自动生成
- `NOT NULL`：不能为空
- `DEFAULT`：默认值

### 5.2 常用约束

| 约束 | 说明 | 示例 |
|------|------|------|
| `PRIMARY KEY` | 主键，唯一且不为空 | `id INT PRIMARY KEY` |
| `AUTO_INCREMENT` | 自动递增 | `id INT AUTO_INCREMENT` |
| `NOT NULL` | 不能为空 | `name VARCHAR(50) NOT NULL` |
| `DEFAULT` | 默认值 | `status INT DEFAULT 1` |
| `UNIQUE` | 唯一，不能重复 | `email VARCHAR(100) UNIQUE` |
| `FOREIGN KEY` | 外键，关联其他表 | 后面详细讲 |

### 5.3 查看表结构

```sql
-- 查看所有表
SHOW TABLES;

-- 查看表结构
DESC students;
-- 或者
DESCRIBE students;

-- 查看建表语句
SHOW CREATE TABLE students;
```

### 5.4 修改表结构

```sql
-- 添加列
ALTER TABLE students ADD phone VARCHAR(20);

-- 删除列
ALTER TABLE students DROP COLUMN phone;

-- 修改列类型
ALTER TABLE students MODIFY age TINYINT;

-- 修改列名和类型
ALTER TABLE students CHANGE gender sex CHAR(1);

-- 重命名表
ALTER TABLE students RENAME TO student;
-- 或者
RENAME TABLE students TO student;
```

### 5.5 删除表

```sql
-- 删除表
DROP TABLE students;

-- 安全写法
DROP TABLE IF EXISTS students;

-- 清空表数据（保留表结构）
TRUNCATE TABLE students;
```

---

## 6. 实战：创建完整的数据库

让我们创建一个简单的**博客系统**数据库：

```sql
-- 1. 创建数据库
CREATE DATABASE IF NOT EXISTS blog_db
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

-- 2. 使用数据库
USE blog_db;

-- 3. 创建用户表
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密后）',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    avatar VARCHAR(255) COMMENT '头像URL',
    status TINYINT DEFAULT 1 COMMENT '状态：0禁用 1正常',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '用户表';

-- 4. 创建文章表
CREATE TABLE articles (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '文章ID',
    title VARCHAR(200) NOT NULL COMMENT '标题',
    content TEXT COMMENT '内容',
    author_id INT NOT NULL COMMENT '作者ID',
    view_count INT DEFAULT 0 COMMENT '浏览量',
    status TINYINT DEFAULT 1 COMMENT '状态：0草稿 1发布',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id)
) COMMENT '文章表';

-- 5. 创建评论表
CREATE TABLE comments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    content TEXT NOT NULL COMMENT '评论内容',
    user_id INT NOT NULL COMMENT '评论者ID',
    article_id INT NOT NULL COMMENT '文章ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (article_id) REFERENCES articles(id)
) COMMENT '评论表';
```

**关键点解释：**
1. `COMMENT`：给字段添加注释，方便理解
2. `ON UPDATE CURRENT_TIMESTAMP`：更新记录时自动更新时间
3. `FOREIGN KEY`：外键约束，确保数据完整性

---

## 7. 练习题

### 练习1：创建数据库
创建一个名为 `shop_db` 的数据库，要求支持中文。

<details>
<summary>点击查看答案</summary>

```sql
CREATE DATABASE IF NOT EXISTS shop_db
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;
```
</details>

---

### 练习2：创建商品表
在 `shop_db` 中创建一个商品表 `products`，包含以下字段：
- id：主键，自增
- name：商品名称，最长100字符，不能为空
- price：价格，精确到小数点后2位
- stock：库存数量，默认为0
- description：商品描述，长文本
- created_at：创建时间，默认当前时间

<details>
<summary>点击查看答案</summary>

```sql
USE shop_db;

CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2),
    stock INT DEFAULT 0,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
</details>

---

### 练习3：创建订单表
创建订单表 `orders`，要求：
- id：主键，自增
- order_no：订单号，唯一，不能为空
- user_id：用户ID，不能为空
- total_amount：订单总金额
- status：订单状态，枚举类型（pending, paid, shipped, completed, cancelled）
- created_at 和 updated_at

<details>
<summary>点击查看答案</summary>

```sql
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    total_amount DECIMAL(12, 2),
    status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
</details>

---

### 练习4：修改表结构
1. 给 `products` 表添加一个 `category` 字段（分类，VARCHAR(50)）
2. 给 `products` 表添加一个 `is_on_sale` 字段（是否上架，默认为1）
3. 将 `description` 字段改名为 `detail`

<details>
<summary>点击查看答案</summary>

```sql
-- 添加 category 字段
ALTER TABLE products ADD category VARCHAR(50);

-- 添加 is_on_sale 字段
ALTER TABLE products ADD is_on_sale TINYINT DEFAULT 1;

-- 修改字段名
ALTER TABLE products CHANGE description detail TEXT;
```
</details>

---

### 练习5：查看表信息
写出以下操作的 SQL：
1. 查看 `shop_db` 中所有的表
2. 查看 `products` 表的结构
3. 查看 `products` 表的建表语句

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 查看所有表
SHOW TABLES;

-- 2. 查看表结构
DESC products;

-- 3. 查看建表语句
SHOW CREATE TABLE products;
```
</details>

---

## 8. 本章小结

本章学习了：
1. **MySQL 基本概念**：数据库、表、字段、记录
2. **数据库操作**：CREATE、USE、DROP、SHOW
3. **数据类型**：数值、字符串、日期时间
4. **表操作**：CREATE TABLE、ALTER TABLE、DROP TABLE
5. **约束**：PRIMARY KEY、AUTO_INCREMENT、NOT NULL、DEFAULT、UNIQUE、FOREIGN KEY

**下一章预告：** SQL 增删改查（CRUD）操作
