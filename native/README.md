# quick.admin native

`native` 是 quick.admin 的原生 Go 后端，也是三套后端实现的业务与 HTTP 契约基线。

## 本地启动

依赖 Go 1.25、PostgreSQL 16、Redis 7。开发环境的本地 PostgreSQL、Redis 和 JWT 参数直接配置在 `cmd/api/conf.dev.yaml`，确认本机服务与配置一致后即可启动：

```bash
cd native
go run ./cmd/api
```

`QUICK_ADMIN_APP_ENV` 只负责选择 `cmd/api/conf.dev.yaml` 或 `cmd/api/conf.prod.yaml`，默认 `dev`，因此本地启动无需 export。配置优先级是：代码默认值 < YAML < `QUICK_ADMIN_*` 环境变量；需要临时覆盖某项配置时仍可使用环境变量。配置结构中没有 `env` 字段，也不识别旧环境变量。

开发配置允许保存本机专用账号和开发密钥，但不得用于生产。`conf.prod.yaml` 中的密钥保持为空，生产环境必须由部署平台注入。

生产环境默认关闭 CORS，前后端通过 Nginx 同源代理；开发环境可在 `conf.dev.yaml` 中开启。

## 数据库迁移与初始管理员

服务启动时自动执行嵌入二进制的 Goose 迁移。数据库结构只能通过 `internal/database/migrations/sql/` 中的版本化 SQL 修改，不使用 GORM AutoMigrate。

迁移会创建逻辑客户端 `web-admin` 和内置角色，但不会写入默认管理员或默认密码。创建首个管理员：

```bash
export QUICK_ADMIN_USERMGR_PASSWORD='replace-with-a-strong-password'
go run ./cmd/usermgr --operation=create --username=admin --nickname=管理员 --role=super_admin
```

重置密码：

```bash
export QUICK_ADMIN_USERMGR_PASSWORD='replace-with-a-new-strong-password'
go run ./cmd/usermgr --operation=reset --username=admin
```

密码不支持命令行参数，工具也不会回显密码。

## 认证

登录接口为 `POST /login`，客户端 ID 为 `web-admin`，支持 `password`、`email`、`sms`、`wechat`。微信登录只接收小程序 `wxCode`，服务端通过微信接口换取 OpenID/UnionID，不接收客户端提交的 OpenID。

Access Token 与 Refresh Token 都关联服务端 Redis 会话；刷新令牌单次使用并在刷新后轮换，登出会立即撤销当前会话。

## 质量检查

```bash
make verify # go test、go vet、race、Swagger freshness
make ci     # verify + Docker build
```

涉及通用工具、配置、认证基础设施和中间件的改动必须补充单元测试。controller/service 中的具体业务逻辑不强制单测，可按风险补充集成或契约测试。

## Docker

```bash
export QUICK_ADMIN_POSTGRES_PASSWORD='replace-me'
export QUICK_ADMIN_JWT_SECRET='replace-with-at-least-32-random-characters'
docker compose -f dockerfile/docker-compose.yml up --build
```

HTTP 端口为 `9009`，探针为 `/health/startup`、`/health/live`、`/health/ready`。
