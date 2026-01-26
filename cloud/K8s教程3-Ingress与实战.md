# Kubernetes 教程3 - Ingress 与实战部署

## 1. Ingress 简介

### 1.1 为什么需要 Ingress？

**Service 的局限：**
- NodePort：端口有限（30000-32767），不够优雅
- LoadBalancer：每个服务都需要一个负载均衡器，成本高

**Ingress 的优势：**
- 统一入口，一个 IP 暴露多个服务
- 支持域名和路径路由
- 支持 HTTPS/TLS
- 支持负载均衡

```
                    Internet
                        |
                   [ Ingress ]
                   /    |    \
                  /     |     \
            Service1  Service2  Service3
               |        |         |
             Pods     Pods      Pods
```

### 1.2 Ingress 架构

```
Ingress 资源（规则定义）
        ↓
Ingress Controller（实际执行，如 Nginx）
        ↓
Service → Pod
```

---

## 2. 安装 Ingress Controller

### 2.1 Minikube 环境

```bash
# 启用 Ingress 插件
minikube addons enable ingress

# 验证
kubectl get pods -n ingress-nginx
```

### 2.2 使用 Helm 安装（通用）

```bash
# 添加仓库
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# 安装
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace

# 验证
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
```

---

## 3. Ingress 配置

### 3.1 基于路径的路由

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.local
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
      - path: /web
        pathType: Prefix
        backend:
          service:
            name: web-service
            port:
              number: 80
```

### 3.2 基于域名的路由

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: multi-host-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
  - host: web.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web-service
            port:
              number: 80
```

### 3.3 配置 HTTPS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-ingress
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - myapp.example.com
    secretName: tls-secret      # 包含证书的 Secret
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: my-service
            port:
              number: 80
```

**创建 TLS Secret：**
```bash
kubectl create secret tls tls-secret \
  --cert=path/to/cert.crt \
  --key=path/to/cert.key
```

### 3.4 常用注解

```yaml
metadata:
  annotations:
    # 重写路径
    nginx.ingress.kubernetes.io/rewrite-target: /
    
    # SSL 重定向
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    
    # 代理超时
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    
    # 限流
    nginx.ingress.kubernetes.io/limit-rps: "100"
    
    # 跨域
    nginx.ingress.kubernetes.io/enable-cors: "true"
```

---

## 4. 实战：部署完整的微服务应用

### 4.1 项目架构

```
                    Ingress
                   /       \
                  /         \
           /api/*          /web/*
              |               |
         api-service     web-service
              |               |
          api-pods        web-pods
              |
          mysql-service
              |
          mysql-pod
```

### 4.2 完整配置文件

**1. namespace.yaml**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: microservice
```

**2. mysql.yaml**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
  namespace: microservice
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: Secret
metadata:
  name: mysql-secret
  namespace: microservice
type: Opaque
stringData:
  MYSQL_ROOT_PASSWORD: rootpassword
  MYSQL_DATABASE: myapp
  MYSQL_USER: appuser
  MYSQL_PASSWORD: apppassword
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: microservice
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        ports:
        - containerPort: 3306
        envFrom:
        - secretRef:
            name: mysql-secret
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: mysql-data
        persistentVolumeClaim:
          claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql
  namespace: microservice
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
  type: ClusterIP
```

**3. api.yaml**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
  namespace: microservice
data:
  DB_HOST: mysql
  DB_PORT: "3306"
  DB_NAME: myapp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: microservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: nginx:alpine    # 替换为你的 API 镜像
        ports:
        - containerPort: 80
        envFrom:
        - configMapRef:
            name: api-config
        env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: MYSQL_USER
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: MYSQL_PASSWORD
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: api
  namespace: microservice
spec:
  selector:
    app: api
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
```

**4. web.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  namespace: microservice
spec:
  replicas: 2
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: web
        image: nginx:alpine    # 替换为你的前端镜像
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "32Mi"
            cpu: "50m"
          limits:
            memory: "64Mi"
            cpu: "100m"
---
apiVersion: v1
kind: Service
metadata:
  name: web
  namespace: microservice
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
```

**5. ingress.yaml**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  namespace: microservice
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.local
    http:
      paths:
      - path: /api(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: api
            port:
              number: 80
      - path: /(.*)
        pathType: Prefix
        backend:
          service:
            name: web
            port:
              number: 80
```

### 4.3 部署步骤

```bash
# 1. 创建命名空间
kubectl apply -f namespace.yaml

# 2. 部署 MySQL
kubectl apply -f mysql.yaml

# 3. 等待 MySQL 就绪
kubectl wait --for=condition=ready pod -l app=mysql -n microservice --timeout=120s

# 4. 部署 API 和 Web
kubectl apply -f api.yaml
kubectl apply -f web.yaml

# 5. 部署 Ingress
kubectl apply -f ingress.yaml

# 6. 查看所有资源
kubectl get all -n microservice

# 7. 查看 Ingress
kubectl get ingress -n microservice
```

### 4.4 本地测试

```bash
# 添加 hosts 记录
echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts

# Minikube 获取 IP
minikube ip

# 或者使用端口转发
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8080:80

# 测试访问
curl http://myapp.local:8080/        # Web
curl http://myapp.local:8080/api/    # API
```

---

## 5. Helm - K8s 包管理器

### 5.1 安装 Helm

```bash
# macOS
brew install helm

# 验证
helm version
```

### 5.2 基本使用

```bash
# 添加仓库
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# 搜索 Chart
helm search repo mysql

# 安装应用
helm install my-mysql bitnami/mysql

# 查看已安装
helm list

# 卸载
helm uninstall my-mysql

# 查看 Chart 的默认值
helm show values bitnami/mysql

# 自定义安装
helm install my-mysql bitnami/mysql \
  --set auth.rootPassword=mypassword \
  --set auth.database=mydb
```

### 5.3 使用 values.yaml

```yaml
# values.yaml
auth:
  rootPassword: mypassword
  database: mydb
  username: myuser
  password: myuserpassword

primary:
  persistence:
    size: 10Gi
```

```bash
helm install my-mysql bitnami/mysql -f values.yaml
```

---

## 6. 练习题

### 练习1：创建 Ingress

创建一个 Ingress，将以下路径路由到对应服务：
- `/api` → api-service:80
- `/` → web-service:80

<details>
<summary>点击查看答案</summary>

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web-service
            port:
              number: 80
```
</details>

---

### 练习2：多域名 Ingress

创建一个 Ingress，支持两个域名：
- api.test.com → api-service
- web.test.com → web-service

<details>
<summary>点击查看答案</summary>

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: multi-domain-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: api.test.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
  - host: web.test.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web-service
            port:
              number: 80
```
</details>

---

### 练习3：使用 Helm 部署 Redis

使用 Helm 部署一个 Redis 实例：
- 设置密码为 `redis123`
- 禁用持久化（测试用）

<details>
<summary>点击查看答案</summary>

```bash
# 添加仓库
helm repo add bitnami https://charts.bitnami.com/bitnami

# 安装 Redis
helm install my-redis bitnami/redis \
  --set auth.password=redis123 \
  --set master.persistence.enabled=false \
  --set replica.persistence.enabled=false

# 验证
kubectl get pods

# 测试连接
kubectl run redis-client --rm -it --image=redis -- redis-cli -h my-redis-master -a redis123
```
</details>

---

### 练习4：完整部署练习

部署一个完整的应用栈：
1. 创建命名空间 `practice`
2. 部署 Redis（使用 Helm）
3. 部署一个 nginx 作为 Web 服务（3副本）
4. 创建 Ingress 暴露 Web 服务

<details>
<summary>点击查看答案</summary>

```bash
# 1. 创建命名空间
kubectl create namespace practice

# 2. 部署 Redis
helm install redis bitnami/redis \
  -n practice \
  --set auth.password=redis123 \
  --set master.persistence.enabled=false \
  --set replica.persistence.enabled=false

# 3. 部署 Web 服务
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  namespace: practice
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
        image: nginx:alpine
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: web
  namespace: practice
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
EOF

# 4. 创建 Ingress
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web-ingress
  namespace: practice
spec:
  ingressClassName: nginx
  rules:
  - host: practice.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web
            port:
              number: 80
EOF

# 验证
kubectl get all -n practice
kubectl get ingress -n practice
```
</details>

---

### 练习5：扩缩容和更新

1. 将 Web 服务扩展到 5 个副本
2. 更新 nginx 镜像到 `nginx:1.25`
3. 查看滚动更新状态
4. 如果有问题，回滚到上一版本

<details>
<summary>点击查看答案</summary>

```bash
# 1. 扩展到 5 副本
kubectl scale deployment web -n practice --replicas=5

# 2. 更新镜像
kubectl set image deployment/web nginx=nginx:1.25 -n practice

# 3. 查看更新状态
kubectl rollout status deployment/web -n practice

# 4. 回滚（如果需要）
kubectl rollout undo deployment/web -n practice
```
</details>

---

## 7. 生产环境最佳实践

### 7.1 资源配置

```yaml
resources:
  requests:           # 最小资源需求
    memory: "128Mi"
    cpu: "100m"
  limits:             # 最大资源限制
    memory: "256Mi"
    cpu: "200m"
```

### 7.2 健康检查

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### 7.3 Pod 反亲和性（分散部署）

```yaml
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app: my-app
        topologyKey: kubernetes.io/hostname
```

### 7.4 Pod 中断预算

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-app-pdb
spec:
  minAvailable: 2      # 至少保持 2 个 Pod 可用
  selector:
    matchLabels:
      app: my-app
```

---

## 8. 本章小结

**Ingress：**
- 统一入口，支持域名和路径路由
- 需要安装 Ingress Controller
- 支持 HTTPS/TLS

**Helm：**
- K8s 的包管理器
- 简化复杂应用的部署
- 支持自定义配置

**生产最佳实践：**
- 配置资源限制
- 配置健康检查
- 使用 Pod 反亲和性
- 配置 Pod 中断预算

**学习路线建议：**
1. ✅ Docker 基础
2. ✅ Dockerfile
3. ✅ Docker Compose
4. ✅ K8s 基础概念
5. ✅ K8s 核心资源
6. ✅ Ingress 和 Helm
7. 下一步：CI/CD（Jenkins/GitLab CI/ArgoCD）
8. 进阶：服务网格（Istio）、监控（Prometheus）
