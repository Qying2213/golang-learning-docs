# 01. 环境搭建：Minikube 安装与启动

> 本章目标：在本地搭建一个可用的 Kubernetes 集群

---

## 1. 什么是 Minikube？

**Minikube** 是一个官方提供的工具，用于在**本地电脑上**运行一个单节点（或多节点）的 Kubernetes 集群。

**它不是玩具**：Minikube 内部运行的是**真正的 Kubernetes**，和阿里云 ACK、AWS EKS 用的是同一套代码。只是规模缩小到了一台机器上。

### Minikube vs 生产集群

| 对比项   | Minikube          | 生产集群 (如 EKS) |
| :------- | :---------------- | :---------------- |
| 节点数量 | 1-3 个 (本地模拟) | 几十到上千个      |
| 用途     | 学习、开发、测试  | 生产环境          |
| 高可用   | 无                | 有 (多 Master)    |
| 费用     | 免费              | 按量付费          |

---

## 2. 环境检查

你的机器上已经安装好了 Minikube 和 kubectl，我们先验证一下：

```bash
# 检查 Minikube 版本
minikube version

# 检查 kubectl 版本
kubectl version --client
```

**期望输出**:

```
minikube version: v1.38.0
Client Version: v1.35.0
```

如果没有安装，请参考下面的安装步骤。

---

## 3. 安装步骤（如果尚未安装）

### 3.1 安装 Minikube (macOS)

```bash
# 使用 Homebrew 安装
brew install minikube
```

### 3.2 安装 kubectl (macOS)

```bash
# 使用 Homebrew 安装
brew install kubectl
```

### 3.3 安装 Docker Desktop

Minikube 需要一个"驱动"来运行虚拟机或容器。推荐使用 Docker。

1. 下载 [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/)
2. 安装并启动
3. 确保 Docker 正在运行（状态栏有 🐳 图标）

---

## 4. 启动你的第一个集群

### 4.1 基本启动

```bash
minikube start
```

这条命令会：

1. 下载 Kubernetes 镜像（首次较慢）
2. 创建一个虚拟机或 Docker 容器
3. 在里面部署 K8s 所有组件
4. 配置 kubectl 连接到这个集群

**首次启动可能需要 3-10 分钟**，取决于网络速度。

### 4.2 国内网络优化启动

如果下载镜像很慢，使用阿里云镜像源：

```bash
minikube start --image-mirror-country='cn'
```

### 4.3 指定资源配置

默认配置可能不够用，建议指定 CPU 和内存：

```bash
minikube start --cpus=4 --memory=4096
```

### 4.4 使用 Docker 驱动（推荐）

```bash
minikube start --driver=docker
```

---

## 5. 验证集群状态

### 5.1 检查 Minikube 状态

```bash
minikube status
```

**期望输出**:

```
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

所有状态都是 `Running` 才算正常。

### 5.2 检查节点

```bash
kubectl get nodes
```

**期望输出**:

```
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   1m    v1.32.0
```

`STATUS` 必须是 `Ready`。

### 5.3 查看集群信息

```bash
kubectl cluster-info
```

**期望输出**:

```
Kubernetes control plane is running at https://192.168.49.2:8443
CoreDNS is running at https://192.168.49.2:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```

---

## 6. 启动 Dashboard（可视化界面）

Kubernetes 自带一个 Web 管理界面，对初学者非常友好。

```bash
minikube dashboard
```

这条命令会自动打开浏览器。如果没有自动打开，复制终端输出的 URL 手动访问。

**注意**：这个命令会占用当前终端，需要时按 `Ctrl+C` 停止。

---

## 7. 常用 Minikube 命令

| 命令                         | 作用                 | 使用场景           |
| :--------------------------- | :------------------- | :----------------- |
| `minikube start`             | 启动集群             | 每天开始工作时     |
| `minikube stop`              | 停止集群（保留数据） | 下班休息时         |
| `minikube delete`            | 删除集群（清空一切） | 想重头再来时       |
| `minikube pause`             | 暂停集群（最省资源） | 临时不用但不想关机 |
| `minikube unpause`           | 恢复暂停的集群       | 继续使用           |
| `minikube status`            | 查看状态             | 检查是否正常       |
| `minikube ssh`               | 登录到集群虚拟机内部 | 高级调试           |
| `minikube dashboard`         | 打开 Web 界面        | 可视化管理         |
| `minikube addons list`       | 列出所有插件         | 查看可用功能       |
| `minikube addons enable xxx` | 启用插件             | 添加功能           |

---

## 8. 实战练习

### 练习 1：完整的启动-停止流程

```bash
# 1. 启动集群
minikube start

# 2. 检查状态
minikube status

# 3. 查看节点
kubectl get nodes

# 4. 停止集群
minikube stop

# 5. 再次检查状态
minikube status
# 此时应该显示 Stopped
```

### 练习 2：SSH 进入集群内部

```bash
# 登录到 Minikube 虚拟机
minikube ssh

# 在虚拟机内部查看运行的容器
docker ps

# 退出虚拟机
exit
```

### 练习 3：启用 Ingress 插件

```bash
# 查看插件列表
minikube addons list

# 启用 Ingress（后面章节会用到）
minikube addons enable ingress

# 验证 Ingress 是否启动
kubectl get pods -n ingress-nginx
```

---

## 9. 常见问题排查

### 问题 1：minikube start 卡住不动

**原因**：网络问题，无法下载镜像

**解决**：

```bash
minikube delete
minikube start --image-mirror-country='cn'
```

### 问题 2：显示 "Exiting due to PROVIDER_DOCKER_NOT_RUNNING"

**原因**：Docker 没有启动

**解决**：启动 Docker Desktop 应用

### 问题 3：kubectl 命令报错 "connection refused"

**原因**：集群没有启动

**解决**：

```bash
minikube start
```

### 问题 4：节点状态是 NotReady

**原因**：集群还在初始化

**解决**：等待 1-2 分钟再检查，或运行：

```bash
kubectl get nodes -w
# 持续观察直到变成 Ready
```

---

## ✅ 本章检查点

完成本章后，请确认你可以：

- [ ] 运行 `minikube start` 并成功启动集群
- [ ] 运行 `kubectl get nodes` 看到节点状态为 Ready
- [ ] 运行 `minikube dashboard` 打开 Web 界面
- [ ] 运行 `minikube ssh` 登录到集群内部
- [ ] 运行 `minikube stop` 停止集群

---

## ⏭️ 下一章

环境准备好了，接下来我们要深入了解 Kubernetes 的内部架构，理解它是如何运作的。

👉 [02_K8s架构深度解析.md](02_K8s架构深度解析.md)
