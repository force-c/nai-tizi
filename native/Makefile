.PHONY: help swagger swagger-check swagger-fmt run build test test-race vet docker-build verify ci clean

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make swagger-fmt  - 格式化 Swagger 注释"
	@echo "  make run          - 运行服务（自动生成文档）"
	@echo "  make build        - 编译项目"
	@echo "  make test         - 运行测试"
	@echo "  make test-race    - 运行竞态检测"
	@echo "  make vet          - 运行 Go 静态检查"
	@echo "  make swagger-check - 检查 Swagger 是否为最新"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make verify       - 执行本地完整质量检查"
	@echo "  make ci           - 执行测试、vet、race、Swagger freshness 和 Docker build"
	@echo "  make clean        - 清理生成的文件"

# 生成 Swagger 文档
swagger:
	@echo "正在生成 Swagger 文档..."
	@swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal
	@echo "✅ Swagger 文档已生成: docs/swagger/"
	@echo "   - swagger.json"
	@echo "   - swagger.yaml"
	@echo "   - docs.go"

swagger-check:
	@tmp_dir=$$(mktemp -d); \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	swag init -g cmd/api/main.go -o "$$tmp_dir/swagger" --parseDependency --parseInternal >/dev/null; \
	diff -ru docs/swagger "$$tmp_dir/swagger"

# 格式化 Swagger 注释
swagger-fmt:
	@echo "正在格式化 Swagger 注释..."
	@swag fmt -g cmd/api/main.go

# 运行服务（先生成文档）
run: swagger
	@echo "正在启动服务..."
	@go run cmd/api/main.go

# 编译项目
build: swagger
	@echo "正在编译项目..."
	@go build -o bin/quick-admin cmd/api/main.go
	@echo "✅ 编译完成: bin/quick-admin"

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

test-race:
	@go test -race ./...

vet:
	@go vet ./...

docker-build:
	@docker build -f dockerfile/Dockerfile -t quick-admin-native:verify .

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
