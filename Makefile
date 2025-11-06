.PHONY: help build run clean test build-frontend build-backend

help:
	@echo "CloudSaver Makefile Commands:"
	@echo "  make build          - 构建前后端完整项目"
	@echo "  make build-frontend - 仅构建前端"
	@echo "  make build-backend  - 仅构建后端"
	@echo "  make run            - 运行应用"
	@echo "  make clean          - 清理构建产物"
	@echo "  make test           - 运行测试"

build: build-frontend build-backend

build-frontend:
	@echo "Building frontend..."
	cd web && pnpm install && pnpm build
	@echo "Copying frontend dist to embed..."
	mkdir -p embed/dist
	cp -r web/dist/* embed/dist/

build-backend: build-frontend
	@echo "Building backend..."
	go mod download
	go build -ldflags="-s -w" -o bin/cloudsaver cmd/server/main.go
	@echo "Build completed! Binary: bin/cloudsaver"

run:
	@echo "Running CloudSaver..."
	./bin/cloudsaver

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf embed/dist/
	rm -rf web/dist/
	rm -rf data/

test:
	@echo "Running tests..."
	go test -v ./...

# 跨平台编译
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/cloudsaver-linux-amd64 cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/cloudsaver-windows-amd64.exe cmd/server/main.go

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/cloudsaver-darwin-amd64 cmd/server/main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/cloudsaver-darwin-arm64 cmd/server/main.go

build-all: build-frontend build-linux build-windows build-darwin-amd64 build-darwin-arm64
	@echo "All platforms built successfully!"
