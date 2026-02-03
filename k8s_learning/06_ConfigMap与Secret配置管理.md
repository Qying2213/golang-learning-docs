# 06. ConfigMap 与 Secret 配置管理

> 本章目标：掌握配置与敏感信息的分离管理

---

## 1. 为什么需要配置分离？

**问题**：应用配置硬编码在代码/镜像里

```dockerfile
# 不好的做法
ENV DATABASE_HOST=192.168.1.100
ENV DATABASE_PASSWORD=mypassword123
```

**问题**：

- 换环境要重新构建镜像
- 密码暴露在镜像里，不安全
- 修改配置需要重新部署

**解决**：使用 ConfigMap 和 Secret

---

## 2. ConfigMap 与 Secret 对比

| 特性   | ConfigMap          | Secret              |
| :----- | :----------------- | :------------------ |
| 用途   | 非敏感配置         | 敏感信息            |
| 存储   | 明文               | Base64 编码         |
| 示例   | 日志级别、服务地址 | 密码、API Key、证书 |
| 安全性 | 低                 | 较高（可加密存储）  |

---

## 3. ConfigMap

### 3.1 创建 ConfigMap

**方法 1：命令行创建**

```bash
# 从字面值创建
kubectl create configmap app-config \
  --from-literal=LOG_LEVEL=info \
  --from-literal=APP_ENV=production

# 从文件创建
kubectl create configmap nginx-config \
  --from-file=nginx.conf

# 从目录创建（每个文件变成一个 key）
kubectl create configmap configs \
  --from-file=./config-dir/
```

**方法 2：YAML 创建**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  # 简单键值对
  LOG_LEVEL: "info"
  APP_ENV: "production"
  DATABASE_HOST: "mysql.default.svc.cluster.local"

  # 多行配置文件
  nginx.conf: |
    server {
      listen 80;
      server_name localhost;
      location / {
        root /usr/share/nginx/html;
      }
    }

  # JSON 配置
  config.json: |
    {
      "debug": false,
      "port": 8080
    }
```

### 3.2 使用 ConfigMap

**方式 1：环境变量注入（单个值）**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
    - name: app
      image: nginx:alpine
      env:
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: LOG_LEVEL
```

**方式 2：环境变量注入（全部值）**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
    - name: app
      image: nginx:alpine
      envFrom:
        - configMapRef:
            name: app-config
```

**方式 3：挂载为文件**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
    - name: app
      image: nginx:alpine
      volumeMounts:
        - name: config-volume
          mountPath: /etc/config # 挂载目录
          readOnly: true
  volumes:
    - name: config-volume
      configMap:
        name: app-config
```

挂载后，`/etc/config/` 目录下会有：

- `LOG_LEVEL` 文件，内容是 `info`
- `nginx.conf` 文件，内容是配置内容
- ...

**方式 4：挂载单个文件**

```yaml
volumeMounts:
  - name: config-volume
    mountPath: /etc/nginx/nginx.conf
    subPath: nginx.conf # 只挂载这一个文件
```

---

## 4. Secret

### 4.1 创建 Secret

**方法 1：命令行创建**

```bash
# 从字面值创建
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=S3cr3t!

# 从文件创建
kubectl create secret generic tls-secret \
  --from-file=cert.pem \
  --from-file=key.pem
```

**方法 2：YAML 创建（需要 Base64 编码）**

```bash
# 先编码
echo -n 'admin' | base64
# YWRtaW4=

echo -n 'S3cr3t!' | base64
# UzNjcjN0IQ==
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  username: YWRtaW4= # base64 编码
  password: UzNjcjN0IQ==


# 或者使用 stringData（自动编码）
---
apiVersion: v1
kind: Secret
metadata:
  name: db-secret-v2
type: Opaque
stringData: # 明文，K8s 自动编码
  username: admin
  password: S3cr3t!
```

### 4.2 Secret 类型

| 类型                             | 用途             |
| :------------------------------- | :--------------- |
| `Opaque`                         | 通用密钥（默认） |
| `kubernetes.io/tls`              | TLS 证书         |
| `kubernetes.io/dockerconfigjson` | Docker 仓库认证  |
| `kubernetes.io/basic-auth`       | 基本认证         |

### 4.3 使用 Secret

**环境变量注入**：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
    - name: app
      image: myapp:v1
      env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: username
        - name: DB_PASS
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
```

**挂载为文件**：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
    - name: app
      image: myapp:v1
      volumeMounts:
        - name: secret-volume
          mountPath: /etc/secrets
          readOnly: true
  volumes:
    - name: secret-volume
      secret:
        secretName: db-secret
```

---

## 5. 实战示例

### 5.1 完整的 Web 应用配置

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp-config
data:
  APP_ENV: production
  LOG_LEVEL: info
  API_URL: http://api.default.svc.cluster.local:8080

---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: webapp-secret
type: Opaque
stringData:
  DB_PASSWORD: SuperSecret123
  API_KEY: abcd-1234-efgh-5678

---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
    spec:
      containers:
        - name: webapp
          image: nginx:alpine
          envFrom:
            - configMapRef:
                name: webapp-config
          env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: webapp-secret
                  key: DB_PASSWORD
            - name: API_KEY
              valueFrom:
                secretKeyRef:
                  name: webapp-secret
                  key: API_KEY
```

### 5.2 应用并验证

```bash
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml

# 进入 Pod 验证环境变量
kubectl exec -it deployment/webapp -- env | grep -E "(APP_|DB_|API_|LOG_)"
```

---

## 6. 热更新

### 6.1 环境变量方式

**不支持热更新**：必须重启 Pod 才能生效

```bash
# 更新 ConfigMap
kubectl edit configmap webapp-config

# 重启 Deployment
kubectl rollout restart deployment webapp
```

### 6.2 Volume 挂载方式

**支持热更新**：大约 1 分钟内自动生效

```bash
# 更新 ConfigMap
kubectl edit configmap webapp-config

# 无需重启，等待约 1 分钟后检查
kubectl exec -it <pod> -- cat /etc/config/LOG_LEVEL
```

---

## 7. 最佳实践

### 7.1 不要把 Secret 提交到 Git

```bash
# .gitignore
*-secret.yaml
secrets/
```

### 7.2 使用外部 Secret 管理

生产环境建议使用：

- HashiCorp Vault
- AWS Secrets Manager
- 阿里云 KMS

### 7.3 设置不可变 ConfigMap/Secret

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: immutable-config
immutable: true # 创建后不可修改
data:
  key: value
```

好处：

- 防止意外修改导致故障
- 提高性能（K8s 不用 watch 变化）

---

## 8. 实战练习

```bash
# 1. 创建 ConfigMap
kubectl create configmap test-config \
  --from-literal=ENV=dev \
  --from-literal=DEBUG=true

# 2. 创建 Secret
kubectl create secret generic test-secret \
  --from-literal=PASSWORD=test123

# 3. 创建使用它们的 Pod
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: test
    image: busybox:1.36
    command: ['sleep', '3600']
    envFrom:
    - configMapRef:
        name: test-config
    env:
    - name: PASSWORD
      valueFrom:
        secretKeyRef:
          name: test-secret
          key: PASSWORD
EOF

# 4. 验证
kubectl exec test-pod -- env

# 5. 清理
kubectl delete pod test-pod
kubectl delete configmap test-config
kubectl delete secret test-secret
```

---

## ✅ 本章检查点

- [ ] 理解 ConfigMap 和 Secret 的区别
- [ ] 能创建 ConfigMap 和 Secret
- [ ] 能通过环境变量注入配置
- [ ] 能通过 Volume 挂载配置文件
- [ ] 理解热更新的限制

---

## ⏭️ 下一章

应用需要持久化数据，如何管理存储？

👉 [07_PV与PVC存储管理.md](07_PV与PVC存储管理.md)
