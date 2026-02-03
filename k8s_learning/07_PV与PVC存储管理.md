# 07. PV 与 PVC 存储管理

> 本章目标：掌握 K8s 存储体系，理解 PV/PVC 的工作原理

---

## 1. 为什么需要持久化存储？

**问题**：Pod 是临时的，数据随时可能丢失

```
Pod 被删除 → 容器内的数据全部丢失
Pod 被调度到其他节点 → 无法访问原节点的数据
```

**解决**：使用持久化存储（Persistent Storage）

---

## 2. 存储体系概览

```
┌─────────────────────────────────────────────────────────────┐
│                        应用层                                │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                       Pod                               │ │
│  │    volumes:                                             │ │
│  │    - name: data                                         │ │
│  │      persistentVolumeClaim:                             │ │
│  │        claimName: mysql-pvc   ◀─────── 申领 (PVC)        │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │ 绑定
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       存储层                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              PersistentVolume (PV)                      │ │
│  │    capacity: 10Gi                                       │ │
│  │    storageClassName: fast                               │ │
│  │    hostPath: /data/mysql     ◀─────── 实际存储位置        │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**三个核心概念**：

- **Volume**：Pod 级别的临时存储
- **PersistentVolume (PV)**：集群级别的存储资源
- **PersistentVolumeClaim (PVC)**：用户对存储的申请

---

## 3. Volume 类型

### 3.1 emptyDir（临时目录）

Pod 内的容器共享，Pod 删除则数据丢失

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: shared-volume
spec:
  containers:
    - name: writer
      image: busybox
      command: ["sh", "-c", "echo hello > /data/message; sleep 3600"]
      volumeMounts:
        - name: shared
          mountPath: /data
    - name: reader
      image: busybox
      command: ["sh", "-c", "sleep 10; cat /data/message; sleep 3600"]
      volumeMounts:
        - name: shared
          mountPath: /data
  volumes:
    - name: shared
      emptyDir: {}
```

### 3.2 hostPath（节点目录）

挂载宿主机的目录，Pod 调度到其他节点数据就没了

```yaml
volumes:
  - name: host-data
    hostPath:
      path: /data/app
      type: DirectoryOrCreate # 不存在则创建
```

**使用场景**：DaemonSet（日志收集、监控代理）

### 3.3 configMap / secret

前一章已介绍，不再赘述

---

## 4. PersistentVolume (PV)

PV 是管理员预先创建的存储资源

### 4.1 创建 PV

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-hostpath
spec:
  capacity:
    storage: 10Gi # 容量
  accessModes:
    - ReadWriteOnce # 访问模式
  persistentVolumeReclaimPolicy: Retain # 回收策略
  storageClassName: manual
  hostPath: # 存储后端（仅用于学习）
    path: /data/pv-hostpath
```

### 4.2 访问模式（accessModes）

| 模式          | 缩写 | 含义       |
| :------------ | :--- | :--------- |
| ReadWriteOnce | RWO  | 单节点读写 |
| ReadOnlyMany  | ROX  | 多节点只读 |
| ReadWriteMany | RWX  | 多节点读写 |

不是所有存储后端都支持所有模式，例如：

- AWS EBS 只支持 RWO
- NFS 支持 RWX

### 4.3 回收策略（persistentVolumeReclaimPolicy）

| 策略    | 行为                              |
| :------ | :-------------------------------- |
| Retain  | PVC 删除后，PV 保留（需手动清理） |
| Delete  | PVC 删除后，PV 和底层存储一起删除 |
| Recycle | 清空数据后重新可用（已废弃）      |

---

## 5. PersistentVolumeClaim (PVC)

PVC 是用户对存储的申请

### 5.1 创建 PVC

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi # 申请 5Gi
  storageClassName: manual # 匹配 PV 的 storageClassName
```

### 5.2 PV 和 PVC 绑定

K8s 自动匹配：

1. `accessModes` 匹配
2. `storageClassName` 匹配
3. `storage` 容量满足（PV >= PVC）

```bash
kubectl get pv
# NAME           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM
# pv-hostpath    10Gi       RWO            Retain           Bound    default/mysql-pvc

kubectl get pvc
# NAME        STATUS   VOLUME         CAPACITY   ACCESS MODES   STORAGECLASS
# mysql-pvc   Bound    pv-hostpath    10Gi       RWO            manual
```

### 5.3 在 Pod 中使用 PVC

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mysql
spec:
  containers:
    - name: mysql
      image: mysql:8.0
      env:
        - name: MYSQL_ROOT_PASSWORD
          value: password
      volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
  volumes:
    - name: mysql-data
      persistentVolumeClaim:
        claimName: mysql-pvc # 引用 PVC
```

---

## 6. StorageClass（动态供给）

每次都要管理员手动创建 PV？太麻烦了。

**StorageClass** 可以动态创建 PV。

### 6.1 查看 StorageClass

```bash
kubectl get storageclass
# NAME                 PROVISIONER                RECLAIMPOLICY
# standard (default)   k8s.io/minikube-hostpath   Delete
```

Minikube 自带一个 `standard` StorageClass。

### 6.2 使用动态供给

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: dynamic-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: standard # 使用 StorageClass
```

**结果**：K8s 自动创建一个 PV 并绑定

```bash
kubectl get pv
# 自动创建的 PV
# pvc-xxxx   1Gi   RWO   Delete   Bound    default/dynamic-pvc   standard
```

---

## 7. 实战练习

### 练习 1：完整的有状态应用部署

```bash
# 1. 创建 PVC
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: standard
EOF

# 2. 创建 MySQL Deployment
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
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
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: rootpassword
        - name: MYSQL_DATABASE
          value: testdb
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-storage
        persistentVolumeClaim:
          claimName: mysql-data
EOF

# 3. 创建 Service
kubectl expose deployment mysql --port=3306

# 4. 验证
kubectl get pods
kubectl get pvc
kubectl get pv

# 5. 测试数据持久化
# 进入 MySQL Pod
kubectl exec -it deployment/mysql -- mysql -uroot -prootpassword

# 在 MySQL 中创建数据
CREATE TABLE testdb.users (id INT, name VARCHAR(50));
INSERT INTO testdb.users VALUES (1, 'Alice');
SELECT * FROM testdb.users;
EXIT;

# 删除 Pod（模拟故障）
kubectl delete pod -l app=mysql

# 等待新 Pod 创建
kubectl get pods -w

# 再次验证数据还在
kubectl exec -it deployment/mysql -- mysql -uroot -prootpassword -e "SELECT * FROM testdb.users"
# 数据还在！

# 6. 清理
kubectl delete deployment mysql
kubectl delete service mysql
kubectl delete pvc mysql-data
```

---

## 8. StatefulSet 简介

对于有状态应用（数据库、消息队列），推荐使用 **StatefulSet**

### 8.1 StatefulSet 特点

- Pod 有固定的网络标识（mysql-0, mysql-1, mysql-2）
- Pod 按顺序创建和删除
- 每个 Pod 可以有独立的 PVC

### 8.2 简单示例

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql
  replicas: 3
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
          volumeMounts:
            - name: data
              mountPath: /var/lib/mysql
  volumeClaimTemplates: # 每个 Pod 自动创建 PVC
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```

**结果**：

- Pod: mysql-0, mysql-1, mysql-2
- PVC: data-mysql-0, data-mysql-1, data-mysql-2

---

## 9. 常见问题

### PVC 一直 Pending

**排查**：

```bash
kubectl describe pvc <pvc-name>
```

**常见原因**：

- 没有匹配的 PV
- StorageClass 不存在
- 容量/访问模式不匹配

### Pod 挂载失败

**排查**：

```bash
kubectl describe pod <pod-name>
# 查看 Events 中的错误信息
```

**常见原因**：

- PVC 未绑定
- 多个 Pod 同时挂载 RWO 的 PV

---

## ✅ 本章检查点

- [ ] 理解 Volume、PV、PVC 的关系
- [ ] 能创建 PV 和 PVC
- [ ] 理解 StorageClass 动态供给
- [ ] 能部署使用持久化存储的应用
- [ ] 了解 StatefulSet 的特点

---

## ⏭️ 下一章

最后一章，汇总日常使用的 kubectl 命令。

👉 [08_kubectl常用命令速查表.md](08_kubectl常用命令速查表.md)
