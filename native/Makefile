.PHONY: help init-config swagger swagger-check swagger-fmt run run-scheduler build test test-race vet docker-build verify ci clean

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make init-config  - 从服务 example 初始化本地配置（不覆盖已有文件）"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make swagger-fmt  - 格式化 Swagger 注释"
	@echo "  make run          - 运行服务（自动生成文档）"
	@echo "  make run-scheduler - 运行独立 Scheduler 服务"
	@echo "  make build        - 编译项目"
	@echo "  make test         - 运行测试"
	@echo "  make test-race    - 运行竞态检测"
	@echo "  make vet          - 运行 Go 静态检查"
	@echo "  make swagger-check - 检查 Swagger 是否为最新"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make verify       - 执行本地完整质量检查"
	@echo "  make ci           - 执行测试、vet、race、Swagger freshness 和 Docker build"
	@echo "  make clean        - 清理生成的文件"

init-config:
	@test -f application/api/conf.dev.yaml || cp application/api/conf.example.yaml application/api/conf.dev.yaml
	@test -f application/api/conf.prod.yaml || cp application/api/conf.example.yaml application/api/conf.prod.yaml
	@test -f application/api/zaplogger.dev.yaml || cp application/api/zaplogger.example.yaml application/api/zaplogger.dev.yaml
	@test -f application/api/zaplogger.prod.yaml || cp application/api/zaplogger.example.yaml application/api/zaplogger.prod.yaml
	@test -f application/scheduler/conf.dev.yaml || cp application/scheduler/conf.example.yaml application/scheduler/conf.dev.yaml
	@test -f application/scheduler/conf.prod.yaml || cp application/scheduler/conf.example.yaml application/scheduler/conf.prod.yaml
	@test -f application/scheduler/zaplogger.dev.yaml || cp application/scheduler/zaplogger.example.yaml application/scheduler/zaplogger.dev.yaml
	@test -f application/scheduler/zaplogger.prod.yaml || cp application/scheduler/zaplogger.example.yaml application/scheduler/zaplogger.prod.yaml

# 生成 Swagger 文档
swagger:
	@echo "正在生成 Swagger 文档..."
	@swag init -g application/api/main.go -o application/api/openapi --parseDependency --parseInternal
	@echo "✅ Swagger 文档已生成: application/api/openapi/"
	@echo "   - swagger.json"
	@echo "   - swagger.yaml"
	@echo "   - docs.go"

swagger-check:
	@tmp_dir=$$(mktemp -d); \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	swag init -g application/api/main.go -o "$$tmp_dir/openapi" --parseDependency --parseInternal >/dev/null; \
	diff -ru application/api/openapi "$$tmp_dir/openapi"

# 格式化 Swagger 注释
swagger-fmt:
	@echo "正在格式化 Swagger 注释..."
	@swag fmt -g application/api/main.go

# 运行服务（先生成文档）
run: swagger
	@echo "正在启动服务..."
	@go run ./application/api

run-scheduler:
	@echo "正在启动 Scheduler 服务..."
	@go run ./application/scheduler

# 编译项目
build: swagger
	@echo "正在编译项目..."
	@go build -o bin/quick-admin-api ./application/api
	@go build -o bin/quick-admin-scheduler ./application/scheduler
	@echo "✅ 编译完成: bin/quick-admin-api, bin/quick-admin-scheduler"

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

test-race:
	@go test -race ./...

vet:
	@go vet ./...

docker-build:
	@docker build -f Dockerfile.api -t quick-admin-native-api:verify .
	@docker build -f Dockerfile.scheduler -t quick-admin-native-scheduler:verify .

verify: test vet test-race swagger-check

ci: verify docker-build

# 清理生成的文件
clean:
	@echo "正在清理..."
	@rm -rf bin/
	@echo "✅ 清理完成"

# 安装 swag 工具
install-swag:
	@echo "正在安装 swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ swag 安装完成"
