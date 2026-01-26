# MySQL 教程4 - 索引与性能优化

## 1. 索引简介

### 1.1 什么是索引？

索引就像书的目录，帮助数据库快速定位数据，而不需要扫描整个表。

**没有索引：** 查找一条记录需要从头到尾扫描整个表（全表扫描）
**有索引：** 通过索引快速定位到数据位置

### 1.2 索引的优缺点

**优点：**
- 大幅提高查询速度
- 加速排序和分组操作
- 加速表连接

**缺点：**
- 占用额外存储空间
- 插入、更新、删除时需要维护索引，会变慢
- 创建和维护索引需要时间

**原则：** 读多写少的场景适合建索引

---

## 2. 索引类型

### 2.1 按数据结构分类

| 类型 | 说明 |
|------|------|
| B+Tree 索引 | 最常用，适合范围查询 |
| Hash 索引 | 只支持等值查询，Memory 引擎支持 |
| Full-text 索引 | 全文索引，用于文本搜索 |

### 2.2 按功能分类

| 类型 | 说明 | 关键字 |
|------|------|--------|
| 主键索引 | 唯一且不为空 | `PRIMARY KEY` |
| 唯一索引 | 值唯一，可以为空 | `UNIQUE` |
| 普通索引 | 最基本的索引 | `INDEX` |
| 组合索引 | 多个列组成的索引 | `INDEX(col1, col2)` |
| 全文索引 | 用于文本搜索 | `FULLTEXT` |

---

## 3. 索引操作

### 3.1 创建索引

```sql
-- 创建表时添加索引
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,        -- 主键索引
    username VARCHAR(50) UNIQUE,              -- 唯一索引
    email VARCHAR(100),
    age INT,
    created_at DATETIME,
    INDEX idx_email (email),                  -- 普通索引
    INDEX idx_age_created (age, created_at)   -- 组合索引
);

-- 给已有表添加索引
CREATE INDEX idx_email ON users(email);

-- 添加唯一索引
CREATE UNIQUE INDEX idx_username ON users(username);

-- 使用 ALTER TABLE 添加索引
ALTER TABLE users ADD INDEX idx_age (age);
ALTER TABLE users ADD UNIQUE INDEX idx_email (email);
```

### 3.2 查看索引

```sql
-- 查看表的索引
SHOW INDEX FROM users;

-- 简洁显示
SHOW KEYS FROM users;
```

### 3.3 删除索引

```sql
-- 删除索引
DROP INDEX idx_email ON users;

-- 使用 ALTER TABLE 删除
ALTER TABLE users DROP INDEX idx_email;
```

---

## 4. 索引使用原则

### 4.1 适合建索引的场景

```sql
-- 1. WHERE 条件中经常使用的列
SELECT * FROM users WHERE email = 'test@example.com';
-- 应该给 email 建索引

-- 2. ORDER BY 排序的列
SELECT * FROM users ORDER BY created_at DESC;
-- 应该给 created_at 建索引

-- 3. JOIN 连接的列
SELECT * FROM orders o JOIN users u ON o.user_id = u.id;
-- user_id 和 id 应该有索引

-- 4. GROUP BY 分组的列
SELECT department, COUNT(*) FROM employees GROUP BY department;
-- department 应该有索引
```

### 4.2 不适合建索引的场景

1. **数据量小的表**：全表扫描可能更快
2. **频繁更新的列**：维护索引开销大
3. **区分度低的列**：如性别（只有男/女）
4. **很少用于查询的列**

### 4.3 组合索引的最左前缀原则

```sql
-- 创建组合索引
CREATE INDEX idx_abc ON users(a, b, c);

-- 以下查询可以使用索引：
WHERE a = 1                    -- ✅ 使用 a
WHERE a = 1 AND b = 2          -- ✅ 使用 a, b
WHERE a = 1 AND b = 2 AND c = 3 -- ✅ 使用 a, b, c
WHERE a = 1 AND c = 3          -- ✅ 只使用 a（c 不连续）

-- 以下查询不能使用索引：
WHERE b = 2                    -- ❌ 没有 a
WHERE b = 2 AND c = 3          -- ❌ 没有 a
WHERE c = 3                    -- ❌ 没有 a
```

**记忆口诀：** 最左前缀，不能跳过！

---

## 5. EXPLAIN 分析查询

### 5.1 基本用法

```sql
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';
```

### 5.2 重要字段解释

| 字段 | 说明 |
|------|------|
| `id` | 查询序号 |
| `select_type` | 查询类型（SIMPLE/PRIMARY/SUBQUERY等） |
| `table` | 查询的表 |
| `type` | 访问类型（重要！） |
| `possible_keys` | 可能使用的索引 |
| `key` | 实际使用的索引 |
| `key_len` | 索引使用的字节数 |
| `rows` | 预估扫描行数 |
| `Extra` | 额外信息 |

### 5.3 type 字段详解（性能从好到差）

| type | 说明 |
|------|------|
| `system` | 表只有一行 |
| `const` | 通过主键或唯一索引查询，最多一条 |
| `eq_ref` | 连接查询时使用主键或唯一索引 |
| `ref` | 使用普通索引 |
| `range` | 索引范围扫描 |
| `index` | 全索引扫描 |
| `ALL` | 全表扫描（最差！） |

**目标：** 至少达到 `range` 级别，避免 `ALL`

### 5.4 示例分析

```sql
-- 准备测试表
CREATE TABLE test_explain (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    age INT,
    email VARCHAR(100),
    INDEX idx_name (name),
    INDEX idx_age (age)
);

-- 插入测试数据
INSERT INTO test_explain (name, age, email) VALUES 
    ('张三', 25, 'zhangsan@test.com'),
    ('李四', 30, 'lisi@test.com'),
    ('王五', 28, 'wangwu@test.com');

-- 分析查询
EXPLAIN SELECT * FROM test_explain WHERE id = 1;
-- type: const（主键查询，最优）

EXPLAIN SELECT * FROM test_explain WHERE name = '张三';
-- type: ref（使用普通索引）

EXPLAIN SELECT * FROM test_explain WHERE age > 25;
-- type: range（范围查询）

EXPLAIN SELECT * FROM test_explain WHERE email = 'test@test.com';
-- type: ALL（全表扫描，email 没有索引）
```

---

## 6. SQL 优化技巧

### 6.1 避免索引失效

```sql
-- ❌ 在索引列上使用函数
SELECT * FROM users WHERE YEAR(created_at) = 2024;
-- ✅ 改为范围查询
SELECT * FROM users WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';

-- ❌ 在索引列上进行运算
SELECT * FROM users WHERE age + 1 = 26;
-- ✅ 改为
SELECT * FROM users WHERE age = 25;

-- ❌ 使用 != 或 <>
SELECT * FROM users WHERE status != 1;
-- ✅ 如果可能，改为 IN
SELECT * FROM users WHERE status IN (0, 2, 3);

-- ❌ LIKE 以 % 开头
SELECT * FROM users WHERE name LIKE '%三';
-- ✅ LIKE 不以 % 开头
SELECT * FROM users WHERE name LIKE '张%';

-- ❌ OR 连接非索引列
SELECT * FROM users WHERE name = '张三' OR age = 25;
-- ✅ 使用 UNION
SELECT * FROM users WHERE name = '张三'
UNION
SELECT * FROM users WHERE age = 25;

-- ❌ 隐式类型转换
SELECT * FROM users WHERE phone = 13800138000;  -- phone 是 VARCHAR
-- ✅ 使用正确的类型
SELECT * FROM users WHERE phone = '13800138000';
```

### 6.2 SELECT 优化

```sql
-- ❌ 使用 SELECT *
SELECT * FROM users;

-- ✅ 只查询需要的列
SELECT id, name, email FROM users;

-- ❌ 查询大量数据不加 LIMIT
SELECT * FROM logs WHERE created_at > '2024-01-01';

-- ✅ 加上 LIMIT 限制
SELECT * FROM logs WHERE created_at > '2024-01-01' LIMIT 1000;
```

### 6.3 JOIN 优化

```sql
-- 1. 确保连接字段有索引
-- 2. 小表驱动大表
-- 3. 避免过多的表连接（一般不超过3-4个）

-- ❌ 大表在前
SELECT * FROM big_table b
JOIN small_table s ON b.id = s.big_id;

-- ✅ 小表在前
SELECT * FROM small_table s
JOIN big_table b ON s.big_id = b.id;
```

### 6.4 分页优化

```sql
-- ❌ 深分页性能差
SELECT * FROM users ORDER BY id LIMIT 1000000, 10;

-- ✅ 使用游标分页（记住上次的 ID）
SELECT * FROM users WHERE id > 1000000 ORDER BY id LIMIT 10;

-- ✅ 使用子查询优化
SELECT * FROM users
WHERE id >= (SELECT id FROM users ORDER BY id LIMIT 1000000, 1)
LIMIT 10;
```

### 6.5 COUNT 优化

```sql
-- ❌ COUNT(*)（需要扫描所有行）
SELECT COUNT(*) FROM users;

-- ✅ 如果只需要判断是否存在
SELECT 1 FROM users LIMIT 1;

-- ✅ 使用近似值（如果精确度要求不高）
SHOW TABLE STATUS LIKE 'users';  -- 查看 Rows 字段

-- ✅ 维护计数表（高并发场景）
```

---

## 7. 慢查询日志

### 7.1 开启慢查询日志

```sql
-- 查看是否开启
SHOW VARIABLES LIKE 'slow_query_log';

-- 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';

-- 设置慢查询阈值（秒）
SET GLOBAL long_query_time = 2;

-- 查看日志文件位置
SHOW VARIABLES LIKE 'slow_query_log_file';
```

### 7.2 分析慢查询日志

```bash
# 使用 mysqldumpslow 工具分析
mysqldumpslow -s t -t 10 /var/log/mysql/slow.log

# -s t: 按查询时间排序
# -t 10: 显示前10条
```

---

## 8. 实战练习

### 练习1：创建索引
给以下表创建合适的索引：
```sql
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(50),
    user_id INT,
    status TINYINT,
    total_amount DECIMAL(10,2),
    created_at DATETIME
);
```
常见查询：
1. 根据订单号查询
2. 查询某用户的订单
3. 查询某时间段的订单
4. 查询某用户某状态的订单

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 订单号唯一索引
CREATE UNIQUE INDEX idx_order_no ON orders(order_no);

-- 2. 用户ID索引
CREATE INDEX idx_user_id ON orders(user_id);

-- 3. 创建时间索引
CREATE INDEX idx_created_at ON orders(created_at);

-- 4. 用户ID + 状态 组合索引
CREATE INDEX idx_user_status ON orders(user_id, status);
```
</details>

---

### 练习2：分析查询
使用 EXPLAIN 分析以下查询，并说明如何优化：

```sql
SELECT * FROM orders 
WHERE YEAR(created_at) = 2024 
AND status = 1 
ORDER BY created_at DESC;
```

<details>
<summary>点击查看答案</summary>

**问题分析：**
1. `YEAR(created_at)` 在索引列上使用函数，导致索引失效
2. `SELECT *` 查询所有列，可能不必要

**优化后：**
```sql
SELECT id, order_no, user_id, status, total_amount, created_at
FROM orders 
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01'
AND status = 1 
ORDER BY created_at DESC;

-- 并创建组合索引
CREATE INDEX idx_status_created ON orders(status, created_at);
```
</details>

---

### 练习3：优化分页
优化以下深分页查询：

```sql
SELECT * FROM users ORDER BY id LIMIT 500000, 20;
```

<details>
<summary>点击查看答案</summary>

```sql
-- 方法1：使用游标分页（需要知道上一页最后的ID）
SELECT * FROM users WHERE id > 500000 ORDER BY id LIMIT 20;

-- 方法2：子查询优化
SELECT * FROM users
WHERE id >= (SELECT id FROM users ORDER BY id LIMIT 500000, 1)
ORDER BY id
LIMIT 20;

-- 方法3：延迟关联
SELECT u.* FROM users u
INNER JOIN (SELECT id FROM users ORDER BY id LIMIT 500000, 20) t
ON u.id = t.id;
```
</details>

---

### 练习4：索引选择
以下查询应该创建什么索引？

```sql
-- 查询1
SELECT * FROM products WHERE category_id = 5 AND price > 100 ORDER BY created_at DESC;

-- 查询2
SELECT * FROM products WHERE name LIKE '手机%' AND status = 1;

-- 查询3
SELECT user_id, COUNT(*) FROM orders WHERE created_at > '2024-01-01' GROUP BY user_id;
```

<details>
<summary>点击查看答案</summary>

```sql
-- 查询1：组合索引（注意顺序）
CREATE INDEX idx_cat_price_created ON products(category_id, price, created_at);
-- 或者分开建
CREATE INDEX idx_category ON products(category_id);
CREATE INDEX idx_created ON products(created_at);

-- 查询2：组合索引
CREATE INDEX idx_status_name ON products(status, name);
-- 注意：name LIKE '手机%' 可以使用索引，但 LIKE '%手机' 不行

-- 查询3：组合索引
CREATE INDEX idx_created_user ON orders(created_at, user_id);
```
</details>

---

### 练习5：综合优化
分析并优化以下查询：

```sql
SELECT u.name, COUNT(o.id) as order_count, SUM(o.total_amount) as total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.status != 0
AND o.created_at BETWEEN '2024-01-01' AND '2024-12-31'
GROUP BY u.id
HAVING COUNT(o.id) > 5
ORDER BY total DESC
LIMIT 10;
```

<details>
<summary>点击查看答案</summary>

**优化建议：**

1. **索引优化：**
```sql
-- orders 表索引
CREATE INDEX idx_user_status_created ON orders(user_id, status, created_at);
```

2. **SQL 优化：**
```sql
SELECT u.name, t.order_count, t.total
FROM users u
INNER JOIN (
    SELECT user_id, COUNT(id) as order_count, SUM(total_amount) as total
    FROM orders
    WHERE status IN (1, 2, 3)  -- 改为 IN，避免 !=
    AND created_at >= '2024-01-01' AND created_at < '2025-01-01'
    GROUP BY user_id
    HAVING COUNT(id) > 5
    ORDER BY total DESC
    LIMIT 10
) t ON u.id = t.user_id
ORDER BY t.total DESC;
```

**优化点：**
- 将 `!= 0` 改为 `IN (1,2,3)`
- 先在子查询中聚合和过滤，减少数据量
- 使用 INNER JOIN 替代 LEFT JOIN（因为 HAVING 已经过滤了无订单用户）
</details>

---

## 9. 本章小结

**索引类型：**
- 主键索引、唯一索引、普通索引、组合索引

**索引原则：**
- 最左前缀原则
- 避免在索引列上使用函数、运算
- 避免 LIKE '%xxx'
- 注意隐式类型转换

**优化工具：**
- EXPLAIN 分析查询
- 慢查询日志

**优化技巧：**
- 避免 SELECT *
- 小表驱动大表
- 深分页优化
- 合理使用组合索引

**下一章预告：** Go 语言操作 MySQL（database/sql 和 GORM）
