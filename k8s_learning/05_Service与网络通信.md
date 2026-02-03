# 05. Service 与网络通信

> 本章目标：掌握 K8s 网络模型，理解 Service 的工作原理

---

## 1. 为什么需要 Service？

### 问题 1：Pod IP 不固定

每次 Pod 重建，IP 都会变：

```
nginx-abc12: 10.244.0.5  → 被删除
nginx-def34: 10.244.0.9  ← 新建，IP 变了
```

前端代码写死 IP？每次都要改？不现实。

### 问题 2：多个 Pod 如何负载均衡？

Deployment 有 3 个副本，访问哪一个？手动轮询？太傻了。

### 解决方案：Service

**Service** 提供：

1. **稳定的虚拟 IP**（ClusterIP）
2. **自动负载均衡**
3. **服务发现**（DNS）

---

## 2. Service 工作原理

```
                    ┌─────────────────────────────────┐
                    │          Service                │
                    │    ClusterIP: 10.96.0.100       │
                    │    Port: 80                     │
                    └───────────────┬─────────────────┘
                                    │
                         kube-proxy 负载均衡
                                    │
              ┌─────────────────────┼─────────────────────┐
              ▼                     ▼                     ▼
        ┌──────────┐          ┌──────────┐          ┌──────────┐
        │  Pod 1   │          │  Pod 2   │          │  Pod 3   │
        │ 10.244.0.5│          │ 10.244.0.6│          │ 10.244.0.7│
        └──────────┘          └──────────┘          └──────────┘
```

**匹配机制**：Service 通过 **标签选择器（selector）** 找到对应的 Pod

---

## 3. Service 的四种类型

| 类型             | 作用             | 访问范围         |
| :--------------- | :--------------- | :--------------- |
| **ClusterIP**    | 集群内部访问     | 仅集群内         |
| **NodePort**     | 通过节点端口访问 | 集群内 + 节点 IP |
| **LoadBalancer** | 云厂商负载均衡器 | 公网             |
| **ExternalName** | 外部域名映射     | -                |

---

## 4. ClusterIP（默认类型）

### 4.1 创建 Service

首先确保有一个 Deployment：

```bash
kubectl create deployment web --image=nginx:alpine --replicas=3
```

创建 `service-clusterip.yaml`：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-service
spec:
  type: ClusterIP # 默认类型，可省略
  selector:
    app: web # 匹配 Pod 的标签
  ports:
    - protocol: TCP
      port: 80 # Service 端口
      targetPort: 80 # Pod 端口
```

```bash
kubectl apply -f service-clusterip.yaml
```

### 4.2 验证

```bash
# 查看 Service
kubectl get svc
# NAME          TYPE        CLUSTER-IP     PORT(S)   AGE
# web-service   ClusterIP   10.96.0.100    80/TCP    30s

# 查看 Endpoints（Service 发现的 Pod 列表）
kubectl get endpoints web-service
# NAME          ENDPOINTS                                      AGE
# web-service   10.244.0.5:80,10.244.0.6:80,10.244.0.7:80      30s
```

### 4.3 集群内访问测试

```bash
# 创建一个临时 Pod 来测试
kubectl run test --image=busybox:1.36 --rm -it -- sh

# 在临时 Pod 内执行
wget -qO- http://web-service
# 或
wget -qO- http://10.96.0.100

# 多执行几次，观察负载均衡效果
```

---

## 5. NodePort

### 5.1 创建 NodePort Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-nodeport
spec:
  type: NodePort
  selector:
    app: web
  ports:
    - protocol: TCP
      port: 80 # Service 端口
      targetPort: 80 # Pod 端口
      nodePort: 30080 # 节点端口（30000-32767）
```

### 5.2 访问方式

```
外部访问：http://<节点IP>:30080
         http://192.168.49.2:30080
```

### 5.3 Minikube 特殊访问

```bash
# 让 Minikube 帮你打开浏览器
minikube service web-nodeport

# 只获取 URL
minikube service web-nodeport --url
```

---

## 6. LoadBalancer

在云环境（AWS、阿里云等）中，会自动创建云厂商的负载均衡器。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-lb
spec:
  type: LoadBalancer
  selector:
    app: web
  ports:
    - port: 80
      targetPort: 80
```

在 Minikube 中：

```bash
# 模拟 LoadBalancer
minikube tunnel
# 保持运行，另开终端操作
```

---

## 7. 服务发现（DNS）

K8s 内置 DNS 服务（CoreDNS），每个 Service 自动获得 DNS 记录。

### 7.1 DNS 格式

```
<service-name>.<namespace>.svc.cluster.local
```

### 7.2 访问示例

```bash
# 同命名空间内，只需要 Service 名
curl http://web-service

# 跨命名空间
curl http://web-service.default.svc.cluster.local

# 简写
curl http://web-service.default
```

### 7.3 测试 DNS 解析

```bash
kubectl run test --image=busybox:1.36 --rm -it -- sh

# 在临时 Pod 内
nslookup web-service
# Server:    10.96.0.10
# Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local
# Name:      web-service
# Address 1: 10.96.0.100 web-service.default.svc.cluster.local
```

---

## 8. Ingress（七层负载均衡）

Service 只能做四层（TCP/UDP）负载均衡。
Ingress 可以做七层（HTTP/HTTPS）负载均衡，支持：

- 基于域名路由
- 基于路径路由
- SSL 终止

### 8.1 启用 Ingress 控制器

```bash
minikube addons enable ingress
kubectl get pods -n ingress-nginx
```

### 8.2 创建 Ingress 规则

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: myapp.local # 域名
      http:
        paths:
          - path: / # 路径
            pathType: Prefix
            backend:
              service:
                name: web-service # 后端 Service
                port:
                  number: 80
```

### 8.3 配置本地 hosts

```bash
# 获取 Minikube IP
minikube ip
# 192.168.49.2

# 添加到 /etc/hosts
echo "192.168.49.2 myapp.local" | sudo tee -a /etc/hosts
```

### 8.4 访问测试

```bash
curl http://myapp.local
```

---

## 9. 多路径 Ingress 示例

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: multi-path-ingress
spec:
  ingressClassName: nginx
  rules:
    - host: shop.local
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## 10. 实战练习

### 练习 1：完整的 Deployment + Service + Ingress

```bash
# 1. 创建 Deployment
kubectl create deployment webapp --image=nginx:alpine --replicas=3

# 2. 暴露 Service
kubectl expose deployment webapp --port=80 --type=ClusterIP

# 3. 创建 Ingress
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: webapp-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: webapp.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: webapp
            port:
              number: 80
EOF

# 4. 配置 hosts
echo "$(minikube ip) webapp.local" | sudo tee -a /etc/hosts

# 5. 测试
curl http://webapp.local

# 6. 清理
kubectl delete ingress webapp-ingress
kubectl delete service webapp
kubectl delete deployment webapp
```

### 练习 2：使用 port-forward 快速调试

```bash
# 创建 Deployment
kubectl create deployment debug-app --image=nginx:alpine

# 不创建 Service，直接端口转发
kubectl port-forward deployment/debug-app 8080:80

# 另一个终端访问
curl http://localhost:8080

# 按 Ctrl+C 停止转发
```

---

## 11. 网络策略（Network Policy）

控制 Pod 之间的网络流量。

### 11.1 默认拒绝所有入站流量

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
spec:
  podSelector: {} # 匹配所有 Pod
  policyTypes:
    - Ingress
```

### 11.2 只允许特定 Pod 访问

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-from-frontend
spec:
  podSelector:
    matchLabels:
      app: api # 对 api Pod 生效
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend # 只允许 frontend 访问
```

---

## 12. 常见问题

### Service 没有 Endpoints

**排查**：

```bash
kubectl get endpoints <service-name>
# 如果是空的，说明没有匹配到 Pod
```

**原因**：

- Service 的 selector 和 Pod 的 labels 不匹配
- Pod 没有处于 Running 状态
- Pod 的 readinessProbe 失败

### 访问 Service 超时

**排查**：

```bash
# 检查 Pod 是否正常
kubectl get pods

# 检查 Service 配置
kubectl describe svc <service-name>

# 进入 Pod 测试端口是否监听
kubectl exec -it <pod-name> -- netstat -tlnp
```

---

## ✅ 本章检查点

- [ ] 理解 Service 解决的问题
- [ ] 能区分 ClusterIP、NodePort、LoadBalancer
- [ ] 能创建 Service 并通过 DNS 访问
- [ ] 能配置 Ingress 实现域名访问
- [ ] 能使用 port-forward 快速调试

---

## ⏭️ 下一章

应用需要配置文件和敏感信息，如何管理？

👉 [06_ConfigMap与Secret配置管理.md](06_ConfigMap与Secret配置管理.md)
