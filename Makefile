BINARY_NAME=vless-reality-proxy
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
LDFLAGS=-ldflags "-X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME}"

.PHONY: build clean test run deps

# 构建
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go

# 交叉编译
build-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}.exe main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin main.go

build-all: build-linux build-windows build-darwin

# 安装依赖
deps:
	go mod download
	go mod tidy

# 运行
run:
	go run main.go

# 测试
test:
	go test -v ./...

# 清理
clean:
	go clean
	rm -f ${BINARY_NAME}*

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...

# 生成配置工具
gen-config:
	@echo "生成示例配置..."
	@go run tools/generate-config.go full-config

# 测试域名
test-domains:
	@echo "测试 Reality 域名..."
	@go run tools/test-domains.go

# 地理位置分析
analyze-geography:
	@echo "分析域名地理位置特征..."
	@go run tools/analyze-geography.go

# 安装到系统
install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/
	sudo chmod +x /usr/local/bin/${BINARY_NAME}

# 卸载
uninstall:
	sudo rm -f /usr/local/bin/${BINARY_NAME}

# 帮助
help:
	@echo "VLESS Reality 代理服务器 - 可用命令:"
	@echo ""
	@echo "构建相关:"
	@echo "  build        - 构建二进制文件"
	@echo "  build-all    - 交叉编译所有平台"
	@echo "  deps         - 安装依赖"
	@echo "  clean        - 清理构建文件"
	@echo ""
	@echo "开发相关:"
	@echo "  run          - 运行程序"
	@echo "  test         - 运行测试"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo ""
	@echo "配置工具:"
	@echo "  gen-config   - 生成完整配置文件"
	@echo "  test-domains - 测试 Reality 回落域名"
	@echo ""
	@echo "部署相关:"
	@echo "  install      - 安装到系统"
	@echo "  uninstall    - 从系统卸载"
	@echo ""
	@echo "帮助:"
	@echo "  help         - 显示此帮助信息"
	@echo ""
	@echo "示例用法:"
	@echo "  make deps && make gen-config && make run"
