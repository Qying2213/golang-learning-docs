# GitHub Actions 入门 ⭐⭐⭐

> 本章目标：掌握 GitHub Actions 的核心概念和基本语法，能写出第一个 workflow

---

## 1. GitHub Actions 是什么？

GitHub Actions 是 GitHub 提供的 CI/CD 服务，让你可以在 GitHub 仓库中自动化执行任务。

**核心特点：**
- 🆓 免费（公开仓库无限制，私有仓库每月2000分钟）
- 📦 开箱即用，无需额外服务器
- 🔌 丰富的插件市场（Actions Marketplace）
- 🔗 与 GitHub 深度集成

---

## 2. 核心概念

### 2.1 概念关系图

```
┌─────────────────────────────────────────────────────────────┐
│                    Workflow（工作流）                         │
│                    .github/workflows/ci.yml                  │
│  ┌─────────────────────────────────────────────────────────┐│
│  │                    Job 1: build                         ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   ││
│  │  │ Step 1  │→ │ Step 2  │→ │ Step 3  │→ │ Step 4  │   ││
│  │  │拉取代码 │  │安装依赖 │  │  编译   │  │  测试   │   ││
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘   ││
│  └─────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────┐│
│  │                    Job 2: deploy                        ││
│  │  ┌─────────┐  ┌─────────┐                              ││
│  │  │ Step 1  │→ │ Step 2  │                              ││
│  │  │构建镜像 │  │  部署   │                              ││
│  │  └─────────┘  └─────────┘                              ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心术语

| 术语 | 说明 | 类比 |
|------|------|------|
| **Workflow** | 工作流，一个完整的自动化流程 | 一条生产线 |
| **Event** | 触发工作流的事件（push、PR等） | 按下启动按钮 |
| **Job** | 工作流中的一个任务 | 生产线上的一个工位 |
| **Step** | Job 中的一个步骤 | 工位上的一个操作 |
| **Action** | 可复用的操作单元 | 标准化的工具 |
| **Runner** | 执行 Job 的服务器 | 工人 |

### 2.3 文件位置

```
你的项目/
├── main.go
├── go.mod
└── .github/
    └── workflows/
        ├── ci.yml        ← 可以有多个 workflow
        ├── deploy.yml
        └── test.yml
```

**重要：** 必须放在 `.github/workflows/` 目录下，文件名随意，后缀是 `.yml` 或 `.yaml`

---

## 3. Workflow 语法详解

### 3.1 最简单的 Workflow

```yaml
# .github/workflows/hello.yml

name: Hello World              # workflow 的名字

on: push                       # 触发条件：push 时触发

jobs:                          # 定义任务
  say-hello:                   # job 的 ID
    runs-on: ubuntu-latest     # 运行环境
    steps:                     # 步骤
      - name: Say Hello        # 步骤名称
        run: echo "Hello, World!"  # 执行的命令
```

**执行效果：** 每次 push 代码，都会在 GitHub 的服务器上执行 `echo "Hello, World!"`

### 3.2 触发条件（on）

```yaml
# 方式1：单个事件
on: push

# 方式2：多个事件
on: [push, pull_request]

# 方式3：详细配置
on:
  push:
    branches:
      - main           # 只有 push 到 main 分支才触发
      - develop
    paths:
      - 'src/**'       # 只有 src 目录下的文件变化才触发
      - '*.go'
    paths-ignore:
      - '**.md'        # 忽略 md 文件的变化
  
  pull_request:
    branches:
      - main           # PR 到 main 分支时触发
  
  schedule:
    - cron: '0 2 * * *'  # 定时任务：每天凌晨2点执行
  
  workflow_dispatch:     # 手动触发（在 GitHub 页面点击按钮）
```

**常用触发事件：**

| 事件 | 说明 |
|------|------|
| `push` | 推送代码时 |
| `pull_request` | 创建/更新 PR 时 |
| `schedule` | 定时执行 |
| `workflow_dispatch` | 手动触发 |
| `release` | 发布 Release 时 |

### 3.3 运行环境（runs-on）

```yaml
jobs:
  build:
    runs-on: ubuntu-latest    # 最常用，Linux 环境
    # runs-on: ubuntu-22.04   # 指定版本
    # runs-on: macos-latest   # macOS 环境
    # runs-on: windows-latest # Windows 环境
```

**推荐：** 大多数情况用 `ubuntu-latest` 就够了

### 3.4 步骤（steps）

```yaml
steps:
  # 方式1：运行命令
  - name: 运行单个命令
    run: echo "Hello"
  
  # 方式2：运行多个命令
  - name: 运行多个命令
    run: |
      echo "第一行"
      echo "第二行"
      go version
  
  # 方式3：使用现成的 Action
  - name: 拉取代码
    uses: actions/checkout@v3
  
  # 方式4：带参数的 Action
  - name: 设置 Go 环境
    uses: actions/setup-go@v4
    with:
      go-version: '1.21'
```

### 3.5 环境变量

```yaml
# 全局环境变量
env:
  APP_NAME: myapp
  GO_VERSION: '1.21'

jobs:
  build:
    runs-on: ubuntu-latest
    # Job 级别环境变量
    env:
      DATABASE_URL: localhost:3306
    
    steps:
      - name: 使用环境变量
        # Step 级别环境变量
        env:
          STEP_VAR: hello
        run: |
          echo "App: $APP_NAME"
          echo "DB: $DATABASE_URL"
          echo "Step: $STEP_VAR"
```

### 3.6 Secrets（密钥）

敏感信息（密码、Token）不能写在配置文件里，要用 Secrets：

**设置 Secrets：**
1. 进入 GitHub 仓库
2. Settings → Secrets and variables → Actions
3. New repository secret
4. 添加 Name 和 Value

**使用 Secrets：**
```yaml
steps:
  - name: 登录 Docker Hub
    run: |
      echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
```

---

## 4. 常用 Actions

### 4.1 必备 Actions

```yaml
# 1. 拉取代码（几乎每个 workflow 都需要）
- uses: actions/checkout@v3

# 2. 设置 Go 环境
- uses: actions/setup-go@v4
  with:
    go-version: '1.21'

# 3. 设置 Node.js 环境
- uses: actions/setup-node@v3
  with:
    node-version: '18'

# 4. 缓存依赖（加速构建）
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

# 5. 上传构建产物
- uses: actions/upload-artifact@v3
  with:
    name: my-artifact
    path: ./dist/
```

### 4.2 Docker 相关

```yaml
# 登录 Docker Hub
- uses: docker/login-action@v2
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}

# 构建并推送镜像
- uses: docker/build-push-action@v4
  with:
    push: true
    tags: username/app:latest
```

---

## 5. 多 Job 配置

### 5.1 并行执行

```yaml
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Job 1"
  
  job2:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Job 2"

# job1 和 job2 会同时执行（并行）
```

### 5.2 顺序执行（依赖）

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Building..."
  
  test:
    needs: build              # test 依赖 build，build 完成后才执行
    runs-on: ubuntu-latest
    steps:
      - run: echo "Testing..."
  
  deploy:
    needs: [build, test]      # deploy 依赖 build 和 test
    runs-on: ubuntu-latest
    steps:
      - run: echo "Deploying..."
```

**执行顺序：**
```
build → test → deploy
```

### 5.3 条件执行

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    # 只在 main 分支执行
    if: github.ref == 'refs/heads/main'
    steps:
      - run: echo "Deploying to production"
  
  deploy-staging:
    runs-on: ubuntu-latest
    # 只在 develop 分支执行
    if: github.ref == 'refs/heads/develop'
    steps:
      - run: echo "Deploying to staging"
```

---

## 6. 实战：第一个完整 Workflow

### 6.1 Go 项目 CI

```yaml
# .github/workflows/ci.yml
name: Go CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'

jobs:
  lint:
    name: 代码检查
    runs-on: ubuntu-latest
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: 运行 golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  test:
    name: 运行测试
    runs-on: ubuntu-latest
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: 缓存依赖
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: 下载依赖
        run: go mod download
      
      - name: 运行测试
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: 上传覆盖率报告
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage.out

  build:
    name: 构建
    runs-on: ubuntu-latest
    needs: [lint, test]  # 依赖 lint 和 test
    steps:
      - name: 拉取代码
        uses: actions/checkout@v3
      
      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: 构建
        run: |
          CGO_ENABLED=0 GOOS=linux go build -o app main.go
      
      - name: 上传构建产物
        uses: actions/upload-artifact@v3
        with:
          name: app-binary
          path: app
```

### 6.2 执行流程图

```
push 代码
    ↓
┌───────────────────────────────────┐
│           并行执行                 │
│  ┌─────────┐    ┌─────────┐      │
│  │  lint   │    │  test   │      │
│  │代码检查 │    │运行测试 │      │
│  └────┬────┘    └────┬────┘      │
│       └──────┬───────┘           │
└──────────────┼───────────────────┘
               ↓
         ┌─────────┐
         │  build  │
         │  构建   │
         └─────────┘
```

---

## 7. 查看执行结果

### 7.1 在 GitHub 上查看

1. 进入你的仓库
2. 点击 **Actions** 标签
3. 选择一个 workflow 运行记录
4. 查看每个 Job 和 Step 的日志

### 7.2 状态徽章

在 README 中添加构建状态徽章：

```markdown
![CI](https://github.com/你的用户名/你的仓库/actions/workflows/ci.yml/badge.svg)
```

效果：![CI](https://img.shields.io/badge/CI-passing-brightgreen)

---

## 8. 练习题

### 练习1：语法理解

**Q1：以下配置会在什么时候触发？**
```yaml
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
```

<details>
<summary>查看答案</summary>

两种情况会触发：
1. 当代码 push 到 main 分支时
2. 当有 PR 请求合并到 main 分支时

</details>

**Q2：`needs: [build, test]` 是什么意思？**

<details>
<summary>查看答案</summary>

表示当前 Job 依赖 build 和 test 两个 Job，只有这两个 Job 都成功完成后，当前 Job 才会执行。

</details>

**Q3：`runs-on: ubuntu-latest` 是什么意思？**

<details>
<summary>查看答案</summary>

指定 Job 运行在 GitHub 提供的最新版 Ubuntu 虚拟机上。GitHub 会自动分配一台云服务器来执行你的任务。

</details>

### 练习2：写一个 Workflow

**任务：** 写一个 workflow，要求：
1. 在 push 到 main 分支时触发
2. 打印 "Hello, CI/CD!"
3. 打印当前日期

<details>
<summary>查看答案</summary>

```yaml
name: Hello CI/CD

on:
  push:
    branches: [ main ]

jobs:
  hello:
    runs-on: ubuntu-latest
    steps:
      - name: Say Hello
        run: echo "Hello, CI/CD!"
      
      - name: Print Date
        run: date
```

</details>

### 练习3：修复错误

**以下配置有什么问题？**
```yaml
name: My Workflow

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: go build
```

<details>
<summary>查看答案</summary>

问题：
1. 没有拉取代码（缺少 `actions/checkout`）
2. 没有设置 Go 环境（缺少 `actions/setup-go`）

修复后：
```yaml
name: My Workflow

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go build
```

</details>

---

## 9. 本章小结

```
GitHub Actions 核心概念：
├── Workflow：完整的自动化流程（.yml 文件）
├── Event：触发条件（push、PR、定时等）
├── Job：一个任务，包含多个 Step
├── Step：一个步骤，执行命令或 Action
└── Action：可复用的操作单元

配置文件位置：.github/workflows/*.yml

常用语法：
├── on：触发条件
├── jobs：定义任务
├── runs-on：运行环境
├── steps：执行步骤
├── uses：使用 Action
├── run：执行命令
├── env：环境变量
├── needs：Job 依赖
└── if：条件执行
```

---

下一章：[03-Go项目实战](./03-Go项目实战.md)
