# native 重构记录与实施结果

> 状态日期：2026-07-22
> 当前状态：核心重构已经实施并完成验证。
> 适用范围：`native/`。除已经完成的 `web-react` 客户端 ID 调整外，本方案不修改 `kratos/`、`gozero/` 和 `web-react/`。

## 1. 文档目的

本文汇总本轮对 `native` 工程的评估、已完成整改，以及 Module、运行期配置和 Scheduler 服务拆分的最终实现。

本文同时作为后续演进依据。讨论过程中出现但最终被否决的方案，不再作为实现目标：

- 不为 `s_config` 增加 `version` 字段。
- 不使用 Redis Pub/Sub 同步 API 实例的配置。
- 不让每个 API 实例启动 Scheduler 或其他异步任务。
- 不在 API 启动时无条件把全部 `s_config` 配置覆盖写入 Redis。
- 不实现旧配置、旧接口、旧数据或旧代码的兼容逻辑。

## 2. 已确认的工程原则

### 2.1 渐进式微服务是首要架构纲领

`native` 定位为**不依赖 Kratos、go-zero 等开源微服务框架的渐进式微服务工程**。

工程从模块化单体开始，在当前业务规模下优先保持简单部署和进程内调用；当某项能力出现独立扩容、独立发布、故障隔离、资源隔离或独立生命周期的真实需求时，再通过增加新的 `application/<service>` 主程序把对应 Module 组合成独立服务。

```text
第一步：模块化单体

application/api
  -> user / auth / config / wechat / captcha Modules

第二步：按真实需求增加进程

application/api       -> Web API Modules
application/scheduler -> Scheduler / Queue Modules
application/<service> -> 某个需要独立部署的业务 Modules
```

渐进式微服务的目标是：

- 不为了使用“微服务”概念而提前拆分进程、数据库和 RPC。
- 不在没有独立部署需求时引入服务发现、远程调用、分布式事务等复杂度。
- 从第一天建立清晰模块边界，使未来拆服务主要是新增主程序、装配模块和增加传输适配器，而不是重写业务逻辑。
- 不把“无缝增加微服务”理解为完全没有接入工作；独立进程仍然需要 HTTP/RPC、部署和监控适配，但核心业务与模块生命周期不应推翻重写。

为保证后续可拆分，当前工程设计必须遵守：

1. `application/<service>` 同时拥有进程入口、依赖装配和该服务专属的领域代码；API 的 controller、router、service、DTO 不再散落到共享目录。
2. `internal/` 只保存多个应用共享、但禁止当前 Go Module 之外导入的能力，不能反向依赖具体 `application/*`。
3. controller/transport 只负责协议转换，核心能力通过稳定的 Go 接口调用。
4. 同进程时使用 Go 接口，不为了未来可能拆分而提前把内部调用改成 RPC。
5. 真正拆分时在既有能力接口外增加 HTTP/RPC Client/Server Adapter，不改写核心业务语义。
6. Module 不直接缓存其他 Module 的底层 Manager；跨模块依赖使用稳定能力接口。
7. 当前服务共享 PostgreSQL；Redis 仅由实际需要共享缓存或状态的服务装配。新增模块应明确数据所有权，跨模块业务调用优先经过能力接口，不扩散对其他模块数据表的直接读写。
8. Logger、配置加载、数据库、Redis、健康检查、信号处理、优雅关闭等进程基础能力复用统一实现，每个新应用按需装配。
9. Module 的 `Init/Start/Stop/Refresh` 生命周期不能假定自己永远运行在 API 进程内。
10. 是否拆为独立服务由真实运行需求决定，不按照 controller、service、model 等代码分层机械拆分。

本次 Module、Redis 运行期配置和 API/Scheduler 拆分设计均以此为最高优先级。若后续局部设计与该纲领冲突，应优先保证渐进演进能力。

### 2.2 工程边界

- `native/` 是业务语义和 HTTP 契约基线。
- 当前重构只针对 `native/`。
- `kratos/`、`gozero/` 后续需要迁移同等能力时，再单独实施。
- 数据库迁移由拥有 Schema 的服务管理；当前只有 `application/api/migrations` 执行迁移，Scheduler 和另外两套后端不隐式复用或执行该包。
- `web-react` 只面向一套统一 HTTP 契约，不为不同后端编写兼容分支。

### 2.3 不做向后兼容

本仓库是全新的试验工程。任何改动默认不考虑旧代码、旧配置、旧接口或旧数据的兼容，只有开发者明确提出兼容要求时才实现兼容逻辑。

因此，后续把验证码、微信、短信、邮件和调度器配置迁出 YAML 时，应直接删除旧配置结构、旧环境变量绑定和旧初始化路径，不保留双读、回退或弃用期。

### 2.4 测试边界

- 通用工具、配置解析、基础设施、中间件、锁、缓存、生命周期管理必须有单元测试。
- controller、service 中的具体业务代码不强制编写单元测试。
- 高风险业务流程仍应通过契约测试或集成测试验证。
- 涉及并发状态的实现必须通过 `go test -race`。

### 2.5 登录与 CORS

- 小程序登录采用纯微信登录：前端提交 `wxCode`，后端调用微信 `Code2Session` 获取 OpenID/UnionID。
- 不接受前端直接提交 OpenID 作为微信登录凭证。
- 逻辑客户端 ID 统一为 `web-admin`，不使用客户端密钥。
- 开发环境可以开启 CORS。
- 生产环境默认关闭应用层 CORS，由 Nginx 同源代理前端和 API。

## 3. 已完成的整改

以下内容已经在当前仓库中落地，不属于下一阶段待实现事项。

### 3.1 强类型启动配置

启动配置已经统一为强类型结构，通过 Viper 加载：

```text
代码默认值 < conf.{dev|prod}.yaml < LIGHTNING_* 环境变量
```

- `LIGHTNING_APP_ENV` 只用于选择 `dev` 或 `prod`，默认 `dev`。
- 配置结构体中没有 `env` 属性。
- API 与 Scheduler 在各自的 `application/<service>/` 中管理配置，不再使用顶层 `configs/`。
- Git 只管理 `conf.example.yaml`、`zaplogger.example.yaml`；真实的 dev/prod 配置被忽略并通过 `make init-config` 初始化。
- 生产 example 中数据库、Redis、RabbitMQ、S3 属性具有明确默认值，部署时可由 `LIGHTNING_*` 环境变量覆盖。
- 配置加载、环境变量映射和配置校验已有单元测试。

这些 YAML 只保存对应进程实际需要的启动基础设施配置。API 配置 HTTP、数据库、Redis、JWT 等能力；Scheduler 当前只配置数据库。

### 3.2 Goose 版本化迁移

- Goose 是唯一数据库结构变更入口。
- 已移除 GORM `AutoMigrate` 与少量 Goose 混用的模式。
- 迁移包和 SQL 位于 `application/api/migrations/`，嵌入 API 二进制并仅在 API 启动时执行。
- Scheduler 只使用已经初始化的表，不执行迁移；未来某个服务拥有独立 Schema 时，在自己的 `application/<service>/migrations/` 管理。
- 已删除不再需要的 `native/deploy/sql` 初始化脚本。
- 当前迁移包含基线结构、移除 Casbin 和系统参考数据。
- 迁移写入逻辑客户端与内置角色，不写入默认管理员和默认密码。
- 初始管理员由 `application/usermgr` 显式创建。

### 3.3 Token 与认证

- JWT 使用唯一 `jti`、固定 issuer 和明确的签名算法校验。
- Access Token 与 Refresh Token 均关联 Redis 服务端会话。
- Refresh Token 单次使用，刷新后立即轮换。
- 登出会撤销当前会话，旧 Access Token 不能继续访问接口。
- 已移除重复的并发登录管理实现。
- 登录业务已从臃肿 controller 下沉到 AuthService。
- 微信登录已经采用 `wxCode -> Code2Session -> OpenID/UnionID` 流程。
- 新微信用户创建和默认角色分配在事务中完成。

### 3.4 HTTP 契约、日志和错误处理

- 业务错误返回真实 HTTP 状态码，不再全部返回 HTTP 200。
- 已补充 HTTP 响应和 OpenAPI 契约测试。
- 操作日志组件具有明确启动、停止和队列生命周期。
- 操作日志会脱敏敏感字段并限制请求体大小。
- Swagger 已重新生成，并加入 freshness 检查。

### 3.5 管理工具与部署

- `application/usermgr` 使用显式 flags。
- 管理员密码只通过 `LIGHTNING_USERMGR_PASSWORD` 传入，不作为命令行参数，也不回显。
- API 与 Scheduler 分别使用根目录 `Dockerfile.api`、`Dockerfile.scheduler`，根目录只保留一个 `docker-compose.yml`。
- Compose 使用 RustFS 提供 S3 兼容存储，API 启动时幂等确保目标 Bucket 存在。
- 已删除旧 `dockerfile/`、`monitoring/` 和过时 `docs/`，Swagger 生成契约迁入 `application/api/openapi/`。
- 容器使用非 root 用户，HTTP 服务端口统一为 `9009`。
- 已配置 startup、live、ready 健康检查。

### 3.6 CI 验证入口

`native/Makefile` 已提供：

```bash
make verify # go test、go vet、race、Swagger freshness
make ci     # verify + Docker build
```

### 3.7 已完成验证

本轮整改完成后已经验证：

- `go test ./...`
- `go vet ./...`
- `go test -race ./...`
- Swagger freshness
- PostgreSQL 16 空库执行完整 Goose 迁移
- Docker build
- Docker Compose 配置解析
- RustFS 健康检查、API 自动创建 S3 Bucket，并在验证后清理隔离容器与数据卷
- API 与 Scheduler 两个进程可同时连接同一套 PostgreSQL 并正常启动和优雅停止，Scheduler 不连接 Redis
- 两个 API 实例使用不同端口同时启动并分别通过 readiness，API 进程未启动 Scheduler Job
- `s_config` 已升级到迁移版本 4，五类运行期配置均存在，`code` 唯一且没有 `version` 字段
- 通过配置管理接口修改 `auth.captcha` 后，验证码能力下一次调用立即读取到 Redis 中的新配置；验证后已恢复默认关闭状态
- Scheduler Module 在运行期配置关闭后停止已注册 Job，停止 Module 后不再执行任务
- 创建测试管理员并完成登录
- 使用 Access Token 访问配置接口
- 登出后原 Token 返回未授权

## 4. 已实施的最终架构

运行期业务模块已经从静态 `container` 初始化方式中拆出，Web API 和异步任务由两个独立进程承载。

```text
                    +-------------------+
                    |     s_config      |
                    | 持久化、管理、审计 |
                    +----+---------+----+
                         |         |
                 缓存缺失时加载   周期读取调度配置
                         |         |
                    +----v----+    +v---------------------+
                    |  Redis  |    | application/scheduler |
                    | API共享 |    | 数据库驱动异步任务     |
                    +----+----+    +-----------------------+
                         |
                    +----v-----------+
                    | application/api |
                    | 无状态 Web      |
                    +-----------------+
```

### 4.1 演进方式

API/Scheduler 拆分是渐进式微服务原则的第一次落地，不代表立即把全部业务改造成分布式系统：

- API 与 Scheduler 服务复用 Module 和 Container 生命周期抽象，但各自只装配真实需要的基础设施；API 使用 ConfigStore/Redis，Scheduler 只使用数据库。
- 每个 `application/<service>` 拥有自己的入口和领域代码，跨服务复用内容才进入 `internal/`，从目录层面减少并行开发时的 Git 冲突。
- 未来增加新的微服务时沿用同一模式：新增应用目录，选择需要的 Module，补充领域代码、transport、健康检查和部署配置。
- Module 在同一进程内时继续通过 Go 接口调用；只有真实跨进程后才增加远程适配器。
- 新服务不能通过复制 API 工程再独立演化的方式创建，否则会形成新的多套业务实现。

顶层不保留 `pkg/`。当前共享包不对外承诺稳定 API，因此使用 Go 编译器强制约束的 `internal/`；只有未来明确需要被其他 Go Module 导入的稳定库才进入 `pkg/`。

推荐保持以下依赖方向：

```text
application/<service>
    -> 进程装配 / transport / bootstrap
        -> Module 能力接口与实现
            -> domain / repository 接口
                -> PostgreSQL、Redis、第三方 SDK 等基础设施适配
```

上层可以依赖下层接口，下层不能反向依赖具体主程序或 HTTP controller。

### 4.2 服务职责

#### `application/api`

API 服务只负责同步 Web 能力：

- 用户、角色、菜单等管理接口。
- 登录、登出、Token 刷新。
- 微信登录。
- 验证码申请与校验。
- 短信、邮件业务请求提交。
- 配置管理接口。

API 服务不启动 Scheduler，不运行定时任务和常驻业务协程，可以在 Nginx 或负载均衡器后水平扩容。

“无状态”表示业务正确性不依赖某个 API 进程的内存状态；并不禁止进程复用数据库连接池、Redis Client 或根据配置生成的第三方 SDK Client。

#### `application/scheduler`

Scheduler 服务负责：

- Scheduler。
- 定时任务。
- 队列消费。
- 异步短信、邮件等后台任务。
- 后续需要持续运行的模块协程。

Scheduler 服务默认单实例运行。需要多实例时，任务必须通过 Redis 分布式锁、可靠队列或任务租约保证不会重复执行。

Scheduler 只能调度代码中显式注册的 Job，不提供从数据库执行任意命令或任意代码的能力。

## 5. Module 与 Container 设计

### 5.1 Module 接口

```go
type Module interface {
	Name() string
	Init(ctx context.Context, cont Container) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Refresh(ctx context.Context, req ModuleRefreshRequest) error
}

type ModuleRefreshRequest struct {
	Codes  []string
	Reason string
}
```

`ModuleRefreshRequest` 不包含版本号。Redis 和数据库中每个 `code` 只保存最新配置。

### 5.2 Container 接口

```go
type Container interface {
	RegisterModules(ctx context.Context, modules ...Module) error
	GetModule(name string) Module
	StartModules(ctx context.Context) error
	StopModules(ctx context.Context) error
	RefreshModule(ctx context.Context, name string, req ModuleRefreshRequest) error
}
```

Container 实现需要满足：

- 模块注册和启动顺序确定，不依赖 Go map 遍历顺序。
- 停止顺序与启动顺序相反。
- 启动中途失败时，回滚已经成功启动的模块。
- 停止全部模块并合并错误，不能遇到第一个错误就中断。
- 调用 Module 生命周期方法时不长期持有 Container 全局锁。
- 同一模块的 `Refresh` 串行执行。
- Module 名称不可重复。

### 5.3 稳定能力入口

业务层依赖 Module 提供的能力接口，不直接依赖 Viper、`s_config` 数据模型或某次初始化创建出的 Manager：

```go
type WeChatModule interface {
	Module
	Code2Session(ctx context.Context, code string) (*WeChatSession, error)
}
```

AuthService 始终调用同一个 `WeChatModule`。Module 内部可以使用最新 Redis 配置创建或替换微信 Client，业务层不缓存底层 Client。

短信、邮件和验证码采用相同方式。Captcha Module 依赖 SMS Module、Email Module 的稳定发送接口，不保存某次启动时获得的具体 Manager 指针。

### 5.4 模块划分

第一批运行期模块：

- `wechat`
- `sms`
- `email`
- `captcha`
- `scheduler`

API 注册：

```text
wechat、sms、email、captcha
```

Scheduler 服务注册：

```text
scheduler
```

模块根据实际依赖按确定顺序初始化。例如 Captcha 在 SMS、Email 之后初始化。

## 6. 运行期配置设计

### 6.1 哪些配置迁移到 `s_config`

从 `conf.dev.yaml`、`conf.prod.yaml` 和强类型启动配置中删除：

- `captcha`
- `wechat`
- `sms`
- `email`
- `scheduler`

建议配置编码：

| Module | `s_config.code` |
| --- | --- |
| 微信 | `integration.wechat` |
| 短信 | `integration.sms` |
| 邮件 | `integration.email` |
| 验证码 | `auth.captcha` |
| 调度器 | `scheduler` |

数据库、Redis、JWT、HTTP、CORS、RabbitMQ、S3、WebSocket 等基础设施仍属于启动配置，但只出现在实际使用它们的服务 YAML 中。当前 Scheduler 只保留数据库配置。

原有 MQTT 实现绑定门禁设备 Topic 协议和消息重试表，不属于通用后台脚手架的基础能力，已经删除。未来出现明确的 IoT 业务后，应由对应服务或可选模块自行拥有协议、配置、协程和数据表，而不是重新放入共享基础设施。

### 6.2 `s_config` 约束

- `code` 必须唯一。
- 同一个 `code` 只保留一条最新配置。
- `data` 使用非空 JSONB。
- 不增加 `version` 字段。
- 更新采用最后一次成功写入生效的语义。
- 默认配置通过新的 Goose 迁移写入，不通过 Go 启动代码偷偷补表或补数据。

`updated_time` 只用于展示和审计，不承担版本控制或消息顺序控制。

### 6.3 Redis 的职责

Redis 保存所有 API 实例共享的请求驱动模块配置：

```text
quick-admin:config:{code}
```

例如：

```text
quick-admin:config:integration.wechat
quick-admin:config:auth.captcha
```

规则：

- 配置 key 默认不设置 TTL。
- API 启动时不无条件覆盖写入全部配置。
- Redis 存在配置时直接返回。
- Redis key 不存在时，从 `s_config` 加载并回填 Redis。
- “key 不存在”和“Redis 服务不可用”必须区分；Redis 不可用时不能伪装成配置不存在。
- 正常配置修改必须通过配置管理接口完成。
- 如果开发者直接修改数据库，需要同时清理对应 Redis key；不提供自动兼容或自动发现数据库直改的能力。

Redis 是共享运行期配置层，`s_config` 是持久化事实来源。Redis 中不保存微信 SDK Client、SMTP Client 等进程内对象。

### 6.4 `ConfigStore.GetOrLoad`

```text
GET Redis 配置
    |
    +-- 存在 --> 返回
    |
    +-- 不存在 --> 获取该 code 的分布式锁
                       |
                       +-- 再次 GET Redis（双重检查）
                       |
                       +-- 查询 s_config
                       |
                       +-- SET Redis
                       |
                       +-- 释放锁
```

没有获得锁的实例进行带随机抖动的短暂等待，然后重新读取 Redis。等待必须响应 `context.Context` 取消。

该流程避免 Redis 被清空后，多个 API 实例同时查询数据库并重复构建同一个配置。Scheduler 不走该缓存链路。

### 6.5 配置修改

同一个 `code` 的加载、创建、修改和删除共用同一把分布式锁：

```text
quick-admin:lock:config:{code}
```

更新流程：

```text
校验配置数据
    |
获取 code 分布式锁
    |
删除 Redis 旧配置
    |
事务更新 s_config
    |
写入 Redis 新配置
    |
刷新当前进程对应 Module
    |
释放锁
```

删除旧 Redis 值失败时，不更新数据库。数据库更新失败时保持缓存缺失，锁释放后由下一次 `GetOrLoad` 重新加载数据库当前值。数据库更新成功但 Redis 写入失败时同样保持缓存缺失并返回错误，Redis 恢复后由 `GetOrLoad` 回填数据库最新值。

同一 `code` 的并发修改由锁串行化，避免较早请求在较晚请求之后覆盖 Redis，导致 Redis 与数据库结果不一致。

这里不使用 Redis Pub/Sub。其他 API 实例无需接收通知，下一次能力调用会读取同一个 Redis key。

### 6.6 分布式锁要求

加锁使用随机 token：

```text
SET quick-admin:lock:config:{code} {token} NX PX {ttl}
```

释放锁必须通过 Lua 比较 token 后删除：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
end
return 0
```

要求：

- 锁必须有 TTL，进程崩溃后能够自动释放。
- 只有持有相同 token 的调用者可以释放锁。
- 锁等待支持 Context 取消和超时。
- 重试使用小幅随机退避。
- 业务函数出现错误或 panic 时仍需安全释放锁。
- 若单次操作可能超过锁 TTL，需要实现续租；当前配置读写应控制在短操作范围内，避免不必要的续租复杂度。

### 6.7 Module 如何获得最新配置

请求驱动的 Module 在每次能力调用时读取 Redis 配置：

```text
AuthService
    |
WeChatModule.Code2Session
    |
ConfigStore.Get("integration.wechat")
    |
Redis 最新配置
```

如果第三方 SDK Client 创建成本较高，Module 可以保存由配置内容摘要标识的进程内 Client：

- 配置摘要未变化：复用已有 Client。
- 配置摘要发生变化：先创建和校验新 Client，成功后原子替换。
- 创建失败：不把无效 Client 写入运行时快照。

该缓存只是性能优化，不是业务状态来源，因此不破坏 API 无状态部署。

Scheduler 等长期运行模块不能依赖“下一次 HTTP 请求”触发刷新。当前单实例 Scheduler Module 按固定短周期直接读取 `s_config`、比较配置摘要并协调已注册 Job；具体刷新周期由运行期配置确定，因此不需要 Redis 启动配置。

## 7. 启动与关闭顺序

### 7.1 API

```text
加载 YAML/环境变量启动配置
    -> 初始化 Logger
    -> 初始化 PostgreSQL
    -> 执行 Goose 迁移
    -> 初始化 Redis
    -> 初始化 JWT 等基础设施
    -> 注册并 Init API Modules
    -> 注册 Router
    -> StartModules
    -> 启动 HTTP Server
```

关闭：

```text
停止接受 HTTP 请求
    -> 等待在途请求结束
    -> StopModules（反向顺序）
    -> 停止日志写入器
    -> 关闭 Redis、数据库等资源
```

### 7.2 Scheduler 服务

```text
加载启动配置
    -> 初始化 Logger、PostgreSQL
    -> 注册并 Init Scheduler Modules
    -> 注册代码 Job
    -> StartModules
    -> 等待退出信号
```

关闭时先停止产生新任务，再等待在途任务完成，最后反向停止 Module 和基础设施。

## 8. 实施记录

### 阶段一：Module 基础设施（已完成）

- 定义 `Module`、`ModuleRefreshRequest` 和新的 `Container` 生命周期接口。
- 用确定顺序替代当前 `Component` 切片的简单启动逻辑。
- 实现启动失败回滚、反向停止、错误合并和刷新串行化。
- 把进程装配与 Module 业务能力分离，保证 Module 不依赖 `application/api`。
- 提取可供未来所有 `application/<service>` 复用的 Logger、数据库、Redis、健康检查和优雅关闭初始化能力。
- 补齐 Container 生命周期单元测试和 race 测试。

### 阶段二：ConfigStore 与分布式锁（已完成）

- 实现通用 Redis 分布式锁。
- 实现 `ConfigStore.GetOrLoad`。
- 实现已知配置 code 的类型解析与校验。
- 补齐锁竞争、安全释放、TTL、Context 取消、双重检查和缓存回填测试。

### 阶段三：数据库与配置管理（已完成）

- 新增 Goose 迁移，为 `s_config.code` 增加唯一约束、收紧 `data` 约束并写入五类默认配置。
- 配置服务按 `code` 管理最新一条配置。
- 创建、修改和删除操作接入同一把 code 分布式锁。
- 更新 Swagger 生成契约和配置管理接口。

### 阶段四：业务 Module（已完成）

- 实现 WeChat、SMS、Email、Captcha Module。
- 将 AuthService、Captcha 等消费者切换为稳定 Module 接口。
- 删除 YAML 中的五类运行期配置和对应环境变量绑定。
- 删除 Container 中原有静态第三方 Manager、Captcha 和 Scheduler 初始化逻辑。
- 不保留旧配置兼容路径。

### 阶段五：Scheduler 服务拆分（已完成）

- 新增 `application/scheduler`。
- 将 Scheduler、定时任务和常驻业务协程迁移到 Scheduler 服务。
- 确保 `application/api` 不再启动异步任务。
- 为 Scheduler 配置检查、Job 注册和优雅停止补齐测试。
- 更新根目录 Dockerfile、Docker Compose、Kubernetes 和 README，分别描述 API 与 Scheduler 服务。

### 阶段六：完整验证（已完成）

- `go test ./...`
- `go vet ./...`
- `go test -race ./...`
- Swagger freshness
- PostgreSQL 空库迁移
- Redis 冷缓存并发回填
- 多 API 实例读取同一配置
- 同一 code 并发更新一致性
- 配置修改后微信、短信、邮件、验证码下一次调用生效
- Scheduler 服务配置变化后重新协调 Job
- API 多实例不重复执行异步任务
- Docker build 与部署配置检查

## 9. 验收标准

重构完成需要同时满足：

1. `application/api` 可以多实例部署，实例内没有 Scheduler 或业务后台协程。
2. `application/scheduler` 可以独立启动、停止和运行异步任务。
3. 五类运行期配置不再出现在 YAML 强类型配置和环境变量绑定中。
4. `s_config.code` 唯一，不存在版本字段，同一 code 只保存最新配置。
5. 所有 API 实例从同一个 Redis key 获得请求驱动模块配置；Scheduler 直接读取 `s_config`，均不依赖 Pub/Sub。
6. Redis 缺少 key 时，通过分布式锁完成一次数据库加载和回填。
7. 同一 code 的并发修改不会造成 Redis 与数据库最终值顺序倒置。
8. 业务代码只依赖稳定 Module 能力接口，不缓存底层 Manager。
9. 配置变化后，新业务调用使用最新 Redis 配置；长期运行的 Scheduler Module 能在规定检查周期内从数据库读取并生效。
10. Module、ConfigStore、分布式锁和其他通用工具具有完整单元测试，并通过 race 检查。
11. 不包含旧配置、旧接口或旧数据兼容代码。
12. `make ci` 全部通过。
13. Module 不依赖具体 `application/*`，同一 Module 能被 API、Scheduler 服务或未来新主程序按需装配。
14. 新增一个独立服务时不需要复制或重写核心业务，只需要新增主程序、transport 和部署配置。
15. 当前进程内模块调用保持 Go 接口，不存在为了“微服务化”而提前引入的内部 RPC。

## 10. 明确不在当前范围内的事项

- 不同步改造 `kratos/` 和 `gozero/`。
- 不为不同后端增加前端兼容分支。
- 不将数据库、Redis、JWT 等进程启动基础设施配置迁入 `s_config`。
- 不增加配置历史表、配置版本号、乐观锁或回滚功能。
- 不使用 Redis Pub/Sub。
- 不把任意脚本或命令存入数据库交给 Scheduler 执行。
- 不为 controller/service 业务代码强制追求单元测试覆盖率。
- 不在缺少独立扩容、发布或隔离需求时机械拆分微服务。
- 不引入 Kratos、go-zero 等完整微服务框架；需要的进程生命周期与基础设施能力由 `native` 自身保持轻量、明确和可测试的实现。
