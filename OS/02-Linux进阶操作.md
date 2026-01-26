# Linux 进阶操作

> **重要程度：⭐⭐⭐ 必须掌握**  
> 权限管理、进程管理、网络操作是日常运维必备技能！不会这些，你的程序部署上去也跑不起来！

## 📚 本章学习目标

学完本章，你将能够：
- 理解 Linux 用户和权限系统
- 熟练管理进程（查看、启动、终止）
- 掌握网络相关命令（端口、连接、测试）
- 学会使用 systemd 管理服务
- 理解环境变量的使用
- 配置定时任务

---

## 1. 用户和权限 ⭐⭐⭐

### 1.1 为什么要学用户和权限？

```
想象一下这个场景：

你在公司服务器上部署了一个 Go 应用
├── 你用 root 用户部署的
├── 配置文件里有数据库密码
├── 日志文件里有用户数据
└── 如果任何人都能访问这些文件...💀

所以我们需要：
1. 不同的用户有不同的权限
2. 文件只能被特定用户访问
3. 程序只能做它该做的事
```

### 1.2 Linux 用户体系

```
Linux 用户分为三类：

┌─────────────────────────────────────────────────────┐
│                    root 用户                         │
│                   (超级管理员)                        │
│              拥有系统的所有权限                        │
│              UID = 0                                 │
└─────────────────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ↓               ↓               ↓
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  系统用户     │ │  普通用户     │ │  普通用户     │
│  (daemon)    │ │  (user1)     │ │  (user2)     │
│  运行服务用   │ │  日常使用     │ │  日常使用     │
│  UID 1-999   │ │  UID >= 1000 │ │  UID >= 1000 │
└──────────────┘ └──────────────┘ └──────────────┘

重要概念：
- UID：用户 ID，每个用户唯一
- GID：组 ID，用户可以属于多个组
- root 的 UID 是 0，拥有最高权限
```

### 1.3 用户管理命令

#### 查看当前用户
```bash
# whoami - 我是谁？
$ whoami
root

# id - 查看用户详细信息
$ id
uid=0(root) gid=0(root) groups=0(root)

# id 用户名 - 查看指定用户
$ id nginx
uid=101(nginx) gid=101(nginx) groups=101(nginx)
```

#### 切换用户
```bash
# su - 切换用户（Switch User）
su - username
#  └─ 加 - 表示同时切换环境变量

# 例子：切换到 nginx 用户
$ su - nginx
$ whoami
nginx

# 切换回 root
$ exit
# 或
$ su - root

# 注意：普通用户切换到其他用户需要输入目标用户的密码
# root 切换到任何用户不需要密码
```

#### sudo - 以管理员身份执行
```bash
# sudo = Super User Do（以超级用户身份执行）

# 为什么用 sudo？
# 1. 不用切换到 root 用户
# 2. 有操作记录，更安全
# 3. 可以精细控制权限

# 用法
sudo 命令

# 例子：
$ sudo apt update           # 更新软件包
$ sudo systemctl restart nginx  # 重启服务
$ sudo vim /etc/hosts       # 编辑系统文件

# 第一次使用会要求输入当前用户的密码
# 密码会缓存一段时间，不用每次都输入
```

#### 创建和删除用户
```bash
# 创建用户
sudo useradd -m username
#            └─ -m 表示创建家目录

# 设置密码
sudo passwd username
# 会提示输入两次密码

# 创建用户并指定 shell
sudo useradd -m -s /bin/bash username

# 删除用户
sudo userdel username        # 只删除用户
sudo userdel -r username     # 删除用户和家目录
```

### 📝 练习题 1.1

**问题**：root 用户和 sudo 有什么区别？为什么推荐用 sudo 而不是直接用 root？

<details>
<summary>点击查看答案</summary>

**答案**：

| 对比项 | 直接用 root | 使用 sudo |
|--------|------------|-----------|
| 安全性 | 低，容易误操作 | 高，每次都要确认 |
| 审计 | 无法追踪谁做了什么 | 有日志记录 |
| 权限控制 | 全部权限 | 可以限制只能执行特定命令 |
| 密码 | 需要知道 root 密码 | 用自己的密码 |

**推荐用 sudo 的原因**：
1. **最小权限原则**：只在需要时获取高权限
2. **可追溯**：`/var/log/auth.log` 记录了谁在什么时候执行了什么命令
3. **防止误操作**：每次都要输入 sudo，给你一个思考的机会
4. **团队协作**：可以给不同用户配置不同的 sudo 权限

</details>

---

### 1.4 文件权限详解

#### 权限的表示方式
```
当你执行 ls -l 时，会看到这样的输出：

-rw-r--r--  1  root  root  1234  Jan 7 10:00  main.go
│└──┬───┘     │      │     │         │          │
│   │         │      │     │         │          └─ 文件名
│   │         │      │     │         └─ 修改时间
│   │         │      │     └─ 文件大小（字节）
│   │         │      └─ 所属组
│   │         └─ 所有者
│   └─ 权限位（9个字符）
└─ 文件类型

文件类型：
- = 普通文件
d = 目录
l = 符号链接（快捷方式）
```

#### 权限位详解
```
权限位分为三组，每组三个字符：

-rwxr-xr--
 └┬┘└┬┘└┬┘
  │  │  └─ 其他用户（Other）的权限
  │  └─ 所属组（Group）的权限
  └─ 所有者（Owner）的权限

每组的三个字符含义：
r = Read（读）    - 可以查看文件内容 / 列出目录内容
w = Write（写）   - 可以修改文件 / 在目录中创建删除文件
x = Execute（执行）- 可以运行程序 / 可以进入目录
- = 没有该权限

例子解读：
-rwxr-xr--
 │││││││││
 │││││││└┴─ 其他用户：r-- = 只能读
 │││││└┴─── 组用户：r-x = 能读和执行
 │││└┴┴──── 所有者：rwx = 能读、写、执行
 └┴┴─────── 这是一个普通文件
```

#### 权限的数字表示法
```
每个权限对应一个数字：
r = 4
w = 2
x = 1
- = 0

计算方法：把三个权限的数字加起来

rwx = 4 + 2 + 1 = 7
rw- = 4 + 2 + 0 = 6
r-x = 4 + 0 + 1 = 5
r-- = 4 + 0 + 0 = 4
--- = 0 + 0 + 0 = 0

所以：
-rwxr-xr-- = 754
-rw-r--r-- = 644
-rwxrwxrwx = 777
drwxr-xr-x = 755（目录）
```

### 📝 练习题 1.2

**问题**：把 `-rw-rw-r--` 转换成数字表示法。

<details>
<summary>点击查看答案</summary>

```
-rw-rw-r--

所有者：rw- = 4 + 2 + 0 = 6
组用户：rw- = 4 + 2 + 0 = 6
其他用户：r-- = 4 + 0 + 0 = 4

答案：664
```

</details>

### 📝 练习题 1.3

**问题**：权限 `755` 代表什么？什么类型的文件通常用这个权限？

<details>
<summary>点击查看答案</summary>

```
755 = rwxr-xr-x

所有者：7 = rwx = 读、写、执行
组用户：5 = r-x = 读、执行
其他用户：5 = r-x = 读、执行

通常用于：
1. 可执行程序（如编译后的 Go 程序）
2. 目录（需要 x 权限才能进入目录）
3. Shell 脚本

为什么是 755 而不是 777？
- 不希望其他人修改你的程序
- 但允许其他人运行它
```

</details>

---

### 1.5 修改权限 - chmod

```bash
# chmod = Change Mode（改变模式/权限）

# 方法1：数字方式（推荐）
chmod 755 script.sh     # rwxr-xr-x
chmod 644 config.yaml   # rw-r--r--
chmod 600 id_rsa        # rw-------（私钥文件）
chmod 777 temp/         # rwxrwxrwx（不推荐！）

# 方法2：符号方式
chmod +x script.sh      # 给所有人添加执行权限
chmod -w file.txt       # 移除所有人的写权限
chmod u+x script.sh     # 只给所有者添加执行权限
chmod g+w file.txt      # 给组用户添加写权限
chmod o-r file.txt      # 移除其他用户的读权限

# 符号说明：
# u = user（所有者）
# g = group（组）
# o = other（其他）
# a = all（所有人）
# + = 添加权限
# - = 移除权限
# = = 设置权限

# 递归修改（目录及其所有内容）
chmod -R 755 myproject/
```

### 1.6 修改所有者 - chown

```bash
# chown = Change Owner（改变所有者）

# 修改所有者
chown nginx file.txt

# 修改所有者和组
chown nginx:nginx file.txt
#     └─ 用户:组

# 递归修改
chown -R nginx:nginx /var/www/

# 只修改组
chgrp nginx file.txt
```

### 1.7 工作中的权限场景

```bash
# 场景1：部署 Go 应用
# 编译后的程序需要执行权限
go build -o myapp main.go
chmod 755 myapp
./myapp

# 场景2：配置文件权限
# 配置文件包含敏感信息，只允许所有者读写
chmod 600 config.yaml

# 场景3：SSH 密钥权限
# SSH 对密钥权限有严格要求
chmod 600 ~/.ssh/id_rsa        # 私钥
chmod 644 ~/.ssh/id_rsa.pub    # 公钥
chmod 700 ~/.ssh               # .ssh 目录

# 场景4：Web 目录权限
# 让 nginx 用户能访问网站文件
chown -R nginx:nginx /var/www/html
chmod -R 755 /var/www/html

# 场景5：日志目录权限
# 应用需要写日志
mkdir -p /var/log/myapp
chown myapp:myapp /var/log/myapp
chmod 755 /var/log/myapp
```

### 📝 练习题 1.4

**问题**：你的 Go 程序编译后无法运行，报错 `Permission denied`，如何解决？

<details>
<summary>点击查看答案</summary>

```bash
# 1. 查看当前权限
ls -l myapp
# 输出：-rw-r--r-- 1 user user 2048000 Jan 7 10:00 myapp
# 发现没有执行权限（没有 x）

# 2. 添加执行权限
chmod +x myapp
# 或
chmod 755 myapp

# 3. 再次查看
ls -l myapp
# 输出：-rwxr-xr-x 1 user user 2048000 Jan 7 10:00 myapp

# 4. 运行
./myapp
```

</details>

---

## 2. 进程管理 ⭐⭐⭐

### 2.1 什么是进程？

```
简单理解：

程序 = 躺在硬盘上的代码文件（静态的）
进程 = 正在运行的程序（动态的）

比喻：
- 程序就像菜谱（写在纸上的步骤）
- 进程就像正在做菜（按照菜谱执行）

一个程序可以启动多个进程：
┌─────────────────────────────────────────┐
│              myapp（程序文件）            │
└─────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ↓           ↓           ↓
   ┌─────────┐ ┌─────────┐ ┌─────────┐
   │ 进程1   │ │ 进程2   │ │ 进程3   │
   │ PID=100 │ │ PID=101 │ │ PID=102 │
   └─────────┘ └─────────┘ └─────────┘

每个进程有唯一的 PID（Process ID）
```

### 2.2 查看进程 - ps

```bash
# ps = Process Status（进程状态）

# 最常用的命令：ps aux
$ ps aux
USER   PID  %CPU %MEM    VSZ   RSS TTY  STAT START   TIME COMMAND
root     1   0.0  0.1 169836 13120 ?    Ss   Jan06   0:05 /sbin/init
root   100   0.5  1.2 123456 65432 ?    Sl   10:00   0:30 ./myapp
nginx  200   0.1  0.5  54321 12345 ?    S    10:00   0:10 nginx: worker

# 参数说明：
# a = 显示所有用户的进程
# u = 显示用户信息
# x = 显示没有终端的进程
```

#### 输出字段详解
```
USER   PID  %CPU %MEM    VSZ   RSS TTY  STAT START   TIME COMMAND
│      │    │    │       │     │   │    │    │       │    │
│      │    │    │       │     │   │    │    │       │    └─ 命令/程序名
│      │    │    │       │     │   │    │    │       └─ 累计 CPU 时间
│      │    │    │       │     │   │    │    └─ 启动时间
│      │    │    │       │     │   │    └─ 进程状态
│      │    │    │       │     │   └─ 终端（? 表示没有终端）
│      │    │    │       │     └─ 实际使用的物理内存（KB）
│      │    │    │       └─ 虚拟内存大小（KB）
│      │    │    └─ 内存使用百分比
│      │    └─ CPU 使用百分比
│      └─ 进程 ID
└─ 运行该进程的用户
```

#### 进程状态（STAT）详解
```
常见状态：
R = Running（运行中）      - 正在使用 CPU
S = Sleeping（睡眠）       - 等待某个事件（如等待网络请求）
D = Disk Sleep（磁盘睡眠） - 等待 I/O，不可中断
Z = Zombie（僵尸）         - 已结束但父进程还没回收
T = Stopped（停止）        - 被暂停了

附加标志：
s = 会话领导者
l = 多线程
+ = 前台进程
< = 高优先级
N = 低优先级

例子：
Ss  = 睡眠状态的会话领导者
Sl  = 睡眠状态的多线程进程
R+  = 运行中的前台进程
```

#### 常用的 ps 命令组合
```bash
# 查看所有进程
ps aux

# 查看指定用户的进程
ps -u nginx

# 查看指定进程名
ps aux | grep myapp

# 查看进程树（父子关系）
ps auxf
# 或
pstree

# 只显示 PID 和命令
ps -eo pid,cmd

# 按 CPU 使用率排序
ps aux --sort=-%cpu | head -10

# 按内存使用率排序
ps aux --sort=-%mem | head -10
```

### 📝 练习题 2.1

**问题**：如何找出系统中最占 CPU 的 5 个进程？

<details>
<summary>点击查看答案</summary>

```bash
# 方法1：ps + sort
ps aux --sort=-%cpu | head -6
# head -6 是因为第一行是标题

# 方法2：ps + sort 命令
ps aux | sort -k3 -rn | head -5
# -k3 表示按第3列（CPU）排序
# -r 表示降序
# -n 表示数字排序

# 方法3：使用 top
top -bn1 | head -12
# -b 批处理模式
# -n1 只刷新一次
```

</details>

---

### 2.3 实时查看进程 - top

```bash
# top 是实时的进程监控工具

$ top

# 输出示例：
top - 10:30:00 up 5 days,  2:30,  1 user,  load average: 0.15, 0.10, 0.05
Tasks: 150 total,   1 running, 149 sleeping,   0 stopped,   0 zombie
%Cpu(s):  2.0 us,  1.0 sy,  0.0 ni, 96.5 id,  0.5 wa,  0.0 hi,  0.0 si
MiB Mem :   7982.4 total,   3245.6 free,   2156.8 used,   2580.0 buff/cache
MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   5432.1 avail Mem

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
  100 root      20   0  123456  65432  12345 S   5.0   0.8   0:30.00 myapp
  200 nginx     20   0   54321  12345   6789 S   1.0   0.2   0:10.00 nginx
```

#### top 界面解读
```
第一行：系统运行时间和负载
top - 10:30:00 up 5 days,  2:30,  1 user,  load average: 0.15, 0.10, 0.05
       │            │           │                  │
       │            │           │                  └─ 负载（1分钟、5分钟、15分钟）
       │            │           └─ 登录用户数
       │            └─ 运行时间
       └─ 当前时间

负载（load average）怎么看？
- 数值 = 等待 CPU 的进程数
- 单核 CPU：负载 1.0 表示 CPU 刚好满载
- 4 核 CPU：负载 4.0 表示 CPU 刚好满载
- 负载持续高于核心数 = CPU 不够用了

第三行：CPU 使用情况
%Cpu(s):  2.0 us,  1.0 sy,  0.0 ni, 96.5 id,  0.5 wa
          │        │        │       │        │
          │        │        │       │        └─ 等待 I/O
          │        │        │       └─ 空闲
          │        │        └─ nice 值调整的进程
          │        └─ 系统（内核）
          └─ 用户空间

重点关注：
- us 高 = 应用程序占用 CPU 多
- sy 高 = 系统调用多（可能有问题）
- wa 高 = I/O 等待多（磁盘慢或网络慢）
- id 高 = CPU 空闲（正常）
```

#### top 中的快捷键
```bash
# 在 top 界面中：
q       # 退出
h       # 帮助
M       # 按内存排序
P       # 按 CPU 排序
k       # 终止进程（输入 PID）
1       # 显示每个 CPU 核心的使用情况
c       # 显示完整命令
```

### 2.4 更好用的 htop

```bash
# htop 是 top 的增强版，更直观

# 安装
sudo apt install htop    # Ubuntu/Debian
sudo yum install htop    # CentOS

# 使用
htop

# htop 的优点：
# 1. 彩色显示，更直观
# 2. 可以用鼠标操作
# 3. 可以直接用方向键选择进程
# 4. 按 F9 可以发送信号终止进程
# 5. 按 F5 可以看进程树
```

### 📝 练习题 2.2

**问题**：`load average: 2.50, 1.80, 1.20` 在一台 2 核 CPU 的服务器上意味着什么？

<details>
<summary>点击查看答案</summary>

**答案**：

```
load average: 2.50, 1.80, 1.20
              │     │     │
              │     │     └─ 15分钟平均负载：1.20
              │     └─ 5分钟平均负载：1.80
              └─ 1分钟平均负载：2.50

对于 2 核 CPU：
- 负载 2.0 = CPU 刚好满载
- 负载 2.50 = 超载 25%（有进程在排队等待 CPU）

分析：
- 1分钟负载 2.50 > 2.0：当前 CPU 有点忙
- 5分钟负载 1.80 < 2.0：最近 5 分钟还好
- 15分钟负载 1.20 < 2.0：长期来看负载不高

结论：可能是短期的负载高峰，需要继续观察。
如果持续高于 2.0，就需要考虑优化或扩容了。
```

</details>

---

### 2.5 终止进程 - kill

```bash
# kill 命令用于向进程发送信号

# 基本用法
kill PID

# 常用信号：
kill -15 PID    # SIGTERM：优雅终止（默认）
kill PID        # 等同于 kill -15
kill -9 PID     # SIGKILL：强制终止（无法被忽略）
kill -1 PID     # SIGHUP：重新加载配置

# 信号说明：
# SIGTERM (15)：请求进程终止，进程可以清理资源后退出
# SIGKILL (9)：强制终止，进程无法捕获，立即死亡
# SIGHUP (1)：挂起信号，很多服务用它来重新加载配置
```

#### 为什么要先用 SIGTERM？
```
SIGTERM（优雅终止）：
┌─────────────────────────────────────────┐
│ 1. 进程收到信号                          │
│ 2. 进程完成当前请求                       │
│ 3. 进程关闭数据库连接                     │
│ 4. 进程保存状态                          │
│ 5. 进程退出                              │
└─────────────────────────────────────────┘

SIGKILL（强制终止）：
┌─────────────────────────────────────────┐
│ 1. 进程收到信号                          │
│ 2. 进程立即死亡                          │
│    - 当前请求中断                        │
│    - 数据库连接没关闭                     │
│    - 数据可能丢失                        │
└─────────────────────────────────────────┘

所以：先用 kill，等几秒，不行再用 kill -9
```

#### 其他终止进程的方式
```bash
# 按进程名终止
pkill myapp
killall myapp

# 查找并终止
ps aux | grep myapp | awk '{print $2}' | xargs kill

# 终止某个端口的进程
kill $(lsof -t -i:8080)
# -t 只输出 PID
```

### 📝 练习题 2.3

**问题**：你的 Go 程序卡死了，`kill PID` 没反应，怎么办？

<details>
<summary>点击查看答案</summary>

```bash
# 1. 先确认进程还在
ps aux | grep myapp

# 2. 尝试 SIGTERM（可能已经试过了）
kill PID

# 3. 等待几秒，如果还没退出，使用 SIGKILL
kill -9 PID

# 4. 确认进程已终止
ps aux | grep myapp

# 注意：
# - kill -9 是最后手段
# - 使用前要考虑数据丢失的风险
# - 如果程序经常需要 kill -9，说明程序有 bug
```

</details>

---

### 2.6 后台运行进程

#### 方法1：& 符号
```bash
# 在命令后加 & 可以让程序在后台运行
./myapp &

# 问题：SSH 断开后，进程会收到 SIGHUP 信号而终止
```

#### 方法2：nohup（推荐）
```bash
# nohup = No Hang Up（不挂断）
# 让进程忽略 SIGHUP 信号

nohup ./myapp &
# 输出会写入 nohup.out 文件

# 指定输出文件
nohup ./myapp > app.log 2>&1 &
#                       │
#                       └─ 2>&1 表示把错误输出也重定向到 app.log

# 查看输出
tail -f app.log
```

#### 方法3：screen 或 tmux
```bash
# screen 可以创建一个虚拟终端，断开后还能重新连接

# 安装
sudo apt install screen

# 创建新会话
screen -S myapp

# 在 screen 中运行程序
./myapp

# 断开会话（程序继续运行）
Ctrl + A, D

# 重新连接
screen -r myapp

# 列出所有会话
screen -ls
```

### 📝 练习题 2.4

**问题**：`nohup ./myapp > app.log 2>&1 &` 这个命令中，`2>&1` 是什么意思？

<details>
<summary>点击查看答案</summary>

**答案**：

```
在 Linux 中，有三个标准文件描述符：
0 = stdin（标准输入）
1 = stdout（标准输出）- 正常输出
2 = stderr（标准错误）- 错误输出

2>&1 的含义：
2  = stderr（错误输出）
>  = 重定向
&1 = 指向 stdout（标准输出）

所以 2>&1 = 把错误输出重定向到标准输出

完整命令解析：
nohup ./myapp > app.log 2>&1 &
│     │       │         │    │
│     │       │         │    └─ 后台运行
│     │       │         └─ 错误输出也写入 app.log
│     │       └─ 标准输出写入 app.log
│     └─ 运行的程序
└─ 忽略 SIGHUP 信号

如果不加 2>&1：
- 正常输出会写入 app.log
- 错误输出会显示在终端（或丢失）
```

</details>

---

## 3. systemd 服务管理 ⭐⭐⭐

### 3.1 什么是 systemd？

```
systemd 是现代 Linux 的服务管理系统

以前（SysV init）：
- 用脚本管理服务
- 启动顺序固定
- 管理复杂

现在（systemd）：
- 统一的服务管理
- 并行启动，更快
- 自动处理依赖关系
- 自动重启崩溃的服务

你的 Go 应用应该用 systemd 来管理！
```

### 3.2 systemctl 基本命令

```bash
# systemctl = system control（系统控制）

# 启动服务
sudo systemctl start nginx

# 停止服务
sudo systemctl stop nginx

# 重启服务
sudo systemctl restart nginx

# 重新加载配置（不重启）
sudo systemctl reload nginx

# 查看状态
sudo systemctl status nginx

# 开机自启
sudo systemctl enable nginx

# 取消开机自启
sudo systemctl disable nginx

# 查看是否开机自启
sudo systemctl is-enabled nginx

# 查看所有服务
sudo systemctl list-units --type=service

# 查看失败的服务
sudo systemctl --failed
```

### 3.3 查看服务状态

```bash
$ sudo systemctl status nginx

● nginx.service - A high performance web server
     Loaded: loaded (/lib/systemd/system/nginx.service; enabled; vendor preset: enabled)
     Active: active (running) since Mon 2024-01-07 10:00:00 UTC; 2h ago
       Docs: man:nginx(8)
   Main PID: 1234 (nginx)
      Tasks: 3 (limit: 4915)
     Memory: 10.5M
        CPU: 1.234s
     CGroup: /system.slice/nginx.service
             ├─1234 nginx: master process /usr/sbin/nginx
             └─1235 nginx: worker process

Jan 07 10:00:00 server systemd[1]: Starting A high performance web server...
Jan 07 10:00:00 server systemd[1]: Started A high performance web server.
```

#### 状态解读
```
Loaded: loaded (...; enabled; ...)
        │            │
        │            └─ enabled = 开机自启
        └─ 配置文件已加载

Active: active (running)
        │       │
        │       └─ 运行中
        └─ 活跃状态

常见状态：
- active (running)：正在运行
- active (exited)：已执行完毕（一次性任务）
- inactive (dead)：未运行
- failed：启动失败
```

### 3.4 为 Go 应用创建 systemd 服务

#### 第一步：创建服务文件
```bash
sudo vim /etc/systemd/system/myapp.service
```

#### 第二步：编写服务配置
```ini
[Unit]
Description=My Go Application
Documentation=https://github.com/yourname/myapp
After=network.target mysql.service redis.service
Wants=mysql.service redis.service

[Service]
Type=simple
User=myapp
Group=myapp
WorkingDirectory=/home/myapp
ExecStart=/home/myapp/myapp
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
LimitNOFILE=65535

# 环境变量
Environment=PORT=8080
Environment=GIN_MODE=release
EnvironmentFile=/home/myapp/.env

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=myapp

[Install]
WantedBy=multi-user.target
```

#### 配置详解
```ini
[Unit] 部分：服务的描述和依赖
Description=My Go Application     # 服务描述
After=network.target              # 在网络启动后再启动
Wants=mysql.service               # 希望 mysql 也启动（非强制）

[Service] 部分：服务的运行配置
Type=simple                       # 服务类型（simple 最常用）
User=myapp                        # 运行用户（不要用 root！）
Group=myapp                       # 运行组
WorkingDirectory=/home/myapp      # 工作目录
ExecStart=/home/myapp/myapp       # 启动命令
ExecReload=/bin/kill -HUP $MAINPID  # 重载命令
Restart=always                    # 崩溃后自动重启
RestartSec=5                      # 重启间隔（秒）
LimitNOFILE=65535                 # 最大文件描述符数

[Install] 部分：安装配置
WantedBy=multi-user.target        # 多用户模式下启动
```

#### 第三步：启用服务
```bash
# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start myapp

# 查看状态
sudo systemctl status myapp

# 设置开机自启
sudo systemctl enable myapp

# 查看日志
sudo journalctl -u myapp -f
```

### 📝 练习题 3.1

**问题**：为什么服务配置中要设置 `User=myapp` 而不是 `User=root`？

<details>
<summary>点击查看答案</summary>

**答案**：

**安全原因**：
1. **最小权限原则**：程序只应该有它需要的权限
2. **限制损害**：如果程序被攻击，攻击者只能获得 myapp 用户的权限，而不是 root
3. **防止误操作**：程序 bug 不会影响系统关键文件

**实际例子**：
```
如果用 root 运行：
- 程序可以读写任何文件
- 程序可以删除系统文件
- 程序被攻击 = 整个系统被攻击

如果用 myapp 运行：
- 程序只能访问 myapp 用户的文件
- 程序无法修改系统配置
- 程序被攻击 = 只影响 myapp 的数据
```

**最佳实践**：
```bash
# 创建专用用户
sudo useradd -r -s /bin/false myapp
# -r 创建系统用户
# -s /bin/false 禁止登录

# 设置目录权限
sudo chown -R myapp:myapp /home/myapp
```

</details>

### 📝 练习题 3.2

**问题**：你的 Go 服务启动失败了，如何排查？

<details>
<summary>点击查看答案</summary>

```bash
# 1. 查看服务状态
sudo systemctl status myapp
# 看 Active 行是否显示 failed
# 看最后几行日志

# 2. 查看详细日志
sudo journalctl -u myapp -n 50
# -n 50 显示最后 50 行

# 3. 查看实时日志
sudo journalctl -u myapp -f

# 4. 常见问题排查：

# 问题1：权限不足
# 检查可执行文件权限
ls -l /home/myapp/myapp
# 检查工作目录权限
ls -ld /home/myapp

# 问题2：端口被占用
ss -tlnp | grep 8080

# 问题3：配置文件错误
# 手动运行程序看报错
/home/myapp/myapp

# 问题4：依赖服务没启动
sudo systemctl status mysql
sudo systemctl status redis

# 5. 修复后重启
sudo systemctl daemon-reload  # 如果改了服务文件
sudo systemctl restart myapp
```

</details>

---

### 3.5 journalctl 查看日志

```bash
# journalctl 是 systemd 的日志查看工具

# 查看指定服务的日志
journalctl -u myapp

# 实时查看（类似 tail -f）
journalctl -u myapp -f

# 查看最后 100 行
journalctl -u myapp -n 100

# 查看今天的日志
journalctl -u myapp --since today

# 查看最近 1 小时的日志
journalctl -u myapp --since "1 hour ago"

# 查看指定时间范围
journalctl -u myapp --since "2024-01-07 10:00" --until "2024-01-07 12:00"

# 查看错误级别的日志
journalctl -u myapp -p err

# 输出为 JSON 格式
journalctl -u myapp -o json
```

---

## 4. 网络操作 ⭐⭐⭐

### 4.1 查看网络信息

```bash
# 查看 IP 地址
$ ip addr
# 或简写
$ ip a

# 输出示例：
1: lo: <LOOPBACK,UP,LOWER_UP>
    inet 127.0.0.1/8 scope host lo
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP>
    inet 192.168.1.100/24 brd 192.168.1.255 scope global eth0

# 解读：
# lo = loopback（本地回环），IP 是 127.0.0.1
# eth0 = 网卡，IP 是 192.168.1.100

# 旧命令（有些系统还在用）
ifconfig
```

### 4.2 查看端口和连接 ⭐⭐⭐

```bash
# ss 命令（推荐，比 netstat 快）
# ss = Socket Statistics

# 查看所有 TCP 监听端口
ss -tlnp
# -t = TCP
# -l = 监听状态
# -n = 显示数字（不解析域名）
# -p = 显示进程

# 输出示例：
State   Recv-Q  Send-Q  Local Address:Port  Peer Address:Port  Process
LISTEN  0       128     0.0.0.0:8080        0.0.0.0:*          users:(("myapp",pid=1234,fd=3))
LISTEN  0       128     0.0.0.0:22          0.0.0.0:*          users:(("sshd",pid=100,fd=3))
LISTEN  0       128     127.0.0.1:3306      0.0.0.0:*          users:(("mysqld",pid=200,fd=3))
```

#### 输出解读
```
Local Address:Port
│              │
│              └─ 监听的端口
└─ 监听的地址

0.0.0.0:8080   = 监听所有网卡的 8080 端口（外部可访问）
127.0.0.1:3306 = 只监听本地的 3306 端口（外部无法访问）
:::8080        = IPv6 的所有地址
```

#### 常用命令
```bash
# 查看指定端口
ss -tlnp | grep 8080

# 查看谁占用了端口
lsof -i :8080

# 查看所有连接（包括已建立的）
ss -anp

# 查看 TCP 连接状态统计
ss -s

# 旧命令 netstat（功能类似）
netstat -tlnp
netstat -anp
```

### 📝 练习题 4.1

**问题**：你的 Go 应用监听 8080 端口，但外部无法访问，可能是什么原因？

<details>
<summary>点击查看答案</summary>

**排查步骤**：

```bash
# 1. 检查应用是否在运行
ps aux | grep myapp

# 2. 检查端口是否在监听
ss -tlnp | grep 8080

# 如果没有输出，说明应用没有正确监听端口
# 如果输出是 127.0.0.1:8080，说明只监听了本地

# 3. 检查监听地址
# 错误：只监听本地
http.ListenAndServe("127.0.0.1:8080", nil)
# 正确：监听所有地址
http.ListenAndServe(":8080", nil)
http.ListenAndServe("0.0.0.0:8080", nil)

# 4. 检查防火墙
sudo ufw status
# 如果 8080 没有开放
sudo ufw allow 8080

# 5. 检查云服务器安全组
# 阿里云/腾讯云等需要在控制台开放端口

# 6. 本地测试
curl http://localhost:8080
# 如果本地能访问，说明是网络/防火墙问题
```

</details>

---

### 4.3 网络测试命令

#### ping - 测试连通性
```bash
# ping 测试网络是否通
ping google.com

# 只 ping 4 次
ping -c 4 google.com

# 输出示例：
PING google.com (142.250.185.78): 56 data bytes
64 bytes from 142.250.185.78: icmp_seq=0 ttl=116 time=10.5 ms
64 bytes from 142.250.185.78: icmp_seq=1 ttl=116 time=11.2 ms

# time = 延迟（越小越好）
# ttl = 生存时间
```

#### curl - 测试 HTTP
```bash
# curl 是最常用的 HTTP 测试工具

# 基本请求
curl http://localhost:8080

# 只看响应头
curl -I http://localhost:8080

# 显示详细信息
curl -v http://localhost:8080

# POST 请求
curl -X POST http://localhost:8080/api/users

# POST JSON 数据
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"test","age":18}' \
  http://localhost:8080/api/users

# 带 Token 的请求
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/api/profile

# 下载文件
curl -O https://example.com/file.zip

# 跟随重定向
curl -L http://example.com

# 设置超时
curl --connect-timeout 5 http://localhost:8080
```

#### telnet/nc - 测试端口
```bash
# 测试端口是否开放
telnet localhost 8080
# 如果连接成功，说明端口开放

# nc（netcat）更好用
nc -zv localhost 8080
# -z = 只扫描，不发送数据
# -v = 详细输出

# 输出：
# Connection to localhost 8080 port [tcp/*] succeeded!
```

### 📝 练习题 4.2

**问题**：用 curl 测试一个 POST 接口，发送 JSON 数据 `{"username":"admin","password":"123456"}`

<details>
<summary>点击查看答案</summary>

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  http://localhost:8080/api/login

# 解释：
# -X POST：指定 HTTP 方法为 POST
# -H "Content-Type: application/json"：设置请求头，告诉服务器发送的是 JSON
# -d '...'：请求体数据

# 如果想看响应头：
curl -i -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  http://localhost:8080/api/login

# 如果想格式化 JSON 输出（需要安装 jq）：
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  http://localhost:8080/api/login | jq
```

</details>

---

### 4.4 防火墙配置

#### Ubuntu (ufw)
```bash
# ufw = Uncomplicated Firewall（简单防火墙）

# 查看状态
sudo ufw status

# 启用防火墙
sudo ufw enable

# 禁用防火墙
sudo ufw disable

# 开放端口
sudo ufw allow 8080
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443

# 开放指定协议的端口
sudo ufw allow 8080/tcp
sudo ufw allow 53/udp

# 删除规则
sudo ufw delete allow 8080

# 拒绝端口
sudo ufw deny 3306

# 允许特定 IP
sudo ufw allow from 192.168.1.100

# 查看详细规则
sudo ufw status verbose
```

#### CentOS (firewalld)
```bash
# 查看状态
sudo firewall-cmd --state
sudo firewall-cmd --list-all

# 开放端口
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload

# 关闭端口
sudo firewall-cmd --remove-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

---

## 5. 系统资源监控 ⭐⭐

### 5.1 内存使用 - free

```bash
$ free -h
# -h = human readable（人类可读格式）

              total        used        free      shared  buff/cache   available
Mem:          7.6Gi       2.1Gi       3.2Gi       100Mi       2.3Gi       5.2Gi
Swap:         2.0Gi          0B       2.0Gi
```

#### 输出解读
```
total     = 总内存（7.6GB）
used      = 已使用（2.1GB）
free      = 完全空闲（3.2GB）
shared    = 共享内存
buff/cache = 缓存（Linux 会用空闲内存做缓存，需要时会释放）
available = 实际可用（5.2GB）= free + 可释放的 buff/cache

重点：看 available，不是 free！
Linux 会把空闲内存用作缓存，这是正常的。
available 才是真正可用的内存。

Swap = 交换分区（把硬盘当内存用，很慢）
如果 Swap 使用很多，说明内存不够了！
```

### 5.2 磁盘使用 - df

```bash
$ df -h
# -h = human readable

Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        50G   20G   28G  42% /
/dev/sda2       100G   60G   35G  64% /home
tmpfs           3.8G     0  3.8G   0% /dev/shm
```

#### 输出解读
```
Filesystem = 文件系统（磁盘分区）
Size       = 总大小
Used       = 已使用
Avail      = 可用
Use%       = 使用百分比
Mounted on = 挂载点

重点关注：
- / 根分区不要超过 80%
- 日志目录 /var/log 容易满
- 数据目录要定期检查
```

### 5.3 目录大小 - du

```bash
# du = Disk Usage（磁盘使用）

# 查看目录大小
du -sh /home/myapp
# -s = summary（汇总）
# -h = human readable

# 查看当前目录下各文件夹大小
du -sh *

# 查看并排序（找出最大的目录）
du -sh * | sort -h

# 查看前 10 大目录
du -sh * | sort -rh | head -10

# 查看指定深度
du -h --max-depth=1 /home
```

### 📝 练习题 5.1

**问题**：服务器磁盘快满了，如何找出占用空间最大的目录？

<details>
<summary>点击查看答案</summary>

```bash
# 1. 先看整体使用情况
df -h

# 2. 找出根目录下最大的目录
du -sh /* | sort -rh | head -10

# 3. 进入最大的目录继续查找
cd /var
du -sh * | sort -rh | head -10

# 4. 常见的大目录：
# /var/log - 日志文件
# /var/lib/docker - Docker 镜像和容器
# /home - 用户数据
# /tmp - 临时文件

# 5. 清理日志
# 查看日志大小
du -sh /var/log/*

# 清理旧日志（谨慎操作）
sudo journalctl --vacuum-size=500M  # 只保留 500MB 日志
sudo find /var/log -name "*.log" -mtime +30 -delete  # 删除 30 天前的日志

# 6. 清理 Docker（如果用 Docker）
docker system prune -a  # 清理未使用的镜像、容器、网络
```

</details>

---

### 5.4 CPU 信息

```bash
# 查看 CPU 信息
cat /proc/cpuinfo | grep "model name" | head -1
# 输出：model name : Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz

# 查看 CPU 核心数
nproc
# 输出：4

# 或者
cat /proc/cpuinfo | grep processor | wc -l

# 查看 CPU 使用率
top -bn1 | grep "Cpu(s)"
# 或用 mpstat（需要安装 sysstat）
mpstat 1 5  # 每秒采样，共 5 次
```

### 5.5 系统信息

```bash
# 查看系统版本
cat /etc/os-release

# 查看内核版本
uname -r
# 输出：5.4.0-100-generic

# 查看所有系统信息
uname -a

# 查看运行时间
uptime
# 输出：10:30:00 up 5 days, 2:30, 1 user, load average: 0.15, 0.10, 0.05

# 查看主机名
hostname
```

---

## 6. 环境变量 ⭐⭐

### 6.1 什么是环境变量？

```
环境变量 = 系统级别的配置

比喻：
- 环境变量就像是"全局设置"
- 所有程序都可以读取
- 用来配置程序的行为

常见用途：
- 配置数据库连接
- 设置 API 密钥
- 指定运行模式（开发/生产）
- 配置路径
```

### 6.2 查看环境变量

```bash
# 查看所有环境变量
env
# 或
printenv

# 查看指定变量
echo $PATH
echo $HOME
echo $USER

# 常见环境变量：
# PATH   - 可执行文件搜索路径
# HOME   - 用户家目录
# USER   - 当前用户名
# SHELL  - 当前 shell
# LANG   - 语言设置
```

### 6.3 设置环境变量

```bash
# 方法1：临时设置（只在当前终端有效）
export MY_VAR="hello"
export PORT=8080

# 验证
echo $MY_VAR
# 输出：hello

# 方法2：运行时设置（只对这个命令有效）
PORT=9090 ./myapp
# myapp 会读取到 PORT=9090
# 但当前终端的 PORT 不变

# 方法3：永久设置（当前用户）
echo 'export MY_VAR="hello"' >> ~/.bashrc
source ~/.bashrc  # 立即生效

# 方法4：永久设置（所有用户）
sudo echo 'export MY_VAR="hello"' >> /etc/profile
source /etc/profile
```

### 6.4 Go 程序读取环境变量

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // 读取环境变量
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"  // 默认值
    }
    
    // 读取数据库配置
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    
    fmt.Printf("Server starting on port %s\n", port)
    fmt.Printf("Database: %s@%s\n", dbUser, dbHost)
    
    // 设置环境变量（很少用）
    os.Setenv("MY_VAR", "hello")
}
```

### 6.5 使用 .env 文件

```bash
# 创建 .env 文件
cat > .env << EOF
PORT=8080
DB_HOST=localhost
DB_USER=root
DB_PASS=secret
GIN_MODE=release
EOF

# 方法1：在启动前加载
export $(cat .env | xargs)
./myapp

# 方法2：在 systemd 服务中使用
# /etc/systemd/system/myapp.service
[Service]
EnvironmentFile=/home/myapp/.env
ExecStart=/home/myapp/myapp

# 方法3：在 Go 中使用 godotenv 库
# go get github.com/joho/godotenv
```

```go
// 使用 godotenv 加载 .env 文件
package main

import (
    "log"
    "os"
    
    "github.com/joho/godotenv"
)

func main() {
    // 加载 .env 文件
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }
    
    port := os.Getenv("PORT")
    // ...
}
```

### 📝 练习题 6.1

**问题**：为什么不应该把数据库密码直接写在代码里，而要用环境变量？

<details>
<summary>点击查看答案</summary>

**答案**：

**安全原因**：
1. **代码会提交到 Git**：密码会被所有人看到
2. **代码可能泄露**：GitHub 上有很多泄露的密码
3. **不同环境不同密码**：开发、测试、生产环境的密码不同

**使用环境变量的好处**：
1. **密码不在代码里**：代码可以公开
2. **方便切换环境**：改环境变量就行，不用改代码
3. **权限控制**：只有运维人员知道生产环境的密码
4. **符合 12-Factor App 原则**

**最佳实践**：
```go
// ❌ 错误做法
db, _ := sql.Open("mysql", "root:password123@tcp(localhost:3306)/mydb")

// ✅ 正确做法
dbUser := os.Getenv("DB_USER")
dbPass := os.Getenv("DB_PASS")
dbHost := os.Getenv("DB_HOST")
dbName := os.Getenv("DB_NAME")
dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName)
db, _ := sql.Open("mysql", dsn)
```

</details>

---

## 7. 定时任务 - crontab ⭐⭐

### 7.1 什么是 crontab？

```
crontab = cron table（定时任务表）

用途：
- 定时备份数据库
- 定时清理日志
- 定时发送报告
- 定时检查服务状态
```

### 7.2 crontab 命令

```bash
# 编辑定时任务
crontab -e

# 查看定时任务
crontab -l

# 删除所有定时任务
crontab -r

# 编辑其他用户的定时任务（需要 root）
sudo crontab -u nginx -e
```

### 7.3 crontab 时间格式

```
┌───────────── 分钟 (0 - 59)
│ ┌───────────── 小时 (0 - 23)
│ │ ┌───────────── 日 (1 - 31)
│ │ │ ┌───────────── 月 (1 - 12)
│ │ │ │ ┌───────────── 星期 (0 - 7，0 和 7 都是周日)
│ │ │ │ │
│ │ │ │ │
* * * * * 命令

特殊字符：
*     任意值
,     列表（1,3,5 = 1、3、5）
-     范围（1-5 = 1到5）
/     步长（*/5 = 每5个单位）
```

### 7.4 常用示例

```bash
# 每分钟执行
* * * * * /home/myapp/check.sh

# 每小时执行（每小时的第 0 分钟）
0 * * * * /home/myapp/hourly.sh

# 每天凌晨 2 点执行
0 2 * * * /home/myapp/backup.sh

# 每天凌晨 2 点 30 分执行
30 2 * * * /home/myapp/backup.sh

# 每周一凌晨 3 点执行
0 3 * * 1 /home/myapp/weekly.sh

# 每月 1 号凌晨 4 点执行
0 4 1 * * /home/myapp/monthly.sh

# 每 5 分钟执行
*/5 * * * * /home/myapp/check.sh

# 每 30 分钟执行
*/30 * * * * /home/myapp/check.sh

# 工作日（周一到周五）每天 9 点执行
0 9 * * 1-5 /home/myapp/workday.sh

# 每天 9 点和 18 点执行
0 9,18 * * * /home/myapp/twice.sh
```

### 7.5 实际工作中的定时任务

```bash
# 1. 每天凌晨备份数据库
0 2 * * * /home/myapp/scripts/backup_db.sh >> /var/log/backup.log 2>&1

# 2. 每小时清理临时文件
0 * * * * find /tmp -type f -mtime +1 -delete

# 3. 每 5 分钟检查服务状态
*/5 * * * * /home/myapp/scripts/health_check.sh

# 4. 每天凌晨清理 7 天前的日志
0 3 * * * find /var/log/myapp -name "*.log" -mtime +7 -delete

# 5. 每周日凌晨重启服务（清理内存）
0 4 * * 0 systemctl restart myapp
```

### 7.6 备份脚本示例

```bash
#!/bin/bash
# /home/myapp/scripts/backup_db.sh

# 配置
DB_USER="root"
DB_PASS="password"
DB_NAME="mydb"
BACKUP_DIR="/home/myapp/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME > $BACKUP_DIR/db_$DATE.sql

# 压缩
gzip $BACKUP_DIR/db_$DATE.sql

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.gz" -mtime +7 -delete

echo "Backup completed: db_$DATE.sql.gz"
```

### 📝 练习题 7.1

**问题**：写一个 crontab 表达式：每天早上 8 点和晚上 8 点执行脚本。

<details>
<summary>点击查看答案</summary>

```bash
0 8,20 * * * /home/myapp/script.sh

# 解释：
# 0     - 第 0 分钟
# 8,20  - 8 点和 20 点
# *     - 每天
# *     - 每月
# *     - 每周

# 或者写成两行：
0 8 * * * /home/myapp/script.sh
0 20 * * * /home/myapp/script.sh
```

</details>

### 📝 练习题 7.2

**问题**：`*/15 * * * *` 表示什么？

<details>
<summary>点击查看答案</summary>

**答案**：每 15 分钟执行一次。

```
*/15 * * * *
│    │ │ │ │
│    │ │ │ └─ 每周
│    │ │ └─ 每月
│    │ └─ 每天
│    └─ 每小时
└─ 每 15 分钟（0, 15, 30, 45 分）

执行时间：
00:00, 00:15, 00:30, 00:45
01:00, 01:15, 01:30, 01:45
...
```

</details>

---

## 8. 实战：完整部署 Go 应用 ⭐⭐⭐

### 8.1 部署前的准备

```bash
# 在本地：
# 1. 确保代码可以正常编译和运行
go build -o myapp main.go
./myapp

# 2. 准备配置文件
cat > config.yaml << EOF
server:
  port: 8080
database:
  host: localhost
  port: 3306
  user: myapp
  name: mydb
EOF

# 3. 交叉编译（在 Mac/Windows 上编译 Linux 版本）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp main.go
```

### 8.2 服务器初始化

```bash
# 登录服务器
ssh root@your-server-ip

# 1. 创建应用用户（不要用 root 运行应用！）
sudo useradd -r -m -s /bin/bash myapp

# 2. 创建目录结构
sudo mkdir -p /home/myapp/{bin,config,logs,data}
sudo chown -R myapp:myapp /home/myapp

# 3. 安装必要软件
sudo apt update
sudo apt install -y nginx mysql-server redis-server

# 4. 配置防火墙
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### 8.3 上传代码

```bash
# 在本地执行：

# 方法1：scp 上传
scp myapp root@your-server-ip:/home/myapp/bin/
scp config.yaml root@your-server-ip:/home/myapp/config/

# 方法2：rsync 同步（推荐）
rsync -avz --exclude='.git' \
  ./myapp \
  ./config.yaml \
  root@your-server-ip:/home/myapp/
```

### 8.4 创建 systemd 服务

```bash
# 在服务器上：
sudo vim /etc/systemd/system/myapp.service
```

```ini
[Unit]
Description=My Go Application
Documentation=https://github.com/yourname/myapp
After=network.target mysql.service redis.service

[Service]
Type=simple
User=myapp
Group=myapp
WorkingDirectory=/home/myapp
ExecStart=/home/myapp/bin/myapp
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
LimitNOFILE=65535

# 环境变量
Environment=GIN_MODE=release
Environment=PORT=8080
EnvironmentFile=-/home/myapp/config/.env

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=myapp

[Install]
WantedBy=multi-user.target
```

### 8.5 配置 Nginx 反向代理

```bash
sudo vim /etc/nginx/sites-available/myapp
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

```bash
# 启用配置
sudo ln -s /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 8.6 启动服务

```bash
# 设置权限
sudo chown -R myapp:myapp /home/myapp
sudo chmod +x /home/myapp/bin/myapp

# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start myapp

# 查看状态
sudo systemctl status myapp

# 设置开机自启
sudo systemctl enable myapp

# 查看日志
sudo journalctl -u myapp -f
```

### 8.7 验证部署

```bash
# 1. 检查进程
ps aux | grep myapp

# 2. 检查端口
ss -tlnp | grep 8080

# 3. 本地测试
curl http://localhost:8080

# 4. 通过 Nginx 测试
curl http://localhost

# 5. 外部测试
curl http://your-domain.com
```

### 8.8 一键部署脚本

```bash
#!/bin/bash
# deploy.sh - 一键部署脚本

set -e  # 遇到错误立即退出

# 配置
SERVER="root@your-server-ip"
APP_NAME="myapp"
REMOTE_DIR="/home/$APP_NAME"

echo "========== 开始部署 =========="

# 1. 编译
echo "[1/5] 编译中..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $APP_NAME main.go

# 2. 上传
echo "[2/5] 上传中..."
rsync -avz $APP_NAME $SERVER:$REMOTE_DIR/bin/
rsync -avz config.yaml $SERVER:$REMOTE_DIR/config/

# 3. 设置权限
echo "[3/5] 设置权限..."
ssh $SERVER "chown -R $APP_NAME:$APP_NAME $REMOTE_DIR && chmod +x $REMOTE_DIR/bin/$APP_NAME"

# 4. 重启服务
echo "[4/5] 重启服务..."
ssh $SERVER "systemctl restart $APP_NAME"

# 5. 检查状态
echo "[5/5] 检查状态..."
ssh $SERVER "systemctl status $APP_NAME --no-pager"

echo "========== 部署完成 =========="
echo "访问地址: http://your-domain.com"
```

```bash
# 使用
chmod +x deploy.sh
./deploy.sh
```

---

## 9. 常见问题排查 ⭐⭐⭐

### 9.1 服务启动失败

```bash
# 1. 查看状态
sudo systemctl status myapp

# 2. 查看详细日志
sudo journalctl -u myapp -n 100

# 常见原因：
# - 权限不足：检查文件权限和用户
# - 端口被占用：ss -tlnp | grep 8080
# - 配置文件错误：手动运行程序看报错
# - 依赖服务没启动：检查数据库、Redis 等
```

### 9.2 无法访问

```bash
# 1. 检查进程
ps aux | grep myapp

# 2. 检查端口
ss -tlnp | grep 8080

# 3. 本地测试
curl http://localhost:8080

# 4. 检查防火墙
sudo ufw status

# 5. 检查 Nginx
sudo nginx -t
sudo systemctl status nginx

# 6. 检查云服务器安全组
# 在云控制台检查是否开放了端口
```

### 9.3 内存/CPU 过高

```bash
# 1. 查看资源使用
top
htop

# 2. 查看进程详情
ps aux | grep myapp

# 3. 查看 Go 程序的 pprof
# 在代码中添加：
# import _ "net/http/pprof"
# go func() { http.ListenAndServe(":6060", nil) }()

# 访问 http://localhost:6060/debug/pprof/

# 4. 重启服务（临时解决）
sudo systemctl restart myapp
```

### 9.4 磁盘满了

```bash
# 1. 查看磁盘使用
df -h

# 2. 找出大文件
du -sh /* | sort -rh | head -10

# 3. 清理日志
sudo journalctl --vacuum-size=500M
sudo find /var/log -name "*.log" -mtime +30 -delete

# 4. 清理 Docker（如果用）
docker system prune -a
```

---

## 10. 本章总结

### 必须掌握的命令

| 类别 | 命令 | 用途 |
|------|------|------|
| 用户权限 | `chmod`, `chown`, `sudo` | 权限管理 |
| 进程管理 | `ps aux`, `top`, `kill`, `nohup` | 查看和管理进程 |
| 服务管理 | `systemctl start/stop/status/enable` | 管理 systemd 服务 |
| 日志查看 | `journalctl -u xxx -f` | 查看服务日志 |
| 网络 | `ss -tlnp`, `curl`, `ping` | 网络诊断 |
| 系统资源 | `free -h`, `df -h`, `du -sh` | 资源监控 |
| 环境变量 | `export`, `echo $VAR` | 配置管理 |
| 定时任务 | `crontab -e`, `crontab -l` | 定时执行 |

### 部署 Go 应用的标准流程

```
1. 创建专用用户（不用 root）
2. 上传编译好的程序
3. 创建 systemd 服务
4. 配置 Nginx 反向代理
5. 启动服务并设置开机自启
6. 验证部署
```

### 记忆口诀

```
权限 755 能执行，644 只读写
进程 ps aux 看状态，kill -9 强制杀
服务 systemctl 来管理，journalctl 看日志
端口 ss -tlnp 查监听，curl 测试 HTTP
```

---

## 📚 下一章预告

下一章我们将学习 **进程基础**，深入理解：
- 什么是进程？进程和程序的区别
- 进程的状态和生命周期
- PCB（进程控制块）
- 进程的创建和终止
- Go 语言中的进程操作

这些是理解操作系统的基础，也是面试常考内容！
