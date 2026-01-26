# MySQL 教程3 - 多表查询（JOIN）

## 1. 为什么需要多表查询？

在实际项目中，数据通常分散在多个表中。比如：
- 用户信息在 `users` 表
- 订单信息在 `orders` 表
- 商品信息在 `products` 表

当我们需要"查询用户的订单详情"时，就需要**关联多个表**进行查询。

---

## 2. 准备测试数据

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS join_test DEFAULT CHARACTER SET utf8mb4;
USE join_test;

-- 学生表
CREATE TABLE students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    class_id INT,
    age INT
);

-- 班级表
CREATE TABLE classes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    teacher VARCHAR(50)
);

-- 课程表
CREATE TABLE courses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    credit INT
);

-- 成绩表（学生-课程 多对多关系）
CREATE TABLE scores (
    id INT PRIMARY KEY AUTO_INCREMENT,
    student_id INT,
    course_id INT,
    score DECIMAL(5, 2)
);

-- 插入班级数据
INSERT INTO classes (name, teacher) VALUES 
    ('一班', '张老师'),
    ('二班', '李老师'),
    ('三班', '王老师');

-- 插入学生数据
INSERT INTO students (name, class_id, age) VALUES 
    ('小明', 1, 18),
    ('小红', 1, 17),
    ('小刚', 2, 19),
    ('小美', 2, 18),
    ('小强', NULL, 20);  -- 没有分配班级

-- 插入课程数据
INSERT INTO courses (name, credit) VALUES 
    ('数学', 4),
    ('英语', 3),
    ('物理', 4),
    ('化学', 3);

-- 插入成绩数据
INSERT INTO scores (student_id, course_id, score) VALUES 
    (1, 1, 90),   -- 小明 数学 90
    (1, 2, 85),   -- 小明 英语 85
    (1, 3, 88),   -- 小明 物理 88
    (2, 1, 92),   -- 小红 数学 92
    (2, 2, 78),   -- 小红 英语 78
    (3, 1, 75),   -- 小刚 数学 75
    (3, 3, 82),   -- 小刚 物理 82
    (4, 2, 95);   -- 小美 英语 95
```

**表关系图：**
```
students (学生)          classes (班级)
+----+------+----------+    +----+------+---------+
| id | name | class_id |    | id | name | teacher |
+----+------+----------+    +----+------+---------+
| 1  | 小明 | 1        | -> | 1  | 一班 | 张老师  |
| 2  | 小红 | 1        | -> | 1  | 一班 | 张老师  |
| 3  | 小刚 | 2        | -> | 2  | 二班 | 李老师  |
| 4  | 小美 | 2        | -> | 2  | 二班 | 李老师  |
| 5  | 小强 | NULL     |    |    |      |         |
+----+------+----------+    +----+------+---------+

scores (成绩)                courses (课程)
+----+------------+-----------+-------+    +----+------+--------+
| id | student_id | course_id | score |    | id | name | credit |
+----+------------+-----------+-------+    +----+------+--------+
| 1  | 1          | 1         | 90    | -> | 1  | 数学 | 4      |
| 2  | 1          | 2         | 85    | -> | 2  | 英语 | 3      |
+----+------------+-----------+-------+    +----+------+--------+
```

---

## 3. JOIN 类型概览

```
+-------------------+-------------------+
|     表A           |       表B         |
|   +-------+       |    +-------+      |
|   |       |       |    |       |      |
|   |   A   | A∩B   |    |   B   |      |
|   |       |       |    |       |      |
|   +-------+       |    +-------+      |
+-------------------+-------------------+

INNER JOIN: 只返回 A∩B（交集）
LEFT JOIN:  返回 A 全部 + A∩B
RIGHT JOIN: 返回 B 全部 + A∩B
FULL JOIN:  返回 A + B 全部（MySQL不直接支持）
```

---

## 4. INNER JOIN（内连接）

**只返回两个表中匹配的记录**

### 4.1 基本语法

```sql
SELECT 列名
FROM 表A
INNER JOIN 表B ON 表A.列 = 表B.列;
```

### 4.2 示例

```sql
-- 查询学生及其班级信息
SELECT 
    s.id,
    s.name AS 学生姓名,
    s.age,
    c.name AS 班级名称,
    c.teacher AS 班主任
FROM students s
INNER JOIN classes c ON s.class_id = c.id;
```

**结果：**
```
+----+----------+-----+----------+---------+
| id | 学生姓名 | age | 班级名称 | 班主任  |
+----+----------+-----+----------+---------+
| 1  | 小明     | 18  | 一班     | 张老师  |
| 2  | 小红     | 17  | 一班     | 张老师  |
| 3  | 小刚     | 19  | 二班     | 李老师  |
| 4  | 小美     | 18  | 二班     | 李老师  |
+----+----------+-----+----------+---------+
```

**注意：** 小强没有出现，因为他的 class_id 是 NULL，无法匹配。

### 4.3 简写形式

```sql
-- 使用 JOIN（默认就是 INNER JOIN）
SELECT s.name, c.name
FROM students s
JOIN classes c ON s.class_id = c.id;

-- 使用 WHERE（旧写法，不推荐）
SELECT s.name, c.name
FROM students s, classes c
WHERE s.class_id = c.id;
```

---

## 5. LEFT JOIN（左连接）

**返回左表所有记录，右表没有匹配的显示 NULL**

### 5.1 基本语法

```sql
SELECT 列名
FROM 表A
LEFT JOIN 表B ON 表A.列 = 表B.列;
```

### 5.2 示例

```sql
-- 查询所有学生及其班级（包括没有班级的学生）
SELECT 
    s.id,
    s.name AS 学生姓名,
    c.name AS 班级名称,
    c.teacher AS 班主任
FROM students s
LEFT JOIN classes c ON s.class_id = c.id;
```

**结果：**
```
+----+----------+----------+---------+
| id | 学生姓名 | 班级名称 | 班主任  |
+----+----------+----------+---------+
| 1  | 小明     | 一班     | 张老师  |
| 2  | 小红     | 一班     | 张老师  |
| 3  | 小刚     | 二班     | 李老师  |
| 4  | 小美     | 二班     | 李老师  |
| 5  | 小强     | NULL     | NULL    |  <- 小强出现了！
+----+----------+----------+---------+
```

### 5.3 查找没有匹配的记录

```sql
-- 查询没有分配班级的学生
SELECT s.name
FROM students s
LEFT JOIN classes c ON s.class_id = c.id
WHERE c.id IS NULL;
```

---

## 6. RIGHT JOIN（右连接）

**返回右表所有记录，左表没有匹配的显示 NULL**

```sql
-- 查询所有班级及其学生（包括没有学生的班级）
SELECT 
    c.name AS 班级名称,
    s.name AS 学生姓名
FROM students s
RIGHT JOIN classes c ON s.class_id = c.id;
```

**结果：**
```
+----------+----------+
| 班级名称 | 学生姓名 |
+----------+----------+
| 一班     | 小明     |
| 一班     | 小红     |
| 二班     | 小刚     |
| 二班     | 小美     |
| 三班     | NULL     |  <- 三班没有学生
+----------+----------+
```

**提示：** RIGHT JOIN 可以通过交换表的顺序用 LEFT JOIN 实现，实际开发中 LEFT JOIN 更常用。

---

## 7. 多表连接

### 7.1 三表连接

```sql
-- 查询学生的成绩详情（学生名、课程名、分数）
SELECT 
    s.name AS 学生,
    c.name AS 课程,
    sc.score AS 分数
FROM students s
INNER JOIN scores sc ON s.id = sc.student_id
INNER JOIN courses c ON sc.course_id = c.id
ORDER BY s.name, c.name;
```

**结果：**
```
+------+------+-------+
| 学生 | 课程 | 分数  |
+------+------+-------+
| 小刚 | 数学 | 75.00 |
| 小刚 | 物理 | 82.00 |
| 小红 | 数学 | 92.00 |
| 小红 | 英语 | 78.00 |
| 小明 | 数学 | 90.00 |
| 小明 | 物理 | 88.00 |
| 小明 | 英语 | 85.00 |
| 小美 | 英语 | 95.00 |
+------+------+-------+
```

### 7.2 四表连接

```sql
-- 查询完整信息：学生、班级、课程、成绩
SELECT 
    s.name AS 学生,
    cl.name AS 班级,
    co.name AS 课程,
    sc.score AS 分数
FROM students s
LEFT JOIN classes cl ON s.class_id = cl.id
LEFT JOIN scores sc ON s.id = sc.student_id
LEFT JOIN courses co ON sc.course_id = co.id
ORDER BY s.name;
```

---

## 8. 自连接

**表与自身连接，常用于层级数据（如员工-上级关系）**

```sql
-- 创建员工表（包含上级ID）
CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    manager_id INT
);

INSERT INTO employees (name, manager_id) VALUES 
    ('老板', NULL),
    ('经理A', 1),
    ('经理B', 1),
    ('员工1', 2),
    ('员工2', 2),
    ('员工3', 3);

-- 查询员工及其上级
SELECT 
    e.name AS 员工,
    m.name AS 上级
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

**结果：**
```
+-------+-------+
| 员工  | 上级  |
+-------+-------+
| 老板  | NULL  |
| 经理A | 老板  |
| 经理B | 老板  |
| 员工1 | 经理A |
| 员工2 | 经理A |
| 员工3 | 经理B |
+-------+-------+
```

---

## 9. 子查询

**在一个查询中嵌套另一个查询**

### 9.1 WHERE 中的子查询

```sql
-- 查询成绩高于平均分的记录
SELECT s.name, sc.score
FROM students s
JOIN scores sc ON s.id = sc.student_id
WHERE sc.score > (SELECT AVG(score) FROM scores);

-- 查询选修了"数学"课程的学生
SELECT name FROM students
WHERE id IN (
    SELECT student_id FROM scores
    WHERE course_id = (SELECT id FROM courses WHERE name = '数学')
);
```

### 9.2 FROM 中的子查询（派生表）

```sql
-- 查询每个学生的平均分，并筛选平均分大于80的
SELECT * FROM (
    SELECT 
        s.name AS 学生,
        AVG(sc.score) AS 平均分
    FROM students s
    JOIN scores sc ON s.id = sc.student_id
    GROUP BY s.id, s.name
) AS avg_scores
WHERE 平均分 > 80;
```

### 9.3 EXISTS 子查询

```sql
-- 查询有成绩记录的学生
SELECT name FROM students s
WHERE EXISTS (
    SELECT 1 FROM scores WHERE student_id = s.id
);

-- 查询没有成绩记录的学生
SELECT name FROM students s
WHERE NOT EXISTS (
    SELECT 1 FROM scores WHERE student_id = s.id
);
```

---

## 10. UNION 合并结果

**将多个 SELECT 结果合并**

```sql
-- UNION：去重合并
SELECT name FROM students WHERE class_id = 1
UNION
SELECT name FROM students WHERE age > 18;

-- UNION ALL：不去重合并（更快）
SELECT name FROM students WHERE class_id = 1
UNION ALL
SELECT name FROM students WHERE age > 18;
```

---

## 11. 实战练习

### 练习1：基础连接
查询所有学生的姓名和班级名称（包括没有班级的学生）

<details>
<summary>点击查看答案</summary>

```sql
SELECT s.name AS 学生姓名, c.name AS 班级名称
FROM students s
LEFT JOIN classes c ON s.class_id = c.id;
```
</details>

---

### 练习2：多表连接
查询每个学生的姓名、班级、选修的课程名称和成绩

<details>
<summary>点击查看答案</summary>

```sql
SELECT 
    s.name AS 学生,
    cl.name AS 班级,
    co.name AS 课程,
    sc.score AS 成绩
FROM students s
LEFT JOIN classes cl ON s.class_id = cl.id
LEFT JOIN scores sc ON s.id = sc.student_id
LEFT JOIN courses co ON sc.course_id = co.id
ORDER BY s.name, co.name;
```
</details>

---

### 练习3：聚合与连接
查询每个班级的学生人数和平均年龄

<details>
<summary>点击查看答案</summary>

```sql
SELECT 
    c.name AS 班级,
    COUNT(s.id) AS 学生人数,
    AVG(s.age) AS 平均年龄
FROM classes c
LEFT JOIN students s ON c.id = s.class_id
GROUP BY c.id, c.name;
```
</details>

---

### 练习4：子查询
查询数学成绩最高的学生姓名和分数

<details>
<summary>点击查看答案</summary>

```sql
-- 方法1：子查询
SELECT s.name, sc.score
FROM students s
JOIN scores sc ON s.id = sc.student_id
WHERE sc.course_id = (SELECT id FROM courses WHERE name = '数学')
AND sc.score = (
    SELECT MAX(score) FROM scores 
    WHERE course_id = (SELECT id FROM courses WHERE name = '数学')
);

-- 方法2：ORDER BY + LIMIT
SELECT s.name, sc.score
FROM students s
JOIN scores sc ON s.id = sc.student_id
JOIN courses c ON sc.course_id = c.id
WHERE c.name = '数学'
ORDER BY sc.score DESC
LIMIT 1;
```
</details>

---

### 练习5：复杂查询
查询每门课程的最高分、最低分、平均分，并显示课程名称

<details>
<summary>点击查看答案</summary>

```sql
SELECT 
    c.name AS 课程,
    MAX(sc.score) AS 最高分,
    MIN(sc.score) AS 最低分,
    ROUND(AVG(sc.score), 2) AS 平均分,
    COUNT(sc.id) AS 选课人数
FROM courses c
LEFT JOIN scores sc ON c.id = sc.course_id
GROUP BY c.id, c.name;
```
</details>

---

### 练习6：查找未选课学生
查询没有选修任何课程的学生

<details>
<summary>点击查看答案</summary>

```sql
-- 方法1：LEFT JOIN + IS NULL
SELECT s.name
FROM students s
LEFT JOIN scores sc ON s.id = sc.student_id
WHERE sc.id IS NULL;

-- 方法2：NOT EXISTS
SELECT name FROM students s
WHERE NOT EXISTS (
    SELECT 1 FROM scores WHERE student_id = s.id
);

-- 方法3：NOT IN
SELECT name FROM students
WHERE id NOT IN (SELECT DISTINCT student_id FROM scores);
```
</details>

---

### 练习7：综合查询
查询一班学生的所有成绩，显示学生姓名、课程名、成绩，按成绩降序排列

<details>
<summary>点击查看答案</summary>

```sql
SELECT 
    s.name AS 学生,
    co.name AS 课程,
    sc.score AS 成绩
FROM students s
JOIN classes cl ON s.class_id = cl.id
JOIN scores sc ON s.id = sc.student_id
JOIN courses co ON sc.course_id = co.id
WHERE cl.name = '一班'
ORDER BY sc.score DESC;
```
</details>

---

### 练习8：统计分析
查询每个学生的总分和平均分，只显示平均分大于80的学生，按总分降序排列

<details>
<summary>点击查看答案</summary>

```sql
SELECT 
    s.name AS 学生,
    SUM(sc.score) AS 总分,
    ROUND(AVG(sc.score), 2) AS 平均分,
    COUNT(sc.id) AS 选课数
FROM students s
JOIN scores sc ON s.id = sc.student_id
GROUP BY s.id, s.name
HAVING AVG(sc.score) > 80
ORDER BY 总分 DESC;
```
</details>

---

## 12. JOIN 性能优化建议

1. **确保连接字段有索引**
   ```sql
   CREATE INDEX idx_class_id ON students(class_id);
   CREATE INDEX idx_student_id ON scores(student_id);
   ```

2. **小表驱动大表**：把数据量小的表放在前面

3. **避免 SELECT ***：只查询需要的字段

4. **使用 EXPLAIN 分析查询**
   ```sql
   EXPLAIN SELECT s.name, c.name
   FROM students s
   JOIN classes c ON s.class_id = c.id;
   ```

---

## 13. 本章小结

| JOIN 类型 | 说明 |
|-----------|------|
| `INNER JOIN` | 返回两表匹配的记录 |
| `LEFT JOIN` | 返回左表全部 + 匹配记录 |
| `RIGHT JOIN` | 返回右表全部 + 匹配记录 |
| `SELF JOIN` | 表与自身连接 |

**子查询位置：**
- `WHERE` 子句中
- `FROM` 子句中（派生表）
- `SELECT` 子句中

**下一章预告：** 索引与性能优化
