# Lightning

Lightning 是一个后台管理系统脚手架 Monorepo。同一套业务能力提供三种 Go 后端实现，并由一个 React 前端通过统一 HTTP 契约接入。

项目以 `native` 作为业务语义和接口行为基线，用于探索不依赖完整开源微服务框架的渐进式微服务工程；`kratos` 和 `gozero` 用于对照实现相同业务能力，而不是形成三套不同的接口。

## 工程组成

```text
lightning/
├── native/      # 原生 Go 实现，业务与 HTTP 契约基线
├── kratos/      # Kratos 框架实现
├── gozero/      # go-zero 框架实现
├── web-react/   # React + TypeScript 管理端
├── AGENTS.md    # 代码代理开发约定
└── README.md
```

### native

`native` 是当前主要演进的后端实现，Go Module 为 `github.com/gcc798/lightning`，使用 Go 1.26.5。

它从模块化单体起步，通过新增 `application/<service>` 组合共享 Module，按真实的扩容、发布和隔离需求渐进拆分服务。当前包含：

- `application/api`：无状态 HTTP API，拥有路由、controller、service、DTO、中间件、OpenAPI 和数据库迁移。
- `application/scheduler`：有状态定时任务进程，不提供 HTTP 服务，也不执行数据库迁移。
- `application/usermgr`：创建管理员、重置密码的一次性管理工具。
- `internal`：多个应用共享、但禁止当前 Go Module 之外导入的基础能力。

API 启动时执行嵌入二进制的 Goose 版本化迁移；数据库结构不使用 GORM `AutoMigrate`。验证码、微信、短信、邮件和 Scheduler 等运行期模块配置保存在 `s_config`，不写入服务 YAML。

详细启动、配置、认证、迁移和部署说明见 [native/README.md](native/README.md)。

### kratos

基于 Kratos 的框架化实现，采用 `sys-api` 和 `sys-rpc` 分层：

- `api/system/v1`：Proto 契约。
- `application/sys-api`：对外 HTTP 服务。
- `application/sys-rpc`：内部 RPC 与数据访问服务。
- `application/sys-rpc/ent/schema`：Ent Schema。

`kratos` 应保持与 `native` 一致的外部 HTTP 契约。

### gozero

基于 go-zero 的框架化实现，保留 `sys-api` 和 `sys-rpc` 分层：

- `application/sys-api`：对外 API 服务。
- `application/sys-rpc`：内部 RPC 服务。
- `application/sys-api/sys.api`：API 描述文件。

`gozero` 应保持与 `native` 一致的外部 HTTP 契约。

### web-react

唯一的前端工程，使用 React、TypeScript、Vite 和 Ant Design。前端只依赖一套 HTTP 契约，不为不同后端实现编写兼容分支。

开发服务默认运行在 `http://localhost:3001`，后端 API 默认地址为 `http://localhost:9009`。详细说明见 [web-react/README.md](web-react/README.md)。

## 快速启动 native

本地依赖 Go 1.26.5、PostgreSQL 16；API 另外使用 Redis 7。首次启动先生成不会提交到 Git 的本地配置：

```bash
cd native
make init-config
```

根据需要修改 `application/api/conf.dev.yaml`，然后启动 API：

```bash
go run ./application/api
```

需要运行定时任务时，另开进程启动 Scheduler：

```bash
go run ./application/scheduler
```

也可以通过 Docker Compose 启动 API、Scheduler、PostgreSQL、Redis 和 RustFS：

```bash
cd native
export LIGHTNING_JWT_SECRET='replace-with-at-least-32-random-characters'
docker compose up --build
```

## 开发与质量检查

各 Go 子工程拥有独立的 `go.mod`，应在各自目录执行命令。

```bash
# native：test、vet、race、Swagger freshness、Docker build
(cd native && make ci)

# Kratos：生成代码、测试与构建
(cd kratos && make proto-all && make ent && make wire && make test && make build-all)

# go-zero
(cd gozero && go test ./... && make build-all)

# React 前端
(cd web-react && pnpm install && pnpm build)
(cd web-react && pnpm dev)
```

## 契约与开发原则

- `native` 是接口路径、参数、响应结构、错误语义和业务规则的基线。
- `kratos` 与 `gozero` 可以采用不同工程分层，但不得发明不同的对外业务语义。
- `web-react` 只面向统一 HTTP 契约；框架后端无法接入时优先修正后端。
- 本仓库是全新的试验工程，默认不兼容旧代码、旧配置、旧接口或旧数据。
- 不为了微服务提前拆分；新增独立服务时复用稳定 Module，避免推翻核心业务。
- 生成代码应通过对应框架命令重新生成，不手工修改生成结果。

完整开发约定见 [AGENTS.md](AGENTS.md)。
