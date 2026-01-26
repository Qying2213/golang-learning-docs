# Docker & Kubernetes 教程目录

## 📚 学习顺序和重要程度

### Docker 部分

| 教程 | 重要程度 | 学习时间 | 说明 |
|------|---------|---------|------|
| **Docker核心概念-镜像与容器.md** | ⭐⭐⭐ | 30分钟 | 必读！理解镜像和容器的区别 |
| **Docker教程1-基础入门.md** | ⭐⭐⭐ | 2-3小时 | 核心基础，必须全部掌握 |
| **Docker教程2-Dockerfile详解.md** | ⭐⭐⭐ | 2-3小时 | 构建镜像必备 |
| **Docker教程3-Compose编排.md** | ⭐⭐⭐ | 2-3小时 | 多容器应用必备 |

### Kubernetes 部分

| 教程 | 重要程度 | 学习时间 | 说明 |
|------|---------|---------|------|
| **K8s教程1-基础概念.md** | ⭐⭐⭐ | 2-3小时 | K8s 入门必读 |
| **K8s教程2-核心资源.md** | ⭐⭐⭐ | 3-4小时 | 最重要的章节 |
| **K8s教程3-Ingress与实战.md** | ⭐⭐ | 2-3小时 | 实战部署 |

### 辅助文档

| 文档 | 重要程度 | 说明 |
|------|---------|------|
| **学习路线图.md** | ⭐⭐⭐ | 先看这个！规划学习路径 |

---

## 📖 各教程重点内容标注

### Docker核心概念-镜像与容器.md ⭐⭐⭐

**必须理解：**
- ⭐⭐⭐ 镜像 vs 容器的区别
- ⭐⭐⭐ Docker 架构（Client-Server）
- ⭐⭐⭐ 镜像分层结构
- ⭐⭐⭐ 容器生命周期
- ⭐⭐⭐ Docker 启动流程

**可以跳过：**
- ⭐ 存储位置细节

---

### Docker教程1-基础入门.md ⭐⭐⭐

**必须掌握：**
- ⭐⭐⭐ 镜像操作：pull、images、rmi
- ⭐⭐⭐ 容器操作：run、ps、stop、rm
- ⭐⭐⭐ docker run 的常用选项：-d、-p、--name、-e、-v
- ⭐⭐⭐ 进入容器：exec -it
- ⭐⭐⭐ 查看日志：logs -f
- ⭐⭐⭐ 数据持久化：-v 数据卷

**了解即可：**
- ⭐ docker search（很少用）
- ⭐⭐ docker inspect（偶尔用）

**必做练习：**
- ⭐⭐⭐ 练习1-5：运行 nginx、MySQL、Redis

---

### Docker教程2-Dockerfile详解.md ⭐⭐⭐

**必须掌握：**
- ⭐⭐⭐ FROM：基础镜像
- ⭐⭐⭐ WORKDIR：工作目录
- ⭐⭐⭐ COPY：复制文件
- ⭐⭐⭐ RUN：执行命令
- ⭐⭐⭐ ENV：环境变量
- ⭐⭐⭐ EXPOSE：声明端口
- ⭐⭐⭐ CMD：启动命令
- ⭐⭐⭐ 多阶段构建（减小镜像体积）
- ⭐⭐⭐ .dockerignore

**了解即可：**
- ⭐⭐ ENTRYPOINT（不常改）
- ⭐ ARG（偶尔用）
- ⭐ ADD（用 COPY 就够了）
- ⭐ VOLUME（一般用 -v）

**必做练习：**
- ⭐⭐⭐ 练习1-5：为 Go/Python/Node 项目写 Dockerfile

---

### Docker教程3-Compose编排.md ⭐⭐⭐

**必须掌握：**
- ⭐⭐⭐ docker-compose.yml 基本结构
- ⭐⭐⭐ services、volumes、networks 配置
- ⭐⭐⭐ 常用命令：up、down、ps、logs
- ⭐⭐⭐ 环境变量配置：environment、env_file
- ⭐⭐⭐ 依赖关系：depends_on

**了解即可：**
- ⭐⭐ 资源限制：deploy.resources
- ⭐⭐ 健康检查：healthcheck
- ⭐ 网络详解（会用默认网络就行）

**必做练习：**
- ⭐⭐⭐ 练习1-5：部署 Go + MySQL + Redis

---

### K8s教程1-基础概念.md ⭐⭐⭐

**必须理解：**
- ⭐⭐⭐ K8s 架构（Control Plane + Node）
- ⭐⭐⭐ Pod：最小部署单元
- ⭐⭐⭐ Deployment：管理 Pod 副本
- ⭐⭐⭐ Service：服务发现和负载均衡
- ⭐⭐⭐ Namespace：资源隔离

**不用深究：**
- ⭐ ReplicaSet（Deployment 会自动管理）
- ⭐ etcd 细节
- ⭐ kube-proxy 原理

**必须掌握的命令：**
- ⭐⭐⭐ kubectl get pods/deploy/svc
- ⭐⭐⭐ kubectl describe
- ⭐⭐⭐ kubectl apply -f
- ⭐⭐⭐ kubectl delete
- ⭐⭐⭐ kubectl logs
- ⭐⭐⭐ kubectl exec -it
- ⭐⭐⭐ kubectl port-forward

**必做练习：**
- ⭐⭐⭐ 练习1-6：部署 nginx、查看日志、扩缩容

---

### K8s教程2-核心资源.md ⭐⭐⭐

**必须掌握：**
- ⭐⭐⭐ Deployment YAML 配置
- ⭐⭐⭐ Service 类型：ClusterIP、NodePort、LoadBalancer
- ⭐⭐⭐ ConfigMap：配置管理
- ⭐⭐⭐ Secret：密钥管理
- ⭐⭐⭐ 健康检查：livenessProbe、readinessProbe
- ⭐⭐⭐ 滚动更新和回滚

**重要但不急：**
- ⭐⭐ PersistentVolume：持久化存储
- ⭐⭐ 资源限制：requests、limits

**了解即可：**
- ⭐ startupProbe（启动探针）
- ⭐ 访问模式细节

**必做练习：**
- ⭐⭐⭐ 练习1-8：完整部署应用（Deployment + Service + ConfigMap + Secret）

---

### K8s教程3-Ingress与实战.md ⭐⭐

**必须掌握：**
- ⭐⭐⭐ Ingress 基本概念
- ⭐⭐⭐ 基于路径的路由
- ⭐⭐⭐ 基于域名的路由
- ⭐⭐ Helm 基本使用

**了解即可：**
- ⭐⭐ HTTPS/TLS 配置
- ⭐ Ingress 注解
- ⭐ Helm Chart 开发

**必做练习：**
- ⭐⭐⭐ 练习1-5：部署完整微服务（前端 + 后端 + 数据库 + Ingress）

---

## 🎯 学习建议

### 第一周：Docker 基础
1. **Day 1-2**：Docker核心概念 + 教程1（1-4章）
2. **Day 3-4**：教程1（5-7章）+ 练习题
3. **Day 5-6**：教程2 Dockerfile + 练习题
4. **Day 7**：教程3 Compose + 练习题

### 第二周：K8s 基础
1. **Day 1-2**：K8s教程1（1-4章）
2. **Day 3-4**：K8s教程1（5-7章）+ 练习题
3. **Day 5-6**：K8s教程2（1-5章）
4. **Day 7**：K8s教程2（6-8章）+ 练习题

### 第三周：K8s 进阶 + 实战
1. **Day 1-2**：K8s教程3 Ingress
2. **Day 3-4**：实战项目1：部署 Go Web 应用
3. **Day 5-7**：实战项目2：部署完整微服务

---

## ✅ 学习检验标准

### Docker 合格标准
- [ ] 能用一句话解释镜像和容器的区别
- [ ] 能独立运行 nginx、mysql、redis 容器
- [ ] 能为 Go 项目写 Dockerfile 并构建镜像
- [ ] 能用 Compose 部署多容器应用
- [ ] 理解数据持久化，会用 -v

### K8s 合格标准
- [ ] 能画出 Pod、Deployment、Service 的关系图
- [ ] 能写 Deployment 和 Service 的 YAML
- [ ] 能用 kubectl 部署、扩缩容、查看日志
- [ ] 能配置 ConfigMap 和 Secret
- [ ] 能部署一个完整的微服务应用

---

## 💡 学习技巧

1. **先会用再深究**：不要纠结原理，先把命令用起来
2. **边学边练**：每学一个知识点，立即动手实践
3. **重点突破**：把 ⭐⭐⭐ 的内容练熟，⭐ 的内容了解即可
4. **遇到问题先查日志**：`docker logs` 和 `kubectl logs` 是最好的老师
5. **做项目巩固**：学完基础后，一定要做实战项目

---

## 📝 快速查询

### Docker 常用命令速查
```bash
# 镜像
docker pull nginx
docker images
docker rmi nginx

# 容器
docker run -d -p 8080:80 --name web nginx
docker ps
docker logs -f web
docker exec -it web bash
docker stop web
docker rm web

# 数据卷
docker volume ls
docker run -v my-data:/data nginx

# Compose
docker compose up -d
docker compose down
docker compose logs -f
```

### K8s 常用命令速查
```bash
# 查看资源
kubectl get pods
kubectl get deploy
kubectl get svc
kubectl describe pod <name>

# 部署
kubectl apply -f app.yaml
kubectl delete -f app.yaml

# 调试
kubectl logs <pod>
kubectl logs -f <pod>
kubectl exec -it <pod> -- bash
kubectl port-forward <pod> 8080:80

# 扩缩容
kubectl scale deployment <name> --replicas=3

# 滚动更新
kubectl set image deployment/<name> app=app:v2
kubectl rollout status deployment/<name>
kubectl rollout undo deployment/<name>
```

---

## 🚀 下一步

学完这些教程后，可以继续学习：
1. CI/CD：Jenkins、GitLab CI、ArgoCD
2. 监控：Prometheus + Grafana
3. 日志：ELK/EFK
4. 服务网格：Istio

但这些都是进阶内容，先把基础打牢！
