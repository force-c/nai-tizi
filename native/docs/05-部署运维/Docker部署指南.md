# Docker 部署指南

## 构建镜像

在 `native` 目录执行：

```bash
docker build -f dockerfile/Dockerfile -t quick-admin-native:local .
```

镜像使用多阶段构建，运行阶段为非 root 用户，包含 `conf.prod.yaml`，监听 `9009`，健康检查访问 `/health/live`。

## Compose 启动

```bash
export QUICK_ADMIN_POSTGRES_PASSWORD='replace-with-a-strong-password'
export QUICK_ADMIN_JWT_SECRET='replace-with-at-least-32-random-characters'
docker compose -f dockerfile/docker-compose.yml up --build
```

Compose 启动 API、PostgreSQL 16 和 Redis 7。数据库与 JWT 密钥没有默认值，缺失时 Compose 会直接报错。

验证：

```bash
curl http://127.0.0.1:9009/health/live
curl http://127.0.0.1:9009/health/ready
```

生产环境不要暴露 PostgreSQL、Redis 端口，也不要使用 Compose 文件作为 Secret 管理方案。应由部署平台注入 `QUICK_ADMIN_*` 环境变量，并让 Nginx 同源代理 API。
