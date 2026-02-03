# 08. kubectl 常用命令速查表

> 日常运维必备命令汇总，建议收藏

---

## 📋 基础命令

### 集群信息

```bash
# 集群信息
kubectl cluster-info

# 查看节点
kubectl get nodes
kubectl get nodes -o wide              # 显示更多信息（IP、版本等）

# 查看 API 版本
kubectl api-versions

# 查看 API 资源
kubectl api-resources
```

### 配置上下文

```bash
# 查看当前上下文
kubectl config current-context

# 查看所有上下文
kubectl config get-contexts

# 切换上下文
kubectl config use-context <context-name>

# 设置默认命名空间
kubectl config set-context --current --namespace=<namespace>
```

---

## 📦 资源查看（GET）

### 基础查看

```bash
# 查看 Pod
kubectl get pods
kubectl get pods -n kube-system        # 指定命名空间
kubectl get pods -A                    # 所有命名空间
kubectl get pods -o wide               # 详细信息
kubectl get pods -o yaml               # YAML 格式
kubectl get pods -o json               # JSON 格式

# 查看多种资源
kubectl get pods,svc,deploy

# 按标签过滤
kubectl get pods -l app=nginx
kubectl get pods -l "app in (nginx, web)"

# 按字段过滤
kubectl get pods --field-selector=status.phase=Running
```

### 常用资源查看

```bash
# Deployment
kubectl get deployments
kubectl get deploy                     # 缩写

# Service
kubectl get services
kubectl get svc                        # 缩写

# ConfigMap
kubectl get configmaps
kubectl get cm                         # 缩写

# Secret
kubectl get secrets

# PV/PVC
kubectl get pv
kubectl get pvc

# Ingress
kubectl get ingress
kubectl get ing                        # 缩写

# Namespace
kubectl get namespaces
kubectl get ns                         # 缩写
```

### 自定义输出

```bash
# 只输出名称
kubectl get pods -o name

# 自定义列
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase

# 排序
kubectl get pods --sort-by=.metadata.creationTimestamp

# 实时监控
kubectl get pods -w
kubectl get pods --watch
```

---

## 🔍 详细查看（DESCRIBE）

```bash
# 查看 Pod 详情（含 Events）
kubectl describe pod <pod-name>

# 查看 Deployment 详情
kubectl describe deployment <deploy-name>

# 查看 Service 详情
kubectl describe service <svc-name>

# 查看 Node 详情
kubectl describe node <node-name>
```

---

## 📝 创建资源

### 命令式创建

```bash
# 创建 Deployment
kubectl create deployment nginx --image=nginx:alpine

# 创建 Service
kubectl create service clusterip nginx --tcp=80:80

# 创建 ConfigMap
kubectl create configmap myconfig --from-literal=key1=value1

# 创建 Secret
kubectl create secret generic mysecret --from-literal=password=123456

# 创建 Namespace
kubectl create namespace dev
```

### 声明式创建

```bash
# 应用 YAML
kubectl apply -f deployment.yaml

# 应用目录下所有 YAML
kubectl apply -f ./manifests/

# 应用远程 YAML
kubectl apply -f https://example.com/deployment.yaml

# 递归应用
kubectl apply -f ./manifests/ -R
```

### 快速生成 YAML 模板

```bash
# 不真正创建，只输出 YAML
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml > deployment.yaml

# 将现有资源导出为 YAML
kubectl get deployment nginx -o yaml > deployment.yaml
```

---

## ✏️ 编辑与更新

### 直接编辑

```bash
# 编辑资源（打开编辑器）
kubectl edit deployment nginx

# 指定编辑器
KUBE_EDITOR="code --wait" kubectl edit deployment nginx
```

### 更新操作

```bash
# 更新镜像
kubectl set image deployment/nginx nginx=nginx:1.19

# 更新资源配置
kubectl set resources deployment/nginx -c=nginx --limits=cpu=200m,memory=512Mi

# 添加/更新标签
kubectl label pod nginx version=v1
kubectl label pod nginx version=v2 --overwrite

# 添加/更新注解
kubectl annotate pod nginx description="My nginx pod"
```

### 滚动更新

```bash
# 查看更新状态
kubectl rollout status deployment/nginx

# 查看更新历史
kubectl rollout history deployment/nginx

# 回滚到上一版本
kubectl rollout undo deployment/nginx

# 回滚到指定版本
kubectl rollout undo deployment/nginx --to-revision=2

# 暂停/恢复更新
kubectl rollout pause deployment/nginx
kubectl rollout resume deployment/nginx

# 重启 Deployment（触发滚动更新）
kubectl rollout restart deployment/nginx
```

---

## 📈 扩缩容

```bash
# 扩缩容
kubectl scale deployment nginx --replicas=5

# 自动扩缩容（需要 metrics-server）
kubectl autoscale deployment nginx --min=2 --max=10 --cpu-percent=80
```

---

## 🗑️ 删除资源

```bash
# 删除 Pod
kubectl delete pod nginx

# 删除 Deployment
kubectl delete deployment nginx

# 通过 YAML 删除
kubectl delete -f deployment.yaml

# 删除命名空间下所有资源
kubectl delete all --all -n dev

# 强制删除（卡住的 Pod）
kubectl delete pod nginx --force --grace-period=0

# 按标签删除
kubectl delete pods -l app=nginx
```

---

## 🔧 调试命令

### 日志

```bash
# 查看日志
kubectl logs nginx

# 实时日志
kubectl logs -f nginx

# 查看容器日志（多容器 Pod）
kubectl logs nginx -c sidecar

# 查看上一个容器的日志（崩溃后）
kubectl logs nginx --previous

# 查看最近 100 行
kubectl logs nginx --tail=100

# 查看最近 1 小时的日志
kubectl logs nginx --since=1h
```

### 执行命令

```bash
# 在 Pod 中执行命令
kubectl exec nginx -- ls /

# 进入 Pod 交互式 Shell
kubectl exec -it nginx -- /bin/sh
kubectl exec -it nginx -- /bin/bash

# 多容器 Pod 指定容器
kubectl exec -it nginx -c sidecar -- /bin/sh
```

### 端口转发

```bash
# 转发 Pod 端口
kubectl port-forward pod/nginx 8080:80

# 转发 Service 端口
kubectl port-forward svc/nginx 8080:80

# 转发 Deployment 端口
kubectl port-forward deployment/nginx 8080:80

# 后台运行
kubectl port-forward pod/nginx 8080:80 &
```

### 复制文件

```bash
# 从 Pod 复制到本地
kubectl cp nginx:/var/log/nginx/access.log ./access.log

# 从本地复制到 Pod
kubectl cp ./config.json nginx:/etc/config/
```

### 资源使用情况

```bash
# 节点资源使用（需要 metrics-server）
kubectl top nodes

# Pod 资源使用
kubectl top pods

# 按命名空间查看
kubectl top pods -n kube-system
```

---

## 🏷️ 标签与注解

```bash
# 添加标签
kubectl label pod nginx tier=frontend

# 修改标签
kubectl label pod nginx tier=backend --overwrite

# 删除标签
kubectl label pod nginx tier-

# 添加注解
kubectl annotate pod nginx description="Main web server"

# 按标签查看
kubectl get pods --show-labels
kubectl get pods -l tier=frontend
```

---

## 🌐 网络调试

```bash
# 创建临时 Pod 调试网络
kubectl run debug --image=busybox:1.36 --rm -it -- sh

# 常用命令（在 debug Pod 内）
nslookup kubernetes
wget -qO- http://nginx-service
ping 10.96.0.1

# 使用 curl 镜像
kubectl run curl --image=curlimages/curl --rm -it -- sh
curl http://nginx-service
```

---

## 📊 快速参考表

### kubectl 资源缩写

| 全名                   | 缩写   |
| :--------------------- | :----- |
| pods                   | po     |
| services               | svc    |
| deployments            | deploy |
| replicasets            | rs     |
| configmaps             | cm     |
| secrets                | (无)   |
| persistentvolumes      | pv     |
| persistentvolumeclaims | pvc    |
| namespaces             | ns     |
| ingresses              | ing    |
| nodes                  | no     |
| events                 | ev     |

### 常用选项

| 选项               | 作用         |
| :----------------- | :----------- |
| `-n`               | 指定命名空间 |
| `-A`               | 所有命名空间 |
| `-o wide`          | 详细输出     |
| `-o yaml`          | YAML 格式    |
| `-o json`          | JSON 格式    |
| `-w`               | 实时监控     |
| `-l`               | 标签选择器   |
| `--dry-run=client` | 模拟运行     |
| `--force`          | 强制操作     |

---

## 🎯 高效技巧

### 设置别名

添加到 `~/.zshrc` 或 `~/.bashrc`：

```bash
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deploy'
alias kga='kubectl get all'
alias kd='kubectl describe'
alias kl='kubectl logs'
alias ke='kubectl exec -it'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'
```

### 自动补全

```bash
# Zsh
source <(kubectl completion zsh)
echo 'source <(kubectl completion zsh)' >> ~/.zshrc

# Bash
source <(kubectl completion bash)
echo 'source <(kubectl completion bash)' >> ~/.bashrc
```

---

## 🎉 学习完成

恭喜你完成了整个 Kubernetes 学习课程！

**回顾一下你学到了什么**：

1. ✅ K8s 架构和核心组件
2. ✅ Pod 的创建和管理
3. ✅ Deployment 的副本管理和滚动更新
4. ✅ Service 的网络通信和服务发现
5. ✅ ConfigMap 和 Secret 的配置管理
6. ✅ PV 和 PVC 的存储管理
7. ✅ 日常运维必备命令

**下一步学习建议**：

- Helm：K8s 的包管理器
- RBAC：权限管理
- 监控：Prometheus + Grafana
- 日志：ELK / Loki
- 服务网格：Istio

祝你在 Kubernetes 的学习之路上越走越远！🚀
