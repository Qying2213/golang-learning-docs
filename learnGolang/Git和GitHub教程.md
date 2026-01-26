# Git 和 GitHub 教程

> 目标：学会用 Git 管理代码，用 GitHub 托管和协作

---

## 1. 基础概念

```
工作区（Working Directory）  →  你电脑上的文件夹
暂存区（Staging Area）       →  准备提交的文件
本地仓库（Local Repository） →  你电脑上的 Git 仓库
远程仓库（Remote Repository）→  GitHub 上的仓库
```

流程：`工作区 → 暂存区 → 本地仓库 → 远程仓库`

---

## 2. 初次配置

```bash
# 设置用户名和邮箱（必须）
git config --global user.name "你的名字"
git config --global user.email "你的邮箱@example.com"

# 查看配置
git config --list
```

---

## 3. 创建仓库

### 3.1 本地新建项目

```bash
# 创建文件夹
mkdir my-project
cd my-project

# 初始化 Git 仓库
git init

# 创建文件
echo "# My Project" > README.md
```

### 3.2 克隆已有仓库

```bash
# 从 GitHub 克隆
git clone https://github.com/用户名/仓库名.git

# 克隆到指定文件夹
git clone https://github.com/用户名/仓库名.git my-folder
```

---

## 4. 日常操作（最常用）

### 4.1 查看状态

```bash
git status

# 简洁模式
git status -s
```

输出解释：
```
M  file.txt   # 已修改（Modified）
A  new.txt    # 新添加（Added）
D  old.txt    # 已删除（Deleted）
?? other.txt  # 未跟踪（Untracked）
```

### 4.2 添加到暂存区

```bash
# 添加单个文件
git add file.txt

# 添加多个文件
git add file1.txt file2.txt

# 添加所有修改的文件
git add .

# 添加所有 .go 文件
git add *.go
```

### 4.3 提交到本地仓库

```bash
# 提交（必须写提交信息）
git commit -m "描述这次改了什么"

# 添加并提交（跳过 git add，只对已跟踪文件有效）
git commit -am "修改了xxx"
```

### 4.4 推送到 GitHub

```bash
# 第一次推送（设置上游分支）
git push -u origin main

# 之后推送
git push
```

### 4.5 拉取最新代码

```bash
# 拉取并合并
git pull

# 只拉取不合并
git fetch
```

---

## 5. 完整流程示例

### 5.1 新项目上传到 GitHub

```bash
# 1. 在 GitHub 上创建新仓库（不要勾选 README）

# 2. 本地初始化
cd my-project
git init
git add .
git commit -m "first commit"

# 3. 关联远程仓库
git remote add origin https://github.com/你的用户名/仓库名.git

# 4. 推送
git branch -M main
git push -u origin main
```

### 5.2 日常开发流程

```bash
# 1. 拉取最新代码
git pull

# 2. 写代码...

# 3. 查看改了什么
git status
git diff

# 4. 添加并提交
git add .
git commit -m "添加了用户登录功能"

# 5. 推送
git push
```

---

## 6. 分支操作

### 6.1 基本操作

```bash
# 查看分支
git branch          # 本地分支
git branch -r       # 远程分支
git branch -a       # 所有分支

# 创建分支
git branch feature-login

# 切换分支
git checkout feature-login
# 或者（新版 Git）
git switch feature-login

# 创建并切换（常用）
git checkout -b feature-login
# 或者
git switch -c feature-login

# 删除分支
git branch -d feature-login      # 已合并的分支
git branch -D feature-login      # 强制删除
```

### 6.2 合并分支

```bash
# 切换到 main 分支
git checkout main

# 合并 feature-login 到 main
git merge feature-login

# 删除已合并的分支
git branch -d feature-login
```

### 6.3 推送分支到远程

```bash
# 推送新分支
git push -u origin feature-login

# 删除远程分支
git push origin --delete feature-login
```

---

## 7. 查看历史

```bash
# 查看提交历史
git log

# 简洁模式（一行一个提交）
git log --oneline

# 图形化显示分支
git log --oneline --graph

# 查看某个文件的历史
git log file.txt

# 查看具体改了什么
git log -p

# 最近 5 条
git log -5
```

---

## 8. 撤销操作

### 8.1 撤销工作区修改

```bash
# 撤销单个文件的修改（还原到上次提交）
git checkout -- file.txt
# 或者（新版）
git restore file.txt

# 撤销所有修改
git checkout -- .
git restore .
```

### 8.2 撤销暂存区

```bash
# 从暂存区移除（但保留修改）
git reset HEAD file.txt
# 或者（新版）
git restore --staged file.txt
```

### 8.3 撤销提交

```bash
# 撤销上次提交，保留修改在工作区
git reset --soft HEAD~1

# 撤销上次提交，保留修改在暂存区
git reset HEAD~1

# 撤销上次提交，丢弃所有修改（危险）
git reset --hard HEAD~1

# 修改上次提交的信息
git commit --amend -m "新的提交信息"
```

---

## 9. .gitignore 文件

告诉 Git 哪些文件不要跟踪：

```bash
# 创建 .gitignore 文件
touch .gitignore
```

常用配置：

```gitignore
# 忽略编译产物
*.exe
*.o
*.out

# 忽略依赖目录
node_modules/
vendor/

# 忽略 IDE 配置
.idea/
.vscode/
*.swp

# 忽略日志
*.log
logs/

# 忽略环境配置
.env
.env.local

# 忽略 macOS 文件
.DS_Store

# 忽略特定文件
config/secret.json

# 不忽略某个文件（例外）
!important.log
```

---

## 10. SSH 配置（免密码推送）

### 10.1 生成 SSH 密钥

```bash
# 生成密钥
ssh-keygen -t ed25519 -C "你的邮箱@example.com"

# 一路回车（使用默认路径，不设密码）

# 查看公钥
cat ~/.ssh/id_ed25519.pub
```

### 10.2 添加到 GitHub

1. 复制公钥内容
2. 打开 GitHub → Settings → SSH and GPG keys
3. 点击 New SSH key
4. 粘贴公钥，保存

### 10.3 测试连接

```bash
ssh -T git@github.com
# 成功会显示：Hi 用户名! You've successfully authenticated
```

### 10.4 使用 SSH 地址

```bash
# 克隆时用 SSH 地址
git clone git@github.com:用户名/仓库名.git

# 修改已有仓库的远程地址
git remote set-url origin git@github.com:用户名/仓库名.git
```

---

## 11. 常见问题

### 11.1 push 被拒绝

```bash
# 原因：远程有新提交，你本地没有
# 解决：先 pull 再 push
git pull
git push

# 如果有冲突，解决冲突后
git add .
git commit -m "解决冲突"
git push
```

### 11.2 合并冲突

```bash
# 冲突文件会显示：
<<<<<<< HEAD
你的代码
=======
别人的代码
>>>>>>> branch-name

# 手动编辑，保留需要的代码，删除标记
# 然后
git add .
git commit -m "解决冲突"
```

### 11.3 不小心提交了敏感信息

```bash
# 如果还没 push，撤销提交
git reset --soft HEAD~1
# 删除敏感文件，重新提交

# 如果已经 push 了
# 1. 立即更换密码/密钥
# 2. 用 git filter-branch 或 BFG 清理历史（复杂）
```

### 11.4 想要放弃所有本地修改

```bash
# 放弃所有修改，回到上次提交
git checkout -- .
git clean -fd  # 删除未跟踪的文件

# 完全同步远程（危险，会丢失本地所有修改）
git fetch origin
git reset --hard origin/main
```

---

## 12. 常用命令速查

| 命令 | 作用 |
|------|------|
| `git init` | 初始化仓库 |
| `git clone url` | 克隆仓库 |
| `git status` | 查看状态 |
| `git add .` | 添加所有文件 |
| `git commit -m "msg"` | 提交 |
| `git push` | 推送到远程 |
| `git pull` | 拉取并合并 |
| `git branch` | 查看分支 |
| `git checkout -b name` | 创建并切换分支 |
| `git merge branch` | 合并分支 |
| `git log --oneline` | 查看历史 |
| `git diff` | 查看修改 |
| `git stash` | 暂存修改 |
| `git stash pop` | 恢复暂存 |

---

## 13. 提交信息规范

好的提交信息让历史更清晰：

```bash
# 格式
<类型>: <描述>

# 类型
feat:     新功能
fix:      修复 bug
docs:     文档修改
style:    代码格式（不影响功能）
refactor: 重构
test:     测试相关
chore:    构建/工具相关

# 例子
git commit -m "feat: 添加用户登录功能"
git commit -m "fix: 修复登录页面崩溃问题"
git commit -m "docs: 更新 README"
```

---

## 练习

1. 创建一个新仓库，推送到 GitHub
2. 创建一个分支，修改代码，合并回 main
3. 配置 SSH，用 SSH 地址推送
4. 故意制造一个冲突，然后解决它
