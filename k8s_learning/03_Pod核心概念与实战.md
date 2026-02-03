# 03. Pod 核心概念与实战

> 本章目标：深入理解 Pod，掌握 YAML 编写和健康检查配置

---

## 1. 什么是 Pod？

**Pod 是 Kubernetes 的最小调度单位**。

注意：最小单位不是容器，而是 Pod。

### 1.1 Pod vs 容器

| 对比项   | 容器 (Container) | Pod              |
| :------- | :--------------- | :--------------- |
| 数量     | 一个进程         | 可以包含多个容器 |
| 网络     | 独立网络命名空间 | 内部容器共享网络 |
| 存储     | 独立             | 可以共享存储卷   |
| 生命周期 | 独立管理         | K8s 统一管理     |

### 1.2 为什么需要 Pod？

有些应用场景需要多个进程紧密协作：

- Web 服务器 + 日志收集器
- 主应用 + Sidecar 代理
- 应用 + 配置重载器

这些容器需要：

1. **共享网络**：通过 localhost 互相通信
2. **共享存储**：访问同一个文件目录
3. **同生共死**：一起创建、一起销毁

Pod 就是实现这个需求的抽象。

---

## 2. Pod 的内部结构

```
┌─────────────────────────────────────────────────┐
│                     Pod                          │
│  ┌───────────────────────────────────────────┐  │
│  │           Pause 容器 (Infra Container)     │  │
│  │     持有网络命名空间，IP: 10.244.0.5       │  │
│  └───────────────────────────────────────────┘  │
│                      │                           │
│        ┌─────────────┴─────────────┐            │
│        ▼                           ▼            │
│  ┌───────────┐              ┌───────────┐       │
│  │ Container │              │ Container │       │
│  │   主应用   │  localhost   │  Sidecar  │       │
│  └───────────┘              └───────────┘       │
│        │                           │            │
│        └───────────┬───────────────┘            │
│                    ▼                            │
│            ┌─────────────┐                      │
│            │   Volume    │                      │
│            │  共享存储    │                      │
│            └─────────────┘                      │
└─────────────────────────────────────────────────┘
```

**Pause 容器**：

- 每个 Pod 都有一个隐藏的 Pause 容器
- 它的唯一职责：持有网络命名空间
- 其他容器"加入"到这个网络空间
- 这就是为什么 Pod 内容器共享 IP

---

## 3. Pod 生命周期

### 3.1 状态（Phase）

| 状态          | 含义                                                |
| :------------ | :-------------------------------------------------- |
| **Pending**   | Pod 已被 K8s 接收，但还没调度到节点，或正在拉取镜像 |
| **Running**   | 至少一个容器正在运行                                |
| **Succeeded** | 所有容器正常退出（exit 0），不会重启                |
| **Failed**    | 至少一个容器异常退出                                |
| **Unknown**   | 无法获取状态（通常是节点失联）                      |

### 3.2 重启策略（restartPolicy）

| 策略        | 行为             | 适用场景                       |
| :---------- | :--------------- | :----------------------------- |
| `Always`    | 总是重启（默认） | Web 服务、API 等长期运行的应用 |
| `OnFailure` | 只在失败时重启   | 批处理任务（Job）              |
| `Never`     | 从不重启         | 一次性脚本                     |

---

## 4. 编写 Pod YAML

### 4.1 基础模板

创建文件 `pod-basic.yaml`：

```yaml
apiVersion: v1 # API 版本
kind: Pod # 资源类型
metadata:
  name: my-nginx # Pod 名称
  labels: # 标签（用于选择和分组）
    app: nginx
    env: dev
spec:
  containers: # 容器列表
    - name: nginx # 容器名称
      image: nginx:alpine # 镜像
      ports:
        - containerPort: 80 # 容器监听的端口
```

**应用它**：

```bash
kubectl apply -f pod-basic.yaml
```

**查看结果**：

```bash
kubectl get pods
kubectl describe pod my-nginx
```

### 4.2 完整生产级模板

创建文件 `pod-production.yaml`：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: production-app
  labels:
    app: myapp
    version: v1
spec:
  restartPolicy: Always

  # 初始化容器（在主容器之前运行）
  initContainers:
    - name: init-check-db
      image: busybox:1.36
      command:
        ["sh", "-c", 'echo "等待数据库就绪..."; sleep 5; echo "检查完成"']

  # 主容器
  containers:
    - name: app
      image: nginx:alpine

      # 端口
      ports:
        - containerPort: 80
          name: http

      # 资源限制（非常重要！）
      resources:
        requests: # 最小需求（用于调度）
          memory: "64Mi"
          cpu: "100m" # 0.1 核
        limits: # 最大限制（超过会被杀掉）
          memory: "128Mi"
          cpu: "200m" # 0.2 核

      # 环境变量
      env:
        - name: APP_ENV
          value: "production"
        - name: LOG_LEVEL
          value: "info"

      # 存活探针（判断容器是否活着）
      livenessProbe:
        httpGet:
          path: /
          port: 80
        initialDelaySeconds: 10 # 启动后等 10 秒再检查
        periodSeconds: 5 # 每 5 秒检查一次
        failureThreshold: 3 # 连续失败 3 次才判定为死亡

      # 就绪探针（判断容器是否能接收流量）
      readinessProbe:
        httpGet:
          path: /
          port: 80
        initialDelaySeconds: 5
        periodSeconds: 3

      # 启动探针（用于启动慢的应用）
      startupProbe:
        httpGet:
          path: /
          port: 80
        failureThreshold: 30 # 允许 30 次失败
        periodSeconds: 10 # 每 10 秒检查，共 300 秒启动时间
```

---

## 5. 三种探针详解

### 5.1 Liveness Probe（存活探针）

**问题**：容器进程还在，但程序死锁了，怎么办？

**解决**：Liveness Probe 定期检查，失败就重启容器

```yaml
livenessProbe:
  httpGet: # HTTP 方式
    path: /healthz
    port: 8080
  tcpSocket: # 或 TCP 方式
    port: 8080
  exec: # 或命令方式
    command:
      - cat
      - /tmp/healthy
```

### 5.2 Readiness Probe（就绪探针）

**问题**：应用刚启动，还在加载数据，不能处理请求

**解决**：Readiness Probe 失败时，从 Service 负载均衡中移除，不给它流量

**区别**：

- Liveness 失败 → 重启容器
- Readiness 失败 → 停止流量，但不重启

### 5.3 Startup Probe（启动探针）

**问题**：老应用启动需要 5 分钟，Liveness 等不了那么久

**解决**：Startup Probe 专门用于启动阶段，成功后才交给 Liveness

---

## 6. 资源限制详解

### 6.1 Requests vs Limits

```yaml
resources:
  requests: # 最小保证，调度时看这个
    memory: "64Mi"
    cpu: "100m"
  limits: # 最大限制，超过就杀掉
    memory: "128Mi"
    cpu: "200m"
```

**CPU 单位**：

- `1` = 1 核
- `500m` = 0.5 核
- `100m` = 0.1 核

**内存单位**：

- `Ki` = 1024 bytes
- `Mi` = 1024 Ki
- `Gi` = 1024 Mi

### 6.2 QoS 等级

K8s 根据 requests/limits 配置，自动给 Pod 分配 QoS 等级：

| QoS 等级   | 条件              | 优先级         |
| :--------- | :---------------- | :------------- |
| Guaranteed | requests = limits | 最高，最后被杀 |
| Burstable  | requests < limits | 中等           |
| BestEffort | 没有设置          | 最低，最先被杀 |

---

## 7. 多容器 Pod 实战

### 7.1 Sidecar 模式

主应用 + 日志收集器：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: sidecar-demo
spec:
  containers:
    # 主应用
    - name: app
      image: nginx:alpine
      volumeMounts:
        - name: logs
          mountPath: /var/log/nginx

    # Sidecar：日志收集
    - name: log-collector
      image: busybox:1.36
      command: ["sh", "-c", "tail -f /logs/access.log"]
      volumeMounts:
        - name: logs
          mountPath: /logs

  # 共享存储卷
  volumes:
    - name: logs
      emptyDir: {}
```

### 7.2 Init Container 模式

等待依赖服务就绪：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  initContainers:
    - name: wait-for-db
      image: busybox:1.36
      command:
        [
          "sh",
          "-c",
          "until nc -z mysql 3306; do echo waiting...; sleep 2; done",
        ]

  containers:
    - name: app
      image: myapp:v1
```

---

## 8. 实战练习

### 练习 1：创建并观察 Pod

```bash
# 创建 Pod
kubectl apply -f pod-production.yaml

# 观察状态变化
kubectl get pods -w

# 查看详细信息
kubectl describe pod production-app

# 查看日志
kubectl logs production-app

# 进入容器
kubectl exec -it production-app -- sh
```

### 练习 2：模拟探针失败

创建 `pod-probe-fail.yaml`：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: probe-fail
spec:
  containers:
    - name: app
      image: nginx:alpine
      livenessProbe:
        httpGet:
          path: /不存在的路径
          port: 80
        initialDelaySeconds: 5
        periodSeconds: 3
```

```bash
kubectl apply -f pod-probe-fail.yaml
kubectl get pods -w
# 观察 RESTARTS 列不断增加
```

### 练习 3：清理资源

```bash
kubectl delete pod my-nginx
kubectl delete pod production-app
kubectl delete pod probe-fail
```

---

## 9. 常见问题

### Pod 一直 Pending

**排查**：

```bash
kubectl describe pod <name>
```

**常见原因**：

- 资源不足（没有节点满足 requests）
- 镜像拉取失败
- 节点选择器不匹配

### Pod 一直 CrashLoopBackOff

**含义**：容器反复崩溃重启

**排查**：

```bash
kubectl logs <pod-name>
kubectl logs <pod-name> --previous  # 查看上一次崩溃的日志
```

### ImagePullBackOff

**含义**：无法拉取镜像

**排查**：

- 镜像名称是否正确
- 私有仓库是否需要认证
- 网络是否可达

---

## ✅ 本章检查点

- [ ] 理解 Pod 和容器的区别
- [ ] 能编写包含探针的 Pod YAML
- [ ] 理解 Liveness、Readiness、Startup 三种探针的区别
- [ ] 能配置资源 requests 和 limits
- [ ] 能使用 `kubectl logs`、`kubectl exec` 调试 Pod

---

## ⏭️ 下一章

实际生产中，我们不会直接创建 Pod，而是通过 Deployment 来管理。

👉 [04_Deployment与副本管理.md](04_Deployment与副本管理.md)
