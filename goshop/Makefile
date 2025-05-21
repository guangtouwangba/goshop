.PHONY: build clean test proto run deps lint docker-up docker-down

# 服务列表
SERVICES := user product inventory order payment marketing cms shipping gateway auth admin

# Build settings
GOBIN := $(shell go env GOPATH)/bin
SHELL=/bin/bash

# Build all services
build: 
	@for service in $(SERVICES); do \
		echo "Building $$service service..." ; \
		go build -o ./build/$$service ./services/$$service/cmd ; \
	done

# Clean build artifacts
clean:
	@rm -rf ./build/*

# Run tests
test:
	@go test ./... -v

# Run tests with coverage
test-cover:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

# Generate protocol buffers
proto:
	@for service in $(SERVICES); do \
		echo "Generating protobuf for $$service service..." ; \
		protoc --go_out=. --go-grpc_out=. ./api/proto/$$service/*.proto ; \
	done

# Install development dependencies
deps:
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	@golangci-lint run

# Start services using docker-compose
docker-up:
	@docker-compose up -d

# Stop services
docker-down:
	@docker-compose down

# Help command
help:
	@echo "GoShop 可用命令:"
	@echo "  make build         - 构建所有服务"
	@echo "  make clean         - 清理构建产物"
	@echo "  make test          - 运行测试"
	@echo "  make test-cover    - 运行测试并生成覆盖率报告"
	@echo "  make proto         - 生成 Protocol Buffers 代码"
	@echo "  make deps          - 安装开发依赖"
	@echo "  make lint          - 运行代码检查"
	@echo "  make docker-up     - 启动 Docker 服务"
	@echo "  make docker-down   - 停止 Docker 服务"

# Default target
default: help
