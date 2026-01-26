# MySQL 教程2 - CRUD 增删改查

## 1. CRUD 概述

CRUD 是数据库操作的四个基本操作：
- **C**reate（创建）→ `INSERT`
- **R**ead（读取）→ `SELECT`
- **U**pdate（更新）→ `UPDATE`
- **D**elete（删除）→ `DELETE`

这是后端开发最常用的 SQL 操作，必须熟练掌握！

---

## 2. 准备工作

先创建测试数据库和表：

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS test_db
DEFAULT CHARACTER SET utf8mb4;

USE test_db;

-- 创建员工表
CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    age INT,
    department VARCHAR(50),
    salary DECIMAL(10, 2),
    hire_date DATE,
    email VARCHAR(100) UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建部门表
CREATE TABLE departments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    manager_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 3. INSERT - 插入数据

### 3.1 插入单条数据

**语法：**
```sql
INSERT INTO 表名 (列1, 列2, ...) VALUES (值1, 值2, ...);
```

**示例：**
```sql
-- 指定列插入
INSERT INTO employees (name, age, department, salary, hire_date, email)
VALUES ('张三', 25, '技术部', 15000.00, '2023-01-15', 'zhangsan@example.com');

-- 插入所有列（不推荐，列顺序可能变化）
INSERT INTO employees 
VALUES (NULL, '李四', 28, '市场部', 12000.00, '2023-03-20', 'lisi@example.com', NOW());
```

### 3.2 插入多条数据

```sql
INSERT INTO employees (name, age, department, salary, hire_date, email)
VALUES 
    ('王五', 30, '技术部', 18000.00, '2022-06-01', 'wangwu@example.com'),
    ('赵六', 26, '人事部', 10000.00, '2023-05-10', 'zhaoliu@example.com'),
    ('钱七', 35, '技术部', 25000.00, '2020-01-01', 'qianqi@example.com'),
    ('孙八', 24, '市场部', 9000.00, '2023-08-15', 'sunba@example.com'),
    ('周九', 32, '财务部', 14000.00, '2021-11-20', 'zhoujiu@example.com');
```

### 3.3 插入部门数据

```sql
INSERT INTO departments (name, manager_id)
VALUES 
    ('技术部', 1),
    ('市场部', 2),
    ('人事部', 4),
    ('财务部', 7);
```

---

## 4. SELECT - 查询数据（重点！）

SELECT 是最常用、最复杂的 SQL 语句，需要重点掌握。

### 4.1 基本查询

```sql
-- 查询所有列
SELECT * FROM employees;

-- 查询指定列
SELECT name, age, salary FROM employees;

-- 使用别名
SELECT name AS 姓名, age AS 年龄, salary AS 工资 FROM employees;
-- 或者省略 AS
SELECT name 姓名, age 年龄, salary 工资 FROM employees;
```

### 4.2 去重查询

```sql
-- 查询所有部门（去重）
SELECT DISTINCT department FROM employees;
```

### 4.3 条件查询 WHERE

**比较运算符：**
| 运算符 | 说明 |
|--------|------|
| `=` | 等于 |
| `<>` 或 `!=` | 不等于 |
| `>` `<` `>=` `<=` | 大于、小于等 |
| `BETWEEN...AND` | 在范围内 |
| `IN (...)` | 在列表中 |
| `LIKE` | 模糊匹配 |
| `IS NULL` | 为空 |
| `IS NOT NULL` | 不为空 |

**示例：**
```sql
-- 查询年龄大于25的员工
SELECT * FROM employees WHERE age > 25;

-- 查询技术部的员工
SELECT * FROM employees WHERE department = '技术部';

-- 查询工资在10000到20000之间
SELECT * FROM employees WHERE salary BETWEEN 10000 AND 20000;

-- 查询技术部或市场部的员工
SELECT * FROM employees WHERE department IN ('技术部', '市场部');

-- 查询姓张的员工（模糊匹配）
SELECT * FROM employees WHERE name LIKE '张%';

-- 查询名字包含"三"的员工
SELECT * FROM employees WHERE name LIKE '%三%';

-- 查询名字第二个字是"四"的员工
SELECT * FROM employees WHERE name LIKE '_四%';
```

**LIKE 通配符：**
- `%`：匹配任意多个字符（包括0个）
- `_`：匹配单个字符

### 4.4 逻辑运算符

```sql
-- AND：同时满足
SELECT * FROM employees 
WHERE department = '技术部' AND salary > 15000;

-- OR：满足其一
SELECT * FROM employees 
WHERE department = '技术部' OR department = '市场部';

-- NOT：取反
SELECT * FROM employees 
WHERE NOT department = '技术部';

-- 组合使用（注意括号）
SELECT * FROM employees 
WHERE (department = '技术部' OR department = '市场部') AND age > 25;
```

### 4.5 排序 ORDER BY

```sql
-- 按工资升序（默认）
SELECT * FROM employees ORDER BY salary;
SELECT * FROM employees ORDER BY salary ASC;

-- 按工资降序
SELECT * FROM employees ORDER BY salary DESC;

-- 多字段排序：先按部门升序，再按工资降序
SELECT * FROM employees 
ORDER BY department ASC, salary DESC;
```

### 4.6 限制结果 LIMIT

```sql
-- 查询前3条
SELECT * FROM employees LIMIT 3;

-- 跳过2条，查询3条（分页常用）
SELECT * FROM employees LIMIT 2, 3;
-- 或者
SELECT * FROM employees LIMIT 3 OFFSET 2;

-- 分页公式：LIMIT (页码-1)*每页数量, 每页数量
-- 第1页，每页5条
SELECT * FROM employees LIMIT 0, 5;
-- 第2页，每页5条
SELECT * FROM employees LIMIT 5, 5;
-- 第3页，每页5条
SELECT * FROM employees LIMIT 10, 5;
```

### 4.7 聚合函数

| 函数 | 说明 |
|------|------|
| `COUNT()` | 计数 |
| `SUM()` | 求和 |
| `AVG()` | 平均值 |
| `MAX()` | 最大值 |
| `MIN()` | 最小值 |

```sql
-- 统计员工总数
SELECT COUNT(*) AS 员工总数 FROM employees;

-- 统计技术部员工数
SELECT COUNT(*) FROM employees WHERE department = '技术部';

-- 计算平均工资
SELECT AVG(salary) AS 平均工资 FROM employees;

-- 计算工资总和
SELECT SUM(salary) AS 工资总和 FROM employees;

-- 查询最高和最低工资
SELECT MAX(salary) AS 最高工资, MIN(salary) AS 最低工资 FROM employees;

-- 组合使用
SELECT 
    COUNT(*) AS 员工数,
    AVG(salary) AS 平均工资,
    MAX(salary) AS 最高工资,
    MIN(salary) AS 最低工资,
    SUM(salary) AS 工资总和
FROM employees;
```

### 4.8 分组查询 GROUP BY

```sql
-- 按部门分组，统计每个部门的员工数
SELECT department, COUNT(*) AS 员工数 
FROM employees 
GROUP BY department;

-- 按部门分组，统计每个部门的平均工资
SELECT department, AVG(salary) AS 平均工资 
FROM employees 
GROUP BY department;

-- 按部门分组，统计多个指标
SELECT 
    department AS 部门,
    COUNT(*) AS 员工数,
    AVG(salary) AS 平均工资,
    MAX(salary) AS 最高工资
FROM employees 
GROUP BY department;
```

### 4.9 分组过滤 HAVING

`HAVING` 用于过滤分组后的结果（`WHERE` 不能用于聚合函数）

```sql
-- 查询员工数大于1的部门
SELECT department, COUNT(*) AS 员工数 
FROM employees 
GROUP BY department
HAVING COUNT(*) > 1;

-- 查询平均工资大于12000的部门
SELECT department, AVG(salary) AS 平均工资 
FROM employees 
GROUP BY department
HAVING AVG(salary) > 12000;

-- WHERE 和 HAVING 的区别
-- WHERE：分组前过滤，不能用聚合函数
-- HAVING：分组后过滤，可以用聚合函数

-- 示例：查询技术部和市场部中，平均工资大于10000的部门
SELECT department, AVG(salary) AS avg_salary
FROM employees
WHERE department IN ('技术部', '市场部')  -- 先过滤部门
GROUP BY department
HAVING AVG(salary) > 10000;               -- 再过滤平均工资
```

### 4.10 SELECT 执行顺序

```sql
SELECT department, AVG(salary)    -- 5. 选择列
FROM employees                     -- 1. 从表中
WHERE age > 20                     -- 2. 先过滤行
GROUP BY department                -- 3. 分组
HAVING AVG(salary) > 10000         -- 4. 过滤分组
ORDER BY AVG(salary) DESC          -- 6. 排序
LIMIT 3;                           -- 7. 限制结果
```

**执行顺序：FROM → WHERE → GROUP BY → HAVING → SELECT → ORDER BY → LIMIT**

---

## 5. UPDATE - 更新数据

### 5.1 基本语法

```sql
UPDATE 表名 SET 列1=值1, 列2=值2, ... WHERE 条件;
```

⚠️ **警告：** 一定要加 WHERE 条件，否则会更新所有记录！

### 5.2 示例

```sql
-- 更新单个字段
UPDATE employees SET salary = 20000 WHERE name = '张三';

-- 更新多个字段
UPDATE employees 
SET salary = 22000, department = '研发部' 
WHERE name = '张三';

-- 工资涨10%
UPDATE employees SET salary = salary * 1.1 WHERE department = '技术部';

-- 更新所有人的年龄+1（危险操作！）
UPDATE employees SET age = age + 1;
```

### 5.3 安全模式

MySQL 默认开启安全模式，不允许没有 WHERE 或 LIMIT 的 UPDATE/DELETE：

```sql
-- 临时关闭安全模式
SET SQL_SAFE_UPDATES = 0;

-- 执行更新...

-- 重新开启
SET SQL_SAFE_UPDATES = 1;
```

---

## 6. DELETE - 删除数据

### 6.1 基本语法

```sql
DELETE FROM 表名 WHERE 条件;
```

⚠️ **警告：** 一定要加 WHERE 条件，否则会删除所有记录！

### 6.2 示例

```sql
-- 删除指定员工
DELETE FROM employees WHERE name = '孙八';

-- 删除工资低于10000的员工
DELETE FROM employees WHERE salary < 10000;

-- 删除所有数据（危险！）
DELETE FROM employees;

-- 清空表（更快，重置自增ID）
TRUNCATE TABLE employees;
```

### 6.3 DELETE vs TRUNCATE

| 特性 | DELETE | TRUNCATE |
|------|--------|----------|
| 可加 WHERE | ✅ | ❌ |
| 触发触发器 | ✅ | ❌ |
| 重置自增ID | ❌ | ✅ |
| 速度 | 慢 | 快 |
| 可回滚 | ✅ | ❌ |

---

## 7. 实战练习

### 练习1：插入数据
向 employees 表插入以下员工：
- 姓名：刘十，年龄：29，部门：技术部，工资：16000，入职日期：2022-09-01，邮箱：liushi@example.com

<details>
<summary>点击查看答案</summary>

```sql
INSERT INTO employees (name, age, department, salary, hire_date, email)
VALUES ('刘十', 29, '技术部', 16000.00, '2022-09-01', 'liushi@example.com');
```
</details>

---

### 练习2：基础查询
1. 查询所有员工的姓名和工资
2. 查询所有不重复的部门名称
3. 查询工资最高的3名员工

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 查询姓名和工资
SELECT name, salary FROM employees;

-- 2. 查询不重复的部门
SELECT DISTINCT department FROM employees;

-- 3. 工资最高的3名员工
SELECT * FROM employees ORDER BY salary DESC LIMIT 3;
```
</details>

---

### 练习3：条件查询
1. 查询技术部工资大于15000的员工
2. 查询2023年入职的员工
3. 查询姓名以"张"或"王"开头的员工

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 技术部工资大于15000
SELECT * FROM employees 
WHERE department = '技术部' AND salary > 15000;

-- 2. 2023年入职
SELECT * FROM employees 
WHERE hire_date BETWEEN '2023-01-01' AND '2023-12-31';
-- 或者
SELECT * FROM employees 
WHERE YEAR(hire_date) = 2023;

-- 3. 姓张或姓王
SELECT * FROM employees 
WHERE name LIKE '张%' OR name LIKE '王%';
```
</details>

---

### 练习4：聚合与分组
1. 统计每个部门的员工数量
2. 计算每个部门的平均工资，并按平均工资降序排列
3. 查询员工数量大于等于2的部门

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 每个部门的员工数量
SELECT department, COUNT(*) AS 员工数 
FROM employees 
GROUP BY department;

-- 2. 每个部门的平均工资，降序
SELECT department, AVG(salary) AS 平均工资 
FROM employees 
GROUP BY department
ORDER BY 平均工资 DESC;

-- 3. 员工数量>=2的部门
SELECT department, COUNT(*) AS 员工数 
FROM employees 
GROUP BY department
HAVING COUNT(*) >= 2;
```
</details>

---

### 练习5：更新数据
1. 将技术部所有员工的工资增加2000
2. 将"张三"的部门改为"研发部"，邮箱改为"zhangsan@newmail.com"

<details>
<summary>点击查看答案</summary>

```sql
-- 1. 技术部工资+2000
UPDATE employees 
SET salary = salary + 2000 
WHERE department = '技术部';

-- 2. 更新张三信息
UPDATE employees 
SET department = '研发部', email = 'zhangsan@newmail.com' 
WHERE name = '张三';
```
</details>

---

### 练习6：综合查询
写一条 SQL 查询：
- 查询每个部门工资最高的员工姓名和工资
- 只显示平均工资大于12000的部门
- 按最高工资降序排列

<details>
<summary>点击查看答案</summary>

```sql
-- 方法1：使用子查询
SELECT e.department, e.name, e.salary
FROM employees e
WHERE e.salary = (
    SELECT MAX(salary) 
    FROM employees 
    WHERE department = e.department
)
AND e.department IN (
    SELECT department 
    FROM employees 
    GROUP BY department 
    HAVING AVG(salary) > 12000
)
ORDER BY e.salary DESC;

-- 方法2：使用窗口函数（MySQL 8.0+）
SELECT department, name, salary
FROM (
    SELECT 
        department, 
        name, 
        salary,
        ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as rn,
        AVG(salary) OVER (PARTITION BY department) as avg_salary
    FROM employees
) t
WHERE rn = 1 AND avg_salary > 12000
ORDER BY salary DESC;
```
</details>

---

### 练习7：分页查询
实现分页功能：每页显示3条记录
1. 查询第1页
2. 查询第2页
3. 查询第3页

<details>
<summary>点击查看答案</summary>

```sql
-- 第1页（跳过0条，取3条）
SELECT * FROM employees LIMIT 0, 3;

-- 第2页（跳过3条，取3条）
SELECT * FROM employees LIMIT 3, 3;

-- 第3页（跳过6条，取3条）
SELECT * FROM employees LIMIT 6, 3;

-- 通用公式：LIMIT (page - 1) * pageSize, pageSize
```
</details>

---

## 8. 本章小结

本章学习了 CRUD 四大操作：

| 操作 | SQL | 说明 |
|------|-----|------|
| 创建 | `INSERT INTO ... VALUES ...` | 插入数据 |
| 读取 | `SELECT ... FROM ... WHERE ...` | 查询数据 |
| 更新 | `UPDATE ... SET ... WHERE ...` | 更新数据 |
| 删除 | `DELETE FROM ... WHERE ...` | 删除数据 |

**SELECT 关键字总结：**
- `DISTINCT`：去重
- `WHERE`：条件过滤
- `ORDER BY`：排序
- `LIMIT`：限制结果
- `GROUP BY`：分组
- `HAVING`：分组后过滤
- 聚合函数：`COUNT`、`SUM`、`AVG`、`MAX`、`MIN`

**下一章预告：** 多表查询（JOIN）
