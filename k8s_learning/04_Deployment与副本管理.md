# 04. Deployment 与副本管理

> 本章目标：掌握 Deployment 的核心功能，包括副本管理、滚动更新和回滚

---

## 1. 为什么不直接用 Pod？

直接创建 Pod 有几个致命问题：

| 问题             | 直接创建 Pod       | 使用 Deployment      |
| :--------------- | :----------------- | :------------------- |
| Pod 挂了怎么办？ | 挂了就没了         | 自动重新创建         |
| 想要 3 个副本？  | 手动创建 3 次      | 设置 replicas=3      |
| 升级版本？       | 删除旧的、创建新的 | 滚动更新，不中断服务 |
| 发布失败？       | 手动恢复           | 一键回滚             |

**结论**：生产环境永远使用 Deployment，不直接创建 Pod

---

## 2. Deployment 控制链

```
┌─────────────────────────────────────────────────────────────────┐
│                      用户创建 Deployment                         │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Deployment                                 │
│   spec:                                                          │
│     replicas: 3                                                  │
│     template:                                                    │
│       spec:                                                      │
│         containers:                                              │
│         - image: nginx:1.14                                      │
└───────────────────────────┬─────────────────────────────────────┘
                            │ 创建
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      ReplicaSet                                  │
│   metadata:                                                      │
│     name: nginx-deploy-5d4f987                                   │
│   spec:                                                          │
│     replicas: 3                                                  │
└───────────────────────────┬─────────────────────────────────────┘
                            │ 创建
                            ▼
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         ┌────────┐   ┌────────┐   ┌────────┐
         │ Pod 1  │   │ Pod 2  │   │ Pod 3  │
         └────────┘   └────────┘   └────────┘
```

**关系**：

- **Deployment** 管理 **ReplicaSet**
- **ReplicaSet** 管理 **Pod**
- 我们只需要操作 Deployment，其他的 K8s 自动处理

---

## 3. 编写 Deployment YAML

### 3.1 基础模板

创建 `deploy-nginx.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deploy # Deployment 名称
  labels:
    app: nginx
spec:
  replicas: 3 # 副本数量

  # 选择器：Deployment 通过标签找到它管理的 Pod
  selector:
    matchLabels:
      app: nginx

  # Pod 模板
  template:
    metadata:
      labels:
        app: nginx # 必须匹配上面的 selector
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
```

### 3.2 部署并验证

```bash
# 应用配置
kubectl apply -f deploy-nginx.yaml

# 查看 Deployment
kubectl get deployments
# NAME           READY   UP-TO-DATE   AVAILABLE   AGE
# nginx-deploy   3/3     3            3           30s

# 查看 ReplicaSet（自动创建的）
kubectl get rs
# NAME                      DESIRED   CURRENT   READY   AGE
# nginx-deploy-5d4f987xxx   3         3         3       30s

# 查看 Pod
kubectl get pods
# NAME                            READY   STATUS    RESTARTS   AGE
# nginx-deploy-5d4f987xxx-abc12   1/1     Running   0          30s
# nginx-deploy-5d4f987xxx-def34   1/1     Running   0          30s
# nginx-deploy-5d4f987xxx-ghi56   1/1     Running   0          30s
```

---

## 4. 副本管理（Scaling）

### 4.1 扩容

```bash
# 方法 1：命令行扩容
kubectl scale deployment nginx-deploy --replicas=5

# 方法 2：修改 YAML 后 apply
# 把 replicas: 3 改成 replicas: 5
kubectl apply -f deploy-nginx.yaml

# 方法 3：使用 edit 直接编辑
kubectl edit deployment nginx-deploy
```

### 4.2 缩容

```bash
kubectl scale deployment nginx-deploy --replicas=2
```

### 4.3 观察变化

```bash
# 实时观察 Pod 变化
kubectl get pods -w
```

---

## 5. 滚动更新（Rolling Update）

这是 Deployment 最强大的功能：**无停机更新**

### 5.1 触发更新

假设要把 nginx 从 1.14.2 升级到 1.16.1：

```bash
# 方法 1：命令行更新
kubectl set image deployment/nginx-deploy nginx=nginx:1.16.1

# 方法 2：修改 YAML 后 apply
# 把 image: nginx:1.14.2 改成 image: nginx:1.16.1
kubectl apply -f deploy-nginx.yaml
```

### 5.2 观察更新过程

```bash
# 查看更新状态
kubectl rollout status deployment/nginx-deploy
# Waiting for deployment "nginx-deploy" rollout to finish: 1 out of 3 new replicas have been updated...
# Waiting for deployment "nginx-deploy" rollout to finish: 2 out of 3 new replicas have been updated...
# deployment "nginx-deploy" successfully rolled out

# 实时观察 Pod 变化
kubectl get pods -w
```

### 5.3 更新过程图解

```
更新前：
  ReplicaSet v1: [Pod1-v1] [Pod2-v1] [Pod3-v1]

更新中 Step 1:
  ReplicaSet v1: [Pod1-v1] [Pod2-v1] [Pod3-v1]
  ReplicaSet v2: [Pod1-v2]                        ← 新建 1 个

更新中 Step 2:
  ReplicaSet v1: [Pod1-v1] [Pod2-v1]              ← 删除 1 个
  ReplicaSet v2: [Pod1-v2] [Pod2-v2]              ← 新建 1 个

更新中 Step 3:
  ReplicaSet v1: [Pod1-v1]                        ← 删除 1 个
  ReplicaSet v2: [Pod1-v2] [Pod2-v2] [Pod3-v2]    ← 新建 1 个

更新完成:
  ReplicaSet v1: (空，但保留用于回滚)
  ReplicaSet v2: [Pod1-v2] [Pod2-v2] [Pod3-v2]
```

---

## 6. 更新策略配置

### 6.1 RollingUpdate 策略（默认）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deploy
spec:
  replicas: 10
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25% # 更新时最多可以多出 25% 的 Pod
      maxUnavailable: 25% # 更新时最多可以有 25% 的 Pod 不可用
```

**参数解释**：

- `maxSurge: 25%`：10 个副本，更新时最多同时存在 13 个 Pod
- `maxUnavailable: 25%`：10 个副本，更新时最少保持 8 个可用

### 6.2 Recreate 策略

先删除所有旧 Pod，再创建新 Pod（会有短暂停机）

```yaml
spec:
  strategy:
    type: Recreate
```

**适用场景**：不支持多版本同时运行的应用

---

## 7. 回滚（Rollback）

更新出问题了？一键回滚！

### 7.1 查看历史版本

```bash
kubectl rollout history deployment/nginx-deploy
# REVISION  CHANGE-CAUSE
# 1         <none>
# 2         <none>
```

### 7.2 回滚到上一个版本

```bash
kubectl rollout undo deployment/nginx-deploy
```

### 7.3 回滚到指定版本

```bash
kubectl rollout undo deployment/nginx-deploy --to-revision=1
```

### 7.4 记录变更原因

```bash
# 更新时加 --record（已废弃，但仍可用）
kubectl set image deployment/nginx-deploy nginx=nginx:1.17.0 --record

# 或使用 annotate
kubectl annotate deployment/nginx-deploy kubernetes.io/change-cause="升级到 1.17.0"
```

---

## 8. 暂停和恢复更新

大规模更改时，可以先暂停，改完再统一生效

```bash
# 暂停
kubectl rollout pause deployment/nginx-deploy

# 做多次修改
kubectl set image deployment/nginx-deploy nginx=nginx:1.17.0
kubectl set resources deployment/nginx-deploy -c=nginx --limits=cpu=300m

# 恢复（所有修改一次性生效）
kubectl rollout resume deployment/nginx-deploy
```

---

## 9. 实战练习

### 练习 1：完整的部署-更新-回滚流程

```bash
# 1. 创建 Deployment（版本 1.14.2）
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
EOF

# 2. 验证
kubectl get deploy,rs,pods

# 3. 升级到 1.16.1
kubectl set image deployment/web-app nginx=nginx:1.16.1

# 4. 观察更新
kubectl rollout status deployment/web-app

# 5. 查看历史
kubectl rollout history deployment/web-app

# 6. 演示回滚
kubectl rollout undo deployment/web-app

# 7. 验证版本回退
kubectl describe deployment web-app | grep Image

# 8. 清理
kubectl delete deployment web-app
```

### 练习 2：测试副本自愈

```bash
# 创建 Deployment
kubectl create deployment healer --image=nginx:alpine --replicas=3

# 手动删除一个 Pod
kubectl delete pod $(kubectl get pods -l app=healer -o name | head -1)

# 观察自动恢复
kubectl get pods -w
# 你会看到新 Pod 立即被创建

# 清理
kubectl delete deployment healer
```

---

## 10. 最佳实践

### 10.1 始终设置资源限制

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 200m
    memory: 256Mi
```

### 10.2 使用明确的镜像标签

```yaml
# 不好
image: nginx
image: nginx:latest

# 好
image: nginx:1.16.1
image: myapp:v2.3.1
```

### 10.3 配置探针

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
```

### 10.4 合理配置更新策略

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1
    maxUnavailable: 0 # 保证零停机
```

---

## ✅ 本章检查点

- [ ] 理解 Deployment → ReplicaSet → Pod 的控制链
- [ ] 能编写 Deployment YAML
- [ ] 能执行扩缩容操作
- [ ] 理解滚动更新的过程
- [ ] 能执行回滚操作

---

## ⏭️ 下一章

应用部署好了，但外部如何访问？Service 来解决网络问题。

👉 [05_Service与网络通信.md](05_Service与网络通信.md)
