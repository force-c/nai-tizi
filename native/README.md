# quick.admin native

`native` 是 quick.admin 的原生 Go 后端，也是三套后端实现的业务与 HTTP 契约基线。

## 本地启动

依赖 Go 1.25 和 PostgreSQL 16；API 另外使用 Redis 7。首次拉取工程后，从各服务的 example 初始化本地配置：

```bash
cd native
make init-config
```

API 配置位于 `application/api/`，Scheduler 配置位于 `application/scheduler/`。修改对应服务的 `conf.dev.yaml` 后即可启动：

```bash
go run ./application/api
```

API 是无状态 Web 进程，不启动定时任务。Scheduler、定时任务和常驻后台协程由独立进程运行：

```bash
cd native
go run ./application/scheduler
```

生产环境可以水平扩容 API；Scheduler 默认运行一个副本。

`LIGHTNING_APP_ENV` 只负责在具体服务目录中选择 `conf.dev.yaml` 或 `conf.prod.yaml`，默认 `dev`，因此本地启动无需 export。配置优先级是：代码默认值 < YAML < `LIGHTNING_*` 环境变量；需要临时覆盖某项配置时仍可使用环境变量。配置结构中没有 `env` 字段，也不识别旧环境变量。

Git 只管理每个服务的 `conf.example.yaml` 和 `zaplogger.example.yaml`；实际的开发、生产配置均被忽略，`make init-config` 不会覆盖已经存在的文件。API example 包含 HTTP、数据库、Redis、JWT、CORS、RabbitMQ、S3 和 WebSocket；Scheduler example 只包含数据库。部署时仍可用 `LIGHTNING_*` 环境变量覆盖。

生产环境默认关闭 CORS，前后端通过 Nginx 同源代理；开发环境可在 `conf.dev.yaml` 中开启。

验证码、微信、短信、邮件和 Scheduler 属于运行期模块配置，不在 YAML 或环境变量中定义。它们持久化在 `s_config`；API 模块通过 Redis 共享读取，缓存缺失时在分布式锁内从数据库加载并回填。单实例 Scheduler 直接周期读取数据库中的调度配置，不依赖 Redis。每个配置编码只保留一条最新记录，不使用配置版本号和 Redis Pub/Sub。

## 数据库迁移与初始管理员

API 服务启动时自动执行嵌入二进制的 Goose 迁移。数据库结构只能通过 `application/api/migrations/sql/` 中的版本化 SQL 修改，不使用 GORM AutoMigrate。Scheduler 只连接已经初始化的数据库，不执行迁移。

迁移会创建逻辑客户端 `web-admin` 和内置角色，但不会写入默认管理员或默认密码。创建首个管理员：

```bash
export LIGHTNING_USERMGR_PASSWORD='replace-with-a-strong-password'
go run ./application/usermgr --operation=create --username=admin --nickname=管理员 --role=super_admin
```

重置密码：

```bash
export LIGHTNING_USERMGR_PASSWORD='replace-with-a-new-strong-password'
go run ./application/usermgr --operation=reset --username=admin
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

`make build` 同时构建 `bin/quick-admin-api` 和 `bin/quick-admin-scheduler`。

涉及通用工具、配置、认证基础设施和中间件的改动必须补充单元测试。controller/service 中的具体业务逻辑不强制单测，可按风险补充集成或契约测试。

## Docker

```bash
export LIGHTNING_JWT_SECRET='replace-with-at-least-32-random-characters'
docker compose up --build
```

PostgreSQL、Redis 和 RustFS 带有本地默认值；需要覆盖时使用 `LIGHTNING_POSTGRES_PASSWORD`、`LIGHTNING_REDIS_PASSWORD` 和 `LIGHTNING_RUSTFS_*`。

根目录的 `Dockerfile.api` 和 `Dockerfile.scheduler` 分别构建两个服务，唯一的 `docker-compose.yml` 编排 API、Scheduler、PostgreSQL、Redis 和 RustFS。RustFS 提供 S3 兼容对象存储；API 启动时会检查并按需创建 `quick-admin` Bucket，多 API 并发创建由二次检查保证幂等。

HTTP 端口为 `9009`，RustFS S3 API 与控制台端口分别为 `9000`、`9001`，API 探针为 `/health/startup`、`/health/live`、`/health/ready`。

## 代码边界

```text
application/api        API 入口及其 controller、router、service、DTO、中间件
application/scheduler  Scheduler 入口与 Job 定义
application/usermgr    一次性管理工具
internal               多个应用共享、但禁止仓库外导入的基础能力与适配器
```

项目不使用顶层 `pkg`。当前共享代码并不是提供给其他 Go Module 使用的公共 SDK，因此使用 Go 编译器强制约束的 `internal` 更准确；只有未来出现明确、稳定且需要被仓库外工程导入的 API 时，才新增 `pkg`。
