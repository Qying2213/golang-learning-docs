# Linux 基础命令

> **重要程度：⭐⭐⭐ 必须掌握**  
> 这是日常工作必用的命令，必须熟练！不会 Linux 命令，后端开发寸步难行！

## 📚 本章学习目标

学完本章，你将能够：
- 理解 Linux 是什么，为什么要学
- 熟练使用文件和目录操作命令
- 掌握 Vim 编辑器的基本使用
- 能够在服务器上传输文件
- 学会搜索文件和文件内容

---

## 1. 为什么要学 Linux？⭐⭐⭐

### 1.1 后端开发离不开 Linux

```
你写的代码在哪里运行？

本地开发（你的电脑：Mac/Windows）
    │
    │ 写代码、测试
    ↓
代码仓库（GitHub/GitLab）
    │
    │ 部署
    ↓
Linux 服务器（阿里云/腾讯云/AWS）
    │
    │ 运行
    ↓
用户访问你的服务

重点：99% 的服务器都是 Linux 系统！
```

### 1.2 工作中必须会的操作

| 场景 | 需要的技能 |
|------|-----------|
| 部署代码 | 上传文件、编译、运行 |
| 查看日志 | 查看文件、搜索内容、实时监控 |
| 排查问题 | 查看进程、查看端口、查看资源 |
| 配置服务 | 编辑配置文件、重启服务 |
| 定时任务 | 设置 crontab |

### 1.3 学习建议

```
🎯 学习方法：

1. 不要死记硬背，要动手练习！
2. 准备一个 Linux 环境：
   - Mac 用户：直接用终端（Terminal）
   - Windows 用户：安装 WSL 或虚拟机
   - 或者：买一台便宜的云服务器练习

3. 每个命令都要亲手敲一遍！
```

### 📝 练习题 1.1

**问题**：为什么大多数服务器使用 Linux 而不是 Windows？

<details>
<summary>点击查看答案</summary>

**答案**：

1. **免费开源**：Linux 不需要付费，Windows Server 需要购买许可证
2. **稳定性高**：Linux 服务器可以运行数年不重启
3. **安全性好**：开源代码，漏洞发现和修复更快
4. **资源占用少**：没有图形界面，更省内存和 CPU
5. **命令行强大**：自动化运维更方便
6. **社区支持**：遇到问题容易找到解决方案

</details>

---

## 2. 连接服务器 ⭐⭐⭐

### 2.1 什么是 SSH？

```
SSH = Secure Shell（安全外壳协议）

作用：让你从自己的电脑远程登录到服务器

┌─────────────────┐                    ┌─────────────────┐
│   你的电脑       │ ───── SSH ─────→ │   Linux 服务器   │
│   (客户端)       │     加密连接      │   (服务端)       │
└─────────────────┘                    └─────────────────┘

为什么用 SSH？
- 加密传输，安全
- 可以远程执行命令
- 可以传输文件
```

### 2.2 SSH 登录命令

```bash
# 最基本的登录命令
ssh 用户名@服务器IP

# 例如：用 root 用户登录 IP 为 192.168.1.100 的服务器
ssh root@192.168.1.100

# 第一次连接会问你是否信任，输入 yes
# 然后输入密码（输入时看不到字符，这是正常的）
```

**实际操作演示**：
```bash
$ ssh root@192.168.1.100
The authenticity of host '192.168.1.100' can't be established.
Are you sure you want to continue connecting (yes/no)? yes   # 输入 yes
root@192.168.1.100's password:                               # 输入密码
Welcome to Ubuntu 20.04 LTS
root@server:~#                                               # 登录成功！
```

### 2.3 指定端口登录

```bash
# 默认 SSH 端口是 22
# 如果服务器改了端口（比如改成 2222），需要指定

ssh -p 2222 root@192.168.1.100
#    └─ -p 参数指定端口
```

### 2.4 使用密钥登录（更安全、更方便）

```
为什么用密钥？
- 不用每次输入密码
- 比密码更安全
- 可以禁用密码登录，防止暴力破解
```

**第一步：生成密钥对**
```bash
# 在你的电脑上执行（不是服务器！）
ssh-keygen -t rsa -b 4096

# 会问你几个问题，直接按回车用默认值就行
# Enter file in which to save the key: 按回车
# Enter passphrase: 按回车（或设置密码）
# Enter same passphrase again: 按回车

# 生成两个文件：
# ~/.ssh/id_rsa      私钥（绝对不能泄露！）
# ~/.ssh/id_rsa.pub  公钥（可以给别人）
```

**第二步：把公钥复制到服务器**
```bash
# 方法1：使用 ssh-copy-id（推荐）
ssh-copy-id root@192.168.1.100

# 方法2：手动复制
# 先查看公钥内容
cat ~/.ssh/id_rsa.pub

# 登录服务器，把公钥内容追加到这个文件
echo "公钥内容" >> ~/.ssh/authorized_keys
```

**第三步：测试免密登录**
```bash
ssh root@192.168.1.100
# 不用输入密码就能登录了！
```

### 📝 练习题 2.1

**问题**：私钥和公钥有什么区别？为什么私钥不能泄露？

<details>
<summary>点击查看答案</summary>

**答案**：

**公钥（id_rsa.pub）**：
- 可以公开，放到服务器上
- 用于加密数据
- 就像一把锁，可以给任何人

**私钥（id_rsa）**：
- 必须保密，只能自己持有
- 用于解密数据
- 就像钥匙，只有你能开锁

**为什么私钥不能泄露？**
- 有了私钥，别人就能登录你的所有服务器
- 相当于把家里钥匙给了陌生人
- 如果泄露了，必须立即重新生成密钥对

</details>

---

## 3. 文件和目录操作 ⭐⭐⭐

### 3.1 Linux 目录结构

```
在学命令之前，先了解 Linux 的目录结构：

/                    # 根目录，所有文件的起点
├── home/            # 普通用户的家目录
│   ├── user1/       # user1 的家目录
│   └── user2/       # user2 的家目录
├── root/            # root 用户的家目录
├── etc/             # 配置文件目录
├── var/             # 可变数据（日志、数据库等）
│   └── log/         # 日志文件
├── tmp/             # 临时文件
├── usr/             # 用户程序
│   ├── bin/         # 用户命令
│   └── local/       # 本地安装的程序
└── opt/             # 第三方软件

常用目录：
- ~     表示当前用户的家目录（/home/用户名 或 /root）
- .     表示当前目录
- ..    表示上级目录
```

### 3.2 pwd - 显示当前目录

```bash
# pwd = Print Working Directory（打印工作目录）

$ pwd
/home/user

# 你现在在哪个目录？用 pwd 看看
```

### 3.3 ls - 列出文件和目录

```bash
# ls = List（列出）

# 最简单的用法：列出当前目录的文件
$ ls
main.go  config  README.md

# 列出指定目录
$ ls /home
user1  user2

# 常用参数组合：ls -la（必须记住！）
$ ls -la
total 16
drwxr-xr-x  3 root root 4096 Jan  7 10:00 .
drwxr-xr-x 18 root root 4096 Jan  7 09:00 ..
-rw-r--r--  1 root root 1234 Jan  7 10:00 main.go
drwxr-xr-x  2 root root 4096 Jan  7 10:00 config
-rw-r--r--  1 root root  567 Jan  7 10:00 README.md
```

**参数详解**：
```bash
ls -l    # 详细信息（权限、大小、时间）
ls -a    # 显示隐藏文件（以 . 开头的文件）
ls -h    # 人性化显示大小（KB、MB、GB）
ls -la   # 组合使用：详细 + 隐藏
ls -lh   # 组合使用：详细 + 人性化大小
ls -lah  # 全部组合
```

**解读 ls -l 的输出**：
```
-rw-r--r--  1  root  root  1234  Jan 7 10:00  main.go
│├─┬─┤├─┬─┤     │     │     │        │          │
││  │  │        │     │     │        │          └─ 文件名
││  │  │        │     │     │        └─ 修改时间
││  │  │        │     │     └─ 文件大小（字节）
││  │  │        │     └─ 所属组
││  │  │        └─ 所有者
││  │  └─ 其他用户权限（r=读，-=无）
││  └─ 组用户权限（r=读，-=无）
│└─ 所有者权限（r=读，w=写，-=无执行）
└─ 文件类型（-=普通文件，d=目录，l=链接）

例子解读：
-rw-r--r--  →  普通文件，所有者可读写，其他人只读
drwxr-xr-x  →  目录，所有者可读写执行，其他人可读和执行
```

### 📝 练习题 3.1

**问题**：`ls -la` 输出中，`.` 和 `..` 分别代表什么？

<details>
<summary>点击查看答案</summary>

**答案**：
- `.`（一个点）：代表当前目录
- `..`（两个点）：代表上级目录（父目录）

**用途**：
```bash
# 复制文件到当前目录
cp /home/file.txt .

# 切换到上级目录
cd ..

# 运行当前目录的程序
./myapp
```

</details>

### 3.4 cd - 切换目录

```bash
# cd = Change Directory（切换目录）

# 切换到指定目录
cd /home/user

# 切换到家目录（三种方式）
cd ~
cd $HOME
cd          # 不带参数也是回家目录

# 切换到上级目录
cd ..

# 切换到上上级目录
cd ../..

# 切换到上一次的目录（很有用！）
cd -

# 实际例子：
$ pwd
/home/user/project
$ cd /var/log
$ pwd
/var/log
$ cd -          # 回到上一次的目录
/home/user/project
```

### 3.5 mkdir - 创建目录

```bash
# mkdir = Make Directory（创建目录）

# 创建单个目录
mkdir myproject

# 创建多级目录（必须加 -p）
mkdir -p myproject/src/main
#      └─ -p 表示创建父目录（如果不存在）

# 不加 -p 会报错：
$ mkdir myproject/src/main
mkdir: cannot create directory 'myproject/src/main': No such file or directory

# 一次创建多个目录
mkdir dir1 dir2 dir3
```

### 3.6 touch - 创建空文件

```bash
# touch 的作用：
# 1. 如果文件不存在，创建空文件
# 2. 如果文件存在，更新修改时间

# 创建空文件
touch main.go
touch config.yaml

# 一次创建多个文件
touch file1.txt file2.txt file3.txt
```

### 3.7 cp - 复制文件和目录

```bash
# cp = Copy（复制）

# 复制文件
cp main.go main_backup.go
cp main.go /home/backup/

# 复制目录（必须加 -r）
cp -r myproject myproject_backup
#  └─ -r 表示递归复制（包括子目录）

# 复制时保留属性（权限、时间等）
cp -rp myproject myproject_backup
#   └─ -p 表示保留属性

# 复制并显示进度
cp -rv myproject myproject_backup
#   └─ -v 表示显示详细信息
```

### 3.8 mv - 移动/重命名

```bash
# mv = Move（移动）

# 重命名文件
mv old_name.go new_name.go

# 移动文件到另一个目录
mv main.go /home/backup/

# 移动并重命名
mv main.go /home/backup/main_old.go

# 移动目录（不需要 -r）
mv myproject /home/backup/

# 批量移动
mv *.go /home/backup/
```

### 3.9 rm - 删除文件和目录 ⚠️

```bash
# rm = Remove（删除）
# ⚠️ Linux 没有回收站，删除后无法恢复！

# 删除文件
rm file.txt

# 删除前确认
rm -i file.txt
#  └─ -i 表示交互式，会问你是否确认

# 删除目录（必须加 -r）
rm -r myproject
#  └─ -r 表示递归删除

# 强制删除（不提示）
rm -f file.txt
#  └─ -f 表示强制，不提示

# 强制递归删除（最危险的命令！）
rm -rf myproject

# ⚠️⚠️⚠️ 绝对不要执行这个命令！！！
# rm -rf /          # 删除整个系统
# rm -rf /*         # 删除整个系统
# rm -rf ~          # 删除家目录所有文件
```

### 📝 练习题 3.2

**问题**：执行以下命令，创建一个项目目录结构：

```
myapp/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/
│   └── service/
├── config/
│   └── config.yaml
└── README.md
```

<details>
<summary>点击查看答案</summary>

```bash
# 创建目录结构
mkdir -p myapp/cmd
mkdir -p myapp/internal/handler
mkdir -p myapp/internal/service
mkdir -p myapp/config

# 创建文件
touch myapp/cmd/main.go
touch myapp/config/config.yaml
touch myapp/README.md

# 验证
ls -la myapp/
ls -la myapp/internal/
```

</details>

---

## 4. 查看文件内容 ⭐⭐⭐

### 4.1 cat - 查看整个文件

```bash
# cat = Concatenate（连接）
# 用于查看文件内容

# 查看文件
cat main.go

# 显示行号
cat -n main.go

# 查看多个文件
cat file1.txt file2.txt

# 合并文件
cat file1.txt file2.txt > combined.txt
```

**cat 的问题**：文件太大时，会刷屏，看不到开头的内容

### 4.2 less - 分页查看（推荐）

```bash
# less 可以分页查看大文件

less main.go

# 在 less 中的操作：
# 空格键    向下翻一页
# b        向上翻一页
# j 或 ↓   向下一行
# k 或 ↑   向上一行
# g        跳到文件开头
# G        跳到文件结尾
# /关键词   搜索（按 n 下一个，N 上一个）
# q        退出

# 例子：查看日志文件
less /var/log/syslog
```

### 4.3 head - 查看文件开头

```bash
# head 查看文件的前几行

# 默认显示前 10 行
head main.go

# 指定行数
head -n 20 main.go    # 前 20 行
head -20 main.go      # 简写

# 查看多个文件
head -n 5 file1.txt file2.txt
```

### 4.4 tail - 查看文件结尾（重要！）

```bash
# tail 查看文件的最后几行

# 默认显示最后 10 行
tail main.go

# 指定行数
tail -n 20 main.go    # 最后 20 行
tail -20 main.go      # 简写

# ⭐⭐⭐ 实时监控文件（看日志必用！）
tail -f /var/log/app.log
#    └─ -f 表示 follow，实时显示新增内容

# 按 Ctrl+C 退出实时监控

# 实时监控并显示最后 100 行
tail -n 100 -f /var/log/app.log
```

**工作中的场景**：
```bash
# 场景1：查看应用日志
tail -f /var/log/myapp/app.log

# 场景2：查看 nginx 访问日志
tail -f /var/log/nginx/access.log

# 场景3：查看系统日志
tail -f /var/log/syslog
```

### 4.5 wc - 统计行数/字数

```bash
# wc = Word Count（字数统计）

# 统计行数
wc -l main.go
# 输出：50 main.go（表示 50 行）

# 统计字数
wc -w main.go

# 统计字节数
wc -c main.go

# 全部统计
wc main.go
# 输出：50  200  1234 main.go
#       行数 字数 字节数
```

### 📝 练习题 4.1

**问题**：如何统计一个日志文件中 "ERROR" 出现了多少次？

<details>
<summary>点击查看答案</summary>

```bash
# 方法1：grep + wc
grep "ERROR" app.log | wc -l

# 方法2：grep -c
grep -c "ERROR" app.log

# 解释：
# grep "ERROR" app.log  → 找出包含 ERROR 的行
# | wc -l               → 统计行数
# grep -c               → 直接统计匹配的行数
```

</details>

---

## 5. Vim 编辑器 ⭐⭐⭐

### 5.1 为什么要学 Vim？

```
在服务器上，你没有 VS Code！
只能用命令行编辑器，Vim 是最常用的。

常见场景：
- 修改配置文件
- 快速编辑代码
- 查看和编辑日志
```

### 5.2 Vim 的三种模式

```
这是 Vim 最重要的概念！

┌─────────────────────────────────────────────────────┐
│                                                      │
│    ┌──────────────┐                                 │
│    │   普通模式   │  ← 默认进入这个模式             │
│    │  (Normal)    │  用于：移动光标、删除、复制     │
│    └──────┬───────┘                                 │
│           │                                          │
│     按 i  │  按 Esc                                 │
│           ↓                                          │
│    ┌──────────────┐                                 │
│    │   插入模式   │  ← 在这个模式下才能输入文字     │
│    │  (Insert)    │                                 │
│    └──────────────┘                                 │
│                                                      │
│    在普通模式下按 : 进入命令模式                    │
│    ┌──────────────┐                                 │
│    │   命令模式   │  ← 用于：保存、退出、搜索替换   │
│    │  (Command)   │                                 │
│    └──────────────┘                                 │
│                                                      │
└─────────────────────────────────────────────────────┘

记住：
- 刚打开文件 → 普通模式
- 想输入文字 → 按 i 进入插入模式
- 输入完了 → 按 Esc 回到普通模式
- 想保存退出 → 按 : 进入命令模式
```

### 5.3 打开和退出 Vim

```bash
# 打开文件
vim main.go

# 打开文件并跳到第 10 行
vim +10 main.go

# 打开文件并搜索关键词
vim +/error main.go
```

**退出 Vim（最常问的问题！）**：
```bash
# 在普通模式下（按 Esc 确保在普通模式）

:q          # 退出（如果没修改）
:q!         # 强制退出（不保存修改）
:w          # 保存（不退出）
:wq         # 保存并退出
:wq!        # 强制保存并退出
ZZ          # 保存并退出（快捷键）
```

### 5.4 插入文字

```bash
# 在普通模式下，按以下键进入插入模式：

i     # 在光标前插入
a     # 在光标后插入
I     # 在行首插入
A     # 在行尾插入
o     # 在下一行插入新行
O     # 在上一行插入新行

# 按 Esc 返回普通模式
```

### 5.5 移动光标

```bash
# 在普通模式下：

# 基本移动
h     # 左
j     # 下
k     # 上
l     # 右

# 快速移动
w     # 下一个单词开头
b     # 上一个单词开头
0     # 行首
$     # 行尾
gg    # 文件开头
G     # 文件结尾
:10   # 跳到第 10 行

# 翻页
Ctrl+f    # 向下翻页
Ctrl+b    # 向上翻页
```

### 5.6 删除和复制

```bash
# 在普通模式下：

# 删除
x         # 删除当前字符
dd        # 删除当前行
5dd       # 删除 5 行
dw        # 删除一个单词
d$        # 删除到行尾
d0        # 删除到行首

# 复制
yy        # 复制当前行
5yy       # 复制 5 行
yw        # 复制一个单词

# 粘贴
p         # 在光标后粘贴
P         # 在光标前粘贴

# 撤销和重做
u         # 撤销
Ctrl+r    # 重做
```

### 5.7 搜索和替换

```bash
# 在普通模式下：

# 搜索
/keyword      # 向下搜索 keyword
?keyword      # 向上搜索 keyword
n             # 下一个匹配
N             # 上一个匹配

# 替换（命令模式）
:s/old/new/           # 替换当前行第一个
:s/old/new/g          # 替换当前行所有
:%s/old/new/g         # 替换全文所有
:%s/old/new/gc        # 替换全文，每次确认
```

### 5.8 实用设置

```bash
# 在命令模式下：

:set number       # 显示行号
:set nonumber     # 隐藏行号
:syntax on        # 语法高亮
:set tabstop=4    # Tab 宽度为 4
:set expandtab    # Tab 转空格
```

### 📝 练习题 5.1

**问题**：用 Vim 完成以下操作：
1. 创建文件 test.txt
2. 输入 "Hello World"
3. 复制这一行
4. 粘贴 3 次
5. 保存退出

<details>
<summary>点击查看答案</summary>

```bash
# 1. 打开/创建文件
vim test.txt

# 2. 输入文字
i                    # 进入插入模式
Hello World          # 输入文字
Esc                  # 回到普通模式

# 3. 复制当前行
yy

# 4. 粘贴 3 次
p
p
p

# 5. 保存退出
:wq
```

**结果**：文件内容是 4 行 "Hello World"

</details>

### 📝 练习题 5.2

**问题**：如何在 Vim 中把所有的 "foo" 替换成 "bar"？

<details>
<summary>点击查看答案</summary>

```bash
# 在普通模式下输入：
:%s/foo/bar/g

# 解释：
# %     表示全文
# s     表示替换
# foo   要被替换的内容
# bar   替换成的内容
# g     表示全局（每行所有匹配）

# 如果想每次确认：
:%s/foo/bar/gc
# c 表示 confirm（确认）
```

</details>

---

## 6. 文件传输 ⭐⭐⭐

### 6.1 为什么要学文件传输？

```
工作中的场景：

你在本地电脑写好了代码
    │
    │ 需要传到服务器上运行
    ↓
┌─────────────────┐                    ┌─────────────────┐
│   你的电脑       │ ───── ？ ─────→  │   Linux 服务器   │
│   (本地)        │     怎么传？      │   (远程)        │
└─────────────────┘                    └─────────────────┘

答案：用 scp 或 rsync 命令！
```

### 6.2 scp - 安全复制（详解）

#### 什么是 scp？
```
scp = Secure Copy（安全复制）

特点：
- 通过 SSH 加密传输，安全
- 用法简单，和 cp 命令类似
- 可以上传，也可以下载
```

#### scp 命令格式
```
理解 scp 的关键是搞清楚"从哪里"到"哪里"：

┌─────────────────────────────────────────────────────────┐
│  scp  [选项]  源文件  目标位置                            │
│              └──┬──┘  └──┬──┘                           │
│                 │        │                              │
│           从哪里复制   复制到哪里                         │
└─────────────────────────────────────────────────────────┘

本地路径：直接写路径，如 ./main.go 或 /home/user/file.txt
远程路径：用户名@服务器IP:路径，如 root@192.168.1.100:/home/
```

#### 场景1：上传文件到服务器

```bash
# 命令格式
scp 本地文件 用户名@服务器IP:远程目录

# 实际例子：把本地的 main.go 上传到服务器的 /home 目录
scp main.go root@192.168.1.100:/home/

# 分解理解：
# scp                    → 命令
# main.go                → 本地文件（当前目录下的 main.go）
# root                   → 服务器的用户名
# @                      → 分隔符
# 192.168.1.100          → 服务器 IP 地址
# :                      → 分隔符（很重要！不能少！）
# /home/                 → 服务器上的目标目录
```

**执行过程演示**：
```bash
$ scp main.go root@192.168.1.100:/home/
root@192.168.1.100's password:        # 输入服务器密码
main.go                 100%  1234     1.2KB/s   00:00
#                       │     │        │         │
#                       │     │        │         └─ 用时
#                       │     │        └─ 传输速度
#                       │     └─ 文件大小
#                       └─ 进度 100%

# 传输完成！文件已经在服务器的 /home/main.go
```

#### 场景2：上传文件并重命名

```bash
# 上传时可以指定新的文件名
scp main.go root@192.168.1.100:/home/app.go
#                                    └─ 上传后叫 app.go

# 对比：
scp main.go root@192.168.1.100:/home/      # 保持原名 main.go
scp main.go root@192.168.1.100:/home/app.go  # 改名为 app.go
```

#### 场景3：上传整个目录

```bash
# 上传目录必须加 -r 参数！
# -r = recursive（递归）

scp -r myproject root@192.168.1.100:/home/
#   └─ 必须加 -r，否则报错！

# 例子：上传整个项目目录
$ ls myproject/
main.go  config/  README.md

$ scp -r myproject root@192.168.1.100:/home/
root@192.168.1.100's password:
main.go                 100%  1234     1.2KB/s   00:00
config/app.yaml         100%   567     0.6KB/s   00:00
README.md               100%   890     0.9KB/s   00:00

# 服务器上就有了 /home/myproject/ 目录
```

#### 场景4：从服务器下载文件

```bash
# 下载就是把源和目标反过来！

# 格式
scp 用户名@服务器IP:远程文件 本地目录

# 例子：下载服务器上的日志文件到当前目录
scp root@192.168.1.100:/var/log/app.log ./
#                       └─ 服务器上的文件    └─ 本地当前目录

# 下载并重命名
scp root@192.168.1.100:/var/log/app.log ./server.log

# 下载整个目录（加 -r）
scp -r root@192.168.1.100:/home/myproject ./
```

#### 场景5：指定端口

```bash
# 默认 SSH 端口是 22
# 如果服务器改了端口，需要用 -P 指定（注意是大写 P！）

scp -P 2222 main.go root@192.168.1.100:/home/
#   └─ 大写 P！

# 为什么是大写 P？
# 因为小写 -p 是保留文件属性的意思
scp -p main.go root@192.168.1.100:/home/  # 保留修改时间等属性
scp -P 2222 main.go root@192.168.1.100:/home/  # 指定端口
```

#### scp 常用参数总结

```bash
-r    递归复制目录（上传/下载目录必须加）
-P    指定端口（大写 P）
-p    保留文件属性（小写 p）
-v    显示详细过程（调试用）
-C    压缩传输（大文件时有用）

# 组合使用
scp -rp myproject root@192.168.1.100:/home/  # 递归 + 保留属性
scp -rv myproject root@192.168.1.100:/home/  # 递归 + 显示详情
```

### 📝 练习题 6.1

**问题**：写出以下操作的 scp 命令：
1. 把本地的 `config.yaml` 上传到服务器 `192.168.1.100` 的 `/etc/myapp/` 目录
2. 把服务器上的 `/var/log/nginx/access.log` 下载到本地当前目录
3. 把本地的 `dist` 目录上传到服务器的 `/var/www/` 目录

<details>
<summary>点击查看答案</summary>

```bash
# 1. 上传配置文件
scp config.yaml root@192.168.1.100:/etc/myapp/

# 2. 下载日志文件
scp root@192.168.1.100:/var/log/nginx/access.log ./

# 3. 上传目录（必须加 -r）
scp -r dist root@192.168.1.100:/var/www/
```

</details>

### 📝 练习题 6.2

**问题**：执行 `scp myproject root@192.168.1.100:/home/` 报错了，可能是什么原因？

<details>
<summary>点击查看答案</summary>

**最可能的原因**：`myproject` 是一个目录，但没有加 `-r` 参数！

```bash
# 错误
scp myproject root@192.168.1.100:/home/
# 报错：myproject: not a regular file

# 正确
scp -r myproject root@192.168.1.100:/home/
```

**其他可能原因**：
1. 服务器 IP 或用户名错误
2. 密码错误
3. 服务器 SSH 服务没启动
4. 防火墙阻止了连接
5. 远程目录不存在或没有写权限

</details>

---

### 6.3 rsync - 同步工具（更强大）

#### 为什么要用 rsync？

```
scp 的问题：

假设你有一个 100MB 的项目目录
第一次上传：传了 100MB ✓
改了一个文件（1KB）
第二次上传：又传了 100MB ✗ （太慢了！）

rsync 的优势：

第一次上传：传了 100MB ✓
改了一个文件（1KB）
第二次上传：只传 1KB ✓ （只传变化的部分！）

这就是"增量传输"！
```

#### rsync 基本用法

```bash
# 格式
rsync [选项] 源 目标

# 最常用的命令（记住这个就够了！）
rsync -avz 本地目录/ 用户名@服务器:远程目录/

# 参数说明：
# -a = archive（归档模式）
#      包含：-r 递归、-l 保留链接、-p 保留权限、-t 保留时间、-g 保留组、-o 保留所有者
# -v = verbose（显示详细信息）
# -z = compress（压缩传输）
```

#### 重要：目录路径的斜杠

```bash
# 这是 rsync 最容易搞错的地方！

# 源目录末尾有斜杠：复制目录里的内容
rsync -avz myproject/ root@server:/home/myproject/
# 结果：/home/myproject/main.go

# 源目录末尾没有斜杠：复制整个目录
rsync -avz myproject root@server:/home/
# 结果：/home/myproject/main.go

# 看起来一样？再看这个例子：
rsync -avz myproject/ root@server:/home/
# 结果：/home/main.go  ← 文件直接放到 /home 了！

# 建议：源和目标都加斜杠，最清晰
rsync -avz ./myproject/ root@server:/home/myproject/
```

#### 实际使用示例

```bash
# 1. 同步项目到服务器
rsync -avz ./myproject/ root@192.168.1.100:/home/myproject/

# 输出示例：
sending incremental file list
main.go
config/app.yaml
sent 2,345 bytes  received 123 bytes  1,645.33 bytes/sec
total size is 10,234  speedup is 4.14

# 2. 第二次同步（只改了 main.go）
rsync -avz ./myproject/ root@192.168.1.100:/home/myproject/

# 输出示例：
sending incremental file list
main.go                    # 只传了这一个文件！
sent 1,234 bytes  received 35 bytes  2,538.00 bytes/sec
total size is 10,234  speedup is 8.06
```

#### 删除远程多余的文件

```bash
# 场景：本地删除了一个文件，希望服务器也删除

# 不加 --delete：服务器上的文件不会被删除
rsync -avz ./myproject/ root@server:/home/myproject/

# 加 --delete：保持完全同步（服务器上多余的文件会被删除）
rsync -avz --delete ./myproject/ root@server:/home/myproject/

# ⚠️ 注意：--delete 要小心使用，可能会删除重要文件！
# 建议先用 --dry-run 预览
rsync -avz --delete --dry-run ./myproject/ root@server:/home/myproject/
#                   └─ 只显示会做什么，不真正执行
```

#### 排除某些文件

```bash
# 不想上传某些文件（如日志、依赖包）

rsync -avz --exclude='*.log' ./myproject/ root@server:/home/myproject/
#          └─ 排除所有 .log 文件

# 排除多个
rsync -avz \
  --exclude='*.log' \
  --exclude='node_modules' \
  --exclude='.git' \
  --exclude='vendor' \
  ./myproject/ root@server:/home/myproject/

# 使用排除文件
echo "*.log" > .rsync-exclude
echo "node_modules" >> .rsync-exclude
echo ".git" >> .rsync-exclude

rsync -avz --exclude-from='.rsync-exclude' ./myproject/ root@server:/home/myproject/
```

#### rsync 常用参数

```bash
-a    归档模式（最常用，包含 -rlptgoD）
-v    显示详细信息
-z    压缩传输
-P    显示进度 + 支持断点续传（等于 --progress --partial）
-n    --dry-run，只预览不执行
--delete    删除目标中源没有的文件
--exclude   排除文件
-e    指定 SSH 命令（如指定端口）

# 指定端口
rsync -avz -e 'ssh -p 2222' ./myproject/ root@server:/home/myproject/
```

### 📝 练习题 6.3

**问题**：scp 和 rsync 有什么区别？什么时候用哪个？

<details>
<summary>点击查看答案</summary>

| 特性 | scp | rsync |
|------|-----|-------|
| 传输方式 | 全量传输（每次都传所有文件） | 增量传输（只传变化的文件） |
| 速度 | 较慢 | 较快 |
| 压缩 | 需要手动压缩 | 内置压缩（-z） |
| 删除同步 | 不支持 | 支持（--delete） |
| 排除文件 | 不支持 | 支持（--exclude） |
| 断点续传 | 不支持 | 支持（-P） |
| 学习成本 | 简单 | 稍复杂 |

**使用建议**：
- **用 scp**：传输单个小文件、临时传一下
- **用 rsync**：同步项目目录、定期备份、大文件传输

**工作中的选择**：
```bash
# 临时传个文件 → scp
scp config.yaml root@server:/home/

# 部署项目 → rsync
rsync -avz --exclude='.git' ./myproject/ root@server:/home/myproject/

# 备份数据 → rsync
rsync -avz --delete /data/ root@backup-server:/backup/data/
```

</details>

### 📝 练习题 6.4

**问题**：写一个 rsync 命令，把本地的 `myapp` 目录同步到服务器，要求：
1. 排除 `.git` 目录
2. 排除所有 `.log` 文件
3. 显示传输进度
4. 服务器 SSH 端口是 2222

<details>
<summary>点击查看答案</summary>

```bash
rsync -avzP \
  --exclude='.git' \
  --exclude='*.log' \
  -e 'ssh -p 2222' \
  ./myapp/ root@192.168.1.100:/home/myapp/

# 参数解释：
# -a    归档模式
# -v    显示详细信息
# -z    压缩传输
# -P    显示进度（等于 --progress --partial）
# --exclude='.git'    排除 .git 目录
# --exclude='*.log'   排除所有 .log 文件
# -e 'ssh -p 2222'    指定 SSH 端口为 2222
```

</details>

---

## 7. 搜索文件和内容 ⭐⭐⭐

### 7.1 find - 查找文件（详解）

#### 为什么要用 find？

```
场景：你知道有个文件叫 main.go，但忘了放在哪个目录

方法1：一个一个目录去找 → 太慢了！
方法2：用 find 命令 → 几秒钟找到！

find 就是用来在目录中搜索文件的命令
```

#### find 命令的结构

```
find 命令和其他命令不太一样，它需要"条件"来告诉它找什么

┌─────────────────────────────────────────────────────────┐
│  find  搜索路径  条件1  条件2  ...                        │
│        └──┬──┘  └─────┬─────┘                           │
│           │           │                                 │
│       在哪里找     找什么样的文件                          │
└─────────────────────────────────────────────────────────┘

为什么要用 -name？
因为 find 可以按很多条件查找：
- 按名字找 → -name
- 按类型找 → -type
- 按大小找 → -size
- 按时间找 → -mtime

所以必须告诉 find 你要按什么条件找！
```

#### 按名称查找 -name

```bash
# -name 表示"按文件名查找"

# 格式
find 搜索路径 -name "文件名"

# 例子：在 /home 目录下找 main.go
find /home -name "main.go"

# 分解理解：
# find       → 命令
# /home      → 在 /home 目录及其子目录中搜索
# -name      → 按名称查找（这个不能省！）
# "main.go"  → 要找的文件名

# 如果不写 -name 会怎样？
find /home "main.go"
# 报错！find 不知道 "main.go" 是什么意思
```

#### 使用通配符

```bash
# * 表示任意字符（零个或多个）

# 找所有 .go 文件
find /home -name "*.go"
#                 └─ * 匹配任意字符，所以匹配 main.go, test.go, xxx.go 等

# 找所有以 test 开头的文件
find /home -name "test*"
# 匹配：test.go, test1.txt, testing.md 等

# 找所有包含 config 的文件
find /home -name "*config*"
# 匹配：config.yaml, app_config.json, myconfig 等

# ⚠️ 重要：通配符要用引号包起来！
find /home -name *.go      # ❌ 可能出错
find /home -name "*.go"    # ✅ 正确
```

#### 忽略大小写 -iname

```bash
# -iname = ignore case name（忽略大小写的名称）

# 问题：Linux 文件名区分大小写
find /home -name "main.go"     # 只能找到 main.go
                               # 找不到 Main.go, MAIN.GO

# 解决：用 -iname
find /home -iname "main.go"    # 能找到 main.go, Main.go, MAIN.GO, MaIn.Go 等
```

#### 按类型查找 -type

```bash
# -type 表示"按文件类型查找"

# 常用类型：
# f = file（普通文件）
# d = directory（目录）
# l = link（符号链接）

# 只找文件，不要目录
find /home -name "*.go" -type f

# 只找目录
find /home -name "config" -type d

# 为什么要指定类型？
# 因为可能有个文件叫 config，也有个目录叫 config
# 指定 -type d 就只找目录
```

#### 按大小查找 -size

```bash
# -size 表示"按文件大小查找"

# 单位：
# c = bytes（字节）
# k = KB
# M = MB
# G = GB

# + 表示大于，- 表示小于

# 找大于 100MB 的文件
find /home -size +100M

# 找小于 1KB 的文件
find /home -size -1k

# 找大于 1GB 的文件
find /var -size +1G

# 组合：找大于 10MB 的日志文件
find /var/log -name "*.log" -size +10M
```

#### 按时间查找 -mtime

```bash
# -mtime = modification time（修改时间）
# 单位是"天"

# -mtime -7  表示 7 天内修改的
# -mtime +30 表示 30 天前修改的
# -mtime 1   表示恰好 1 天前修改的

# 找 7 天内修改的文件
find /home -mtime -7

# 找 30 天前的旧文件
find /home -mtime +30

# 找 7 天内修改的日志文件
find /var/log -name "*.log" -mtime -7
```

#### 组合多个条件

```bash
# 多个条件写在一起，表示"并且"（AND）

# 找 .log 文件，并且大于 10MB
find /var/log -name "*.log" -size +10M

# 找 .go 文件，并且是 7 天内修改的
find /home -name "*.go" -mtime -7

# 找目录，并且名字是 node_modules
find /home -type d -name "node_modules"
```

#### 找到后执行命令 -exec

```bash
# -exec 表示"对找到的文件执行命令"

# 格式
find 路径 条件 -exec 命令 {} \;
#                      └┬┘ └┬┘
#                       │   └─ 结束标记（必须有！）
#                       └─ 代表找到的每个文件

# 例子：找到 .log 文件并显示详细信息
find /var/log -name "*.log" -exec ls -lh {} \;

# 例子：找到 .tmp 文件并删除
find /tmp -name "*.tmp" -exec rm {} \;

# 例子：找到大文件并显示大小
find /home -size +100M -exec ls -lh {} \;

# ⚠️ 注意：{} 和 \; 之间有空格！
find /home -name "*.log" -exec rm {} \;   # ✅ 正确
find /home -name "*.log" -exec rm {}\;    # ❌ 错误
```

#### 常用 find 命令速查

```bash
# 按名称找
find /home -name "main.go"           # 精确匹配
find /home -name "*.go"              # 通配符
find /home -iname "main.go"          # 忽略大小写

# 按类型找
find /home -type f                   # 只找文件
find /home -type d                   # 只找目录

# 按大小找
find /home -size +100M               # 大于 100MB
find /home -size -1k                 # 小于 1KB

# 按时间找
find /home -mtime -7                 # 7 天内修改的
find /home -mtime +30                # 30 天前修改的

# 组合使用
find /var/log -name "*.log" -size +10M -mtime +7
# 找 /var/log 下，.log 文件，大于 10MB，7 天前的
```

### 📝 练习题 7.1

**问题**：`find /home -name "*.go"` 这个命令中，为什么 `*.go` 要用引号包起来？

<details>
<summary>点击查看答案</summary>

**答案**：防止 shell 提前展开通配符！

```bash
# 如果不加引号：
find /home -name *.go

# shell 会先把 *.go 展开成当前目录下的 .go 文件
# 假设当前目录有 main.go 和 test.go
# 命令就变成了：
find /home -name main.go test.go
# 这就出错了！

# 加引号后：
find /home -name "*.go"
# shell 不会展开，*.go 原样传给 find
# find 自己去匹配所有 .go 文件
```

**结论**：用 find 时，文件名模式一定要加引号！

</details>

### 📝 练习题 7.2

**问题**：写一个 find 命令，找出 /var/log 目录下所有大于 50MB 的 .log 文件，并显示它们的大小。

<details>
<summary>点击查看答案</summary>

```bash
find /var/log -name "*.log" -size +50M -exec ls -lh {} \;

# 分解：
# find /var/log        → 在 /var/log 目录搜索
# -name "*.log"        → 文件名以 .log 结尾
# -size +50M           → 大于 50MB
# -exec ls -lh {} \;   → 对每个找到的文件执行 ls -lh

# 输出示例：
# -rw-r--r-- 1 root root 120M Jan 7 10:00 /var/log/syslog
# -rw-r--r-- 1 root root  80M Jan 7 09:00 /var/log/auth.log
```

</details>

---

### 7.2 grep - 搜索文件内容（超重要！）

```bash
# grep 用于在文件中搜索内容

# 基本用法
grep "关键词" 文件名

# 例子：在日志中搜索 error
grep "error" app.log

# 常用参数
grep -i "error" app.log        # 忽略大小写
grep -n "error" app.log        # 显示行号
grep -r "error" ./             # 递归搜索目录
grep -v "debug" app.log        # 反向匹配（不包含）
grep -c "error" app.log        # 统计匹配行数

# 显示上下文
grep -A 3 "error" app.log      # 显示匹配行及后 3 行
grep -B 3 "error" app.log      # 显示匹配行及前 3 行
grep -C 3 "error" app.log      # 显示匹配行及前后各 3 行

# 正则表达式
grep -E "error|warning" app.log    # 匹配 error 或 warning
grep "^2024" app.log               # 以 2024 开头的行
grep "error$" app.log              # 以 error 结尾的行
```

### 7.3 实际工作场景

```bash
# 场景1：在项目中搜索某个函数
grep -rn "func main" ./

# 场景2：查看日志中的错误
grep -i "error" /var/log/app.log | tail -20

# 场景3：统计今天的错误数量
grep "2024-01-07" app.log | grep -c "ERROR"

# 场景4：查看某个 IP 的访问记录
grep "192.168.1.100" /var/log/nginx/access.log

# 场景5：排除某些内容
grep -v "DEBUG" app.log | grep "ERROR"
```

### 📝 练习题 7.3

**问题**：如何找出 /var/log 目录下所有大于 100MB 的 .log 文件？

<details>
<summary>点击查看答案</summary>

```bash
find /var/log -name "*.log" -size +100M

# 如果还想看文件大小：
find /var/log -name "*.log" -size +100M -exec ls -lh {} \;
```

</details>

---

## 8. 管道和重定向 ⭐⭐⭐

### 8.1 管道 |

```bash
# 管道：把前一个命令的输出，作为后一个命令的输入

# 语法
命令1 | 命令2 | 命令3

# 例子1：查看进程并过滤
ps aux | grep "go"

# 例子2：统计日志中 ERROR 的数量
cat app.log | grep "ERROR" | wc -l

# 例子3：查看最占内存的进程
ps aux | sort -k4 -rn | head -10
#        └─ 按第4列（内存）排序

# 例子4：去重并排序
cat file.txt | sort | uniq

# 例子5：查看目录大小并排序
du -sh * | sort -h

```

### 8.2 输出重定向 > 和 >>

```bash
# > 覆盖写入
echo "hello" > file.txt      # 文件内容变成 hello

# >> 追加写入
echo "world" >> file.txt     # 在文件末尾追加 world

# 把命令输出保存到文件
ls -la > filelist.txt
ps aux > processes.txt

# 把错误输出重定向
./myapp 2> error.log         # 只保存错误
./myapp > output.log 2>&1    # 标准输出和错误都保存

# 丢弃输出
./myapp > /dev/null 2>&1     # 不显示任何输出
```

### 8.3 输入重定向 <

```bash
# 从文件读取输入
./myapp < input.txt

# 例子：统计文件行数
wc -l < file.txt
```

### 📝 练习题 8.1

**问题**：解释这个命令的作用：
```bash
cat app.log | grep "ERROR" | grep -v "timeout" | wc -l
```

<details>
<summary>点击查看答案</summary>

**答案**：统计日志中包含 "ERROR" 但不包含 "timeout" 的行数。

**分解**：
1. `cat app.log` - 读取日志文件
2. `| grep "ERROR"` - 过滤出包含 ERROR 的行
3. `| grep -v "timeout"` - 排除包含 timeout 的行
4. `| wc -l` - 统计行数

</details>

---

## 9. 压缩和解压 ⭐⭐

### 9.1 tar 命令

```bash
# tar 是 Linux 最常用的压缩工具

# 压缩目录
tar -czvf myproject.tar.gz myproject/
#    │││└─ 指定文件名
#    ││└─ 显示过程
#    │└─ 使用 gzip 压缩
#    └─ 创建压缩包

# 解压
tar -xzvf myproject.tar.gz
#    └─ 解压

# 解压到指定目录
tar -xzvf myproject.tar.gz -C /home/

# 查看压缩包内容（不解压）
tar -tzvf myproject.tar.gz
```

### 9.2 zip 命令

```bash
# 压缩
zip -r myproject.zip myproject/

# 解压
unzip myproject.zip

# 解压到指定目录
unzip myproject.zip -d /home/

# 查看压缩包内容
unzip -l myproject.zip
```

---

## 10. 本章总结

### 必须熟练的命令

| 类别 | 命令 | 用途 |
|------|------|------|
| 目录 | `pwd`, `cd`, `ls -la` | 查看和切换目录 |
| 文件 | `touch`, `mkdir -p`, `cp -r`, `mv`, `rm -rf` | 文件操作 |
| 查看 | `cat`, `less`, `head`, `tail -f` | 查看文件内容 |
| 编辑 | `vim` | 编辑文件 |
| 传输 | `scp`, `rsync` | 文件传输 |
| 搜索 | `find`, `grep` | 查找文件和内容 |
| 管道 | `|`, `>`, `>>` | 组合命令 |
| 压缩 | `tar -czvf`, `tar -xzvf` | 压缩解压 |

### 记忆技巧

```
ls    = List（列出）
cd    = Change Directory（切换目录）
pwd   = Print Working Directory（打印工作目录）
cp    = Copy（复制）
mv    = Move（移动）
rm    = Remove（删除）
mkdir = Make Directory（创建目录）
cat   = Concatenate（连接/查看）
grep  = Global Regular Expression Print（全局正则表达式打印）
```

---

## 📚 下一章预告

下一章我们将学习 **Linux 进阶操作**，包括：
- 用户和权限管理
- 进程管理
- 网络操作
- 系统服务管理
- 环境变量

这些是部署和运维必备的技能！
