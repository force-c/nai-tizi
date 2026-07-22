# AGENTS.md

本文件供 Codex、Claude Code 和其他代码代理快速理解仓库结构与开发约定。

## 项目定位

`quick.admin` 是一个后台系统脚手架工程。它不是单一后端工程，而是同一套后端业务逻辑的三套 Go 实现，加上一个 React 前端工程：

- `native/`：开发者从零手搓的原生 Go 后端实现，是业务语义和接口行为的基线。
- `kratos/`：基于 Kratos 开源微服务框架的实现。
- `gozero/`：基于 go-zero 开源微服务框架的实现。
- `web-react/`：唯一的前端工程，应尽量无缝对接三套后端实现。

本质上，这是一套后台管理系统能力在三种后端框架里的对照实现。做需求或修 bug 时，先理解 `native` 的业务语义，再把同等行为落到目标框架实现中。

## 核心原则

1. `native/` 是业务基线。接口路径、请求参数、响应结构、错误语义和业务规则优先参考它。
2. `kratos/` 和 `gozero/` 是框架化实现，不应发明不同的业务语义。
3. `web-react/` 是同一套前端，通过 HTTP 契约与后端交互；不应为了不同后端实现让前端做兼容分支。
4. 三套后端可以有不同工程分层和框架代码生成方式，但对外业务能力应保持一致。
5. 只要 `web-react` 能和 `native` 正常交互，`kratos` 和 `gozero` 也必须提供兼容的 HTTP 契约。
6. 修改生成文件时优先使用对应框架命令重新生成，不要只做字符串硬改。
7. 本仓库是全新的试验工程，任何改动默认不考虑旧代码、旧配置、旧接口或旧数据的向后兼容；只有开发者明确提出兼容要求时，才实现兼容逻辑。
8. `native/` 是不依赖完整开源微服务框架的渐进式微服务工程：从模块化单体起步，按真实扩容、发布和隔离需求通过新增 `application/<service>` 组合既有 Module；每个目录拥有该服务的入口与领域代码，不得为了微服务而提前拆分，也不得让未来拆服务必须推翻重写核心业务。

## 根目录结构

```text
quick.admin/
├── native/      # 原生 Go 后端，业务基线，开发者手搓实现
├── kratos/      # Kratos 微服务框架版本
├── gozero/      # go-zero 微服务框架版本
├── web-react/   # React 前端
├── README.md
├── AGENTS.md
└── CLAUDE.md
```

## 后端实现关系

### native

`native/` 是从零手搓的 Go 后端实现，主要用于确认业务事实。

常见入口与边界：

- `native/application/api/`：HTTP API 入口及其 router、controller、service、DTO 和中间件。
- `native/application/scheduler/`：Scheduler 入口与异步 Job。
- `native/application/usermgr/`：管理员维护工具。
- `native/internal/domain/model/`：GORM 模型。
- `native/application/api/migrations/`：API 所拥有数据库 Schema 的 Goose 迁移及版本化 SQL；Scheduler 不执行迁移。
- `native/internal/database/`：多个应用共享的数据库连接与 GORM 插件，不负责执行迁移。
- `native/internal/`：多个应用共享、但通过 Go `internal` 规则禁止仓库外导入的基础能力。
- `native/application/api/openapi/`：API 自有的 Swagger 生成文件，`application/api/main.go` 会导入。

API 与 Scheduler 各自在服务目录保存 `conf.example.yaml` 和 `zaplogger.example.yaml`。真实的 `conf.dev.yaml`、`conf.prod.yaml`、`zaplogger.dev.yaml`、`zaplogger.prod.yaml` 不纳入 Git，使用 `make init-config` 初始化。
每个服务的 YAML 只声明自身实际依赖；当前 API 配置完整基础设施，Scheduler 配置只包含数据库，不得为了复用配置结构加入未使用字段。

`native` 不使用顶层 `pkg/` 存放普通共享代码。仅当某个包明确作为稳定 API 供当前 Go Module 之外的工程导入时，才考虑新增 `pkg/`；仓库内多应用共享代码应保留在 `internal/`。

常用命令：

```bash
cd native
go test ./...
make swagger
```

### kratos

`kratos/` 是基于 Kratos 的微服务实现，采用 `sys-api` + `sys-rpc` 分层。

常见入口：

- `kratos/api/system/v1/`：proto 契约。
- `kratos/application/sys-api/`：对外 HTTP 服务。
- `kratos/application/sys-rpc/`：内部 gRPC 服务和核心数据访问。
- `kratos/application/sys-rpc/ent/schema/`：Ent schema。
- `kratos/pkg/`：Kratos 版本沉淀的通用能力。
- `kratos/Makefile`：代码生成、构建和测试入口。

常用命令：

```bash
cd kratos
make conf
make proto
make ent
make wire
go test ./...
```

注意：proto 相关 Go 代码包含编码后的 descriptor。修改 `go_package` 或 module 路径后，要用 `make conf` / `make proto` 重新生成，不能只做文本替换。

### gozero

`gozero/` 是基于 go-zero 的微服务实现，保留 `sys-api` / `sys-rpc` 分层。

常见入口：

- `gozero/application/sys-api/`：对外 API 服务。
- `gozero/application/sys-rpc/`：内部 RPC 服务。
- `gozero/application/sys-api/sys.api`：API 描述文件。
- `gozero/Makefile`：构建入口。

常用命令：

```bash
cd gozero
go test ./...
```

## 前端工程

### web-react

`web-react/` 是唯一的前端工程。它只应依赖统一 HTTP 契约，不应因为后端实现不同而写兼容逻辑或分支判断。

前端契约原则：

- 前端只面向一套接口语义。
- `native` 是前端契约基线。
- 如果 `web-react` 可以和 `native` 正常交互，`kratos` 和 `gozero` 也应该无缝可用。
- 发现某个框架后端无法被同一套前端调用时，优先修后端契约，而不是改前端兼容。

常用命令：

```bash
cd web-react
pnpm install
pnpm dev
pnpm build
```

默认开发服务端口参考 `web-react/README.md`。

## 开发建议

- 做业务改动前，先在 `native/` 找到对应 controller/service/model/request/response。
- 若目标是 `kratos` 或 `gozero`，实现时保持与 `native` 的 HTTP 契约和接口语义一致。
- 若改动会影响前端接口，优先确认是否破坏了 `native` 契约；不要让 `web-react` 为不同后端实现做特殊兼容。
- 每个 Go 子工程单独运行测试：`native`、`kratos`、`gozero` 各自都有自己的 `go.mod`。
- 不要把 `native`、`kratos`、`gozero` 当成互相引用的包；它们是同一业务的不同实现。
- 不要手工修改 `native/application/api/openapi/`，应通过 `make swagger` 重新生成。

## 快速判断应该看哪里

- 想知道业务原始行为：看 `native/`。
- 想改 Kratos 微服务版本：看 `kratos/api` 和 `kratos/application`。
- 想改 go-zero 微服务版本：看 `gozero/application`。
- 想看前端当前对接方式：看 `web-react/src`。
