# Kubernetes 部署指南

`native/k8s/` 提供最小部署清单：

- `deployment.yaml`：API Deployment、资源限制、启动/存活/就绪探针
- `service.yaml`：ClusterIP Service
- `configmap.yaml`：非敏感配置
- `secret.yaml.example`：Secret 字段示例，不应直接用于生产

部署前将镜像地址替换为带不可变版本或 digest 的真实镜像，并通过集群 Secret 管理系统创建 `quick-admin-secret`：

```text
database-dsn
redis-password
jwt-secret
```

随后应用资源：

```bash
kubectl apply -f k8s/configmap.yaml
kubectl apply -f /secure/path/secret.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/deployment.yaml
```

容器端口为 `9009`。探针路径分别为 `/health/startup`、`/health/live`、`/health/ready`。Service 保持 ClusterIP，由 Ingress/Nginx 同源转发，不在 API 内启用生产 CORS。
