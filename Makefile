# lpcenter Project Management Makefile

# Basic configuration
SystemName="lpcenter"
GoVersion=$(shell go version | awk '{print $$3, $$4}')
BuildTime=$(shell date "+%Y-%m-%d %H:%M:%S")
GitCommitID=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
Version="0.0.1_base"

# Build mode
VER=debug

# Set Gin log mode based on build mode
ifeq ($(VER),debug)
	GinLogMode="debug"
else
	GinLogMode="release"
endif

# Build parameters
LDFLAGS="\
 -X 'github.com/ricky97gr/lpcenter/server/config.GoVersion=$(GoVersion)'\
 -X 'github.com/ricky97gr/lpcenter/server/config.SystemName=$(SystemName)'\
 -X 'github.com/ricky97gr/lpcenter/server/config.CommitID=$(GitCommitID)'\
 -X 'github.com/ricky97gr/lpcenter/server/config.BuildTime=$(BuildTime)'\
 -X 'github.com/ricky97gr/lpcenter/server/config.Version=$(Version)'\
 -X 'github.com/ricky97gr/lpcenter/server/config.GinLogMode=${GinLogMode}'\
"

# Target definitions
.PHONY: help stop all server web prepare revert docker-build docker-run docker-stop release build clean

# Display help information
help:
	@echo "lpcenter Project Management Tool"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Targets:"
	@echo "  help          - Display this help information"
	@echo "  build         - 编译所有产物 (前端+后端)"
	@echo "  build-web     - 仅编译前端"
	@echo "  build-server  - 仅编译后端"
	@echo "  clean         - Clean all build artifacts"
	@echo "  stop          - Clean running processes and build artifacts"
	@echo "  prepare       - Prepare database environment"
	@echo "  revert        - Revert prepare actions (stop and remove containers)"
	@echo "  all           - Clean all processes and run all components"
	@echo ""
	@echo "Docker Targets (前端+后端 统一镜像):"
	@echo "  docker-build  - 一键构建完整Docker镜像 (自动编译前端和后端)"
	@echo "  docker-run    - 运行完整Docker容器 (前端页面+后端API同时启动)"
	@echo "  docker-stop   - 停止并删除Docker容器"
	@echo "  docker-clean  - 清理: 停止容器+删除容器+删除镜像"
	@echo ""
	@echo "Component Targets:"
	@echo "  server        - Clean server process, build and run server"
	@echo "  web           - Clean web process, build and run web frontend"
	@echo ""
	@echo "Build Modes:"
	@echo "  VER=debug     - Debug mode (default)"
	@echo "  VER=release   - Release mode"
	@echo ""
	@echo "Examples:"
	@echo "  make all              # 本地启动所有组件 (开发调试)"
	@echo "  make server           # 本地只启动后端"
	@echo "  make web              # 本地只启动前端"
	@echo "  make stop             # 停止所有进程"
	@echo "  make prepare          # 准备数据库环境"
	@echo ""
	@echo "  make build            # 编译所有产物 (前端+后端)"
	@echo "  make build-web        # 只编译前端"
	@echo "  make build-server     # 只编译后端"
	@echo ""
	@echo "  make docker-build     # 构建Docker镜像 (自动先编译所有产物)"
	@echo "  make docker-run       # 运行Docker容器"
	@echo "  make docker-stop      # 停止Docker容器"
	@echo "  make docker-clean     # 停止+删除容器 + 删除镜像"

# Clean running processes
stop:
	@echo "Cleaning running processes..."
	@# Stop server process
	@ps aux | grep lpcenter_server | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Stop web frontend process (npm run dev)
	@ps aux | grep "npm run dev" | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Stop all lpcenter-related processes
	@ps aux | grep lpcenter | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@echo "Process cleanup completed!"

# Prepare database environment
prepare:
	@echo "Starting to prepare database environment..."
	@# Check if MySQL is running
	@if ! sudo docker ps | grep -q mysql; then \
		echo "MySQL container not running. Starting MySQL..."; \
		sudo docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=lpcenter -p 3306:3306 -d mysql:latest; \
		echo "Waiting for MySQL to be ready..."; \
		sleep 10; \
	fi
	@echo "Database environment preparation completed!"

# Revert prepare actions (stop and remove containers)
revert:
	@echo "Starting to revert prepare actions..."
	@echo "Stopping and removing mysql container..."
	@sudo docker stop mysql 2>/dev/null || true
	@sudo docker rm mysql 2>/dev/null || true
	@echo "Containers stopped and removed, environment reverted!"

# Run all components
all:
	@echo "Starting to run all components..."
	@echo "Note: This command will run all components in the background"
	@# Clean all processes
	@make stop
	@# Build and run all components
	@make server
	@sleep 2
	@make web
	@echo "All components started!"

# Server related target
server:
	@echo "Starting to build and run server..."
	@# Clean server process
	@ps aux | grep lpcenter_server | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Build server
	@cd ./server && mkdir -p bin && go build -ldflags $(LDFLAGS) -o bin/lpcenter_server ./cmd/main.go
	@echo "Server build completed!"
	@# Run server
	@cd ./server && ./bin/lpcenter_server &
	@echo "Server started!"

# Web frontend related target
web:
	@echo "Starting to build and run web frontend..."
	@# Clean web frontend process
	@ps aux | grep "npm run dev" | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Build web frontend
	@cd ./web && npm install && npm run build
	@echo "Web frontend build completed!"
	@# Run web frontend
	@cd ./web && npm run dev &
	@echo "Web frontend started!"

# Build targets

# Build all components in release mode
release:
	@echo "Building all components in release mode..."
	@make build-server
	@echo "All components built successfully!"

# Build server in release mode
build-server:
	@echo "Building server in release mode..."
	@cd ./server && mkdir -p bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o bin/lpcenter_server ./cmd/main.go
	@echo "Server built successfully!"

# Build frontend
build-web:
	@echo "Building frontend..."
	@cd ./web && npm run build
	@echo "Frontend built successfully!"

# Build all components (frontend + backend)
build:
	@echo "======================================"
	@echo "  编译所有产物"
	@echo "======================================"
	@make build-web
	@make build-server
	@echo ""
	@echo "✅ 所有产物编译完成:"
	@echo "   前端: web/dist/"
	@echo "   后端: server/bin/lpcenter_server"
	@echo ""

# Clean build artifacts
clean:
	@echo "Cleaning all build artifacts..."
	@rm -rf ./server/bin 2>/dev/null || true
	@rm -rf ./web/dist 2>/dev/null || true
	@echo "Cleanup completed!"

docker-build:
	@echo "======================================"
	@echo "  构建 lpcenter 统一镜像 (前端+后端)"
	@echo "======================================"
	@echo "→ 第一步: 本地编译所有产物..."
	@make build
	@echo ""
	@echo "→ 第二步: 构建Docker镜像..."
	@sudo docker build -t lpcenter -f Dockerfile .
	@echo ""
	@echo "✅ Docker镜像构建完成: lpcenter:latest"
	@echo "   镜像包含: 前端页面 + 后端服务 + Nginx"
	@echo ""

# Run unified Docker container
docker-run:
	@make docker-build
	@echo "======================================"
	@echo "  启动 lpcenter 统一容器"
	@echo "======================================"
	@# Stop and remove existing container if it exists
	@sudo docker stop lpcenter 2>/dev/null || true
	@sudo docker rm lpcenter 2>/dev/null || true
	@# Run new container (使用host网络模式，直接访问宿主机MySQL)
	@sudo docker run --name lpcenter \
		--restart always \
		--network host \
		-v lpcenter_uploads:/app/uploads/plugins \
		-d lpcenter
	@echo ""
	@echo "✅ 容器启动成功！"
	@echo "   前端页面: http://localhost:9090"
	@echo "   API接口: http://localhost:9090/api"
	@echo ""

# Stop and remove unified Docker container
docker-stop:
	@echo "停止并删除 lpcenter 容器..."
	@sudo docker stop lpcenter 2>/dev/null || true
	@sudo docker rm lpcenter 2>/dev/null || true
	@echo "✅ 容器已停止并删除！"

# Clean Docker: stop + remove container + remove image
docker-clean:
	@echo "======================================"
	@echo "  清理 lpcenter Docker 环境"
	@echo "======================================"
	@echo "→ 停止并删除容器..."
	@sudo docker stop lpcenter 2>/dev/null || true
	@sudo docker rm lpcenter 2>/dev/null || true
	@echo "→ 删除镜像..."
	@sudo docker rmi lpcenter:latest 2>/dev/null || true
	@echo ""
	@echo "✅ Docker 清理完成！"
	@echo "   容器: 已停止并删除"
	@echo "   镜像: 已删除"
	@echo ""