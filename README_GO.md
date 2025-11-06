# CloudSaver - Go重构版

基于 Go + Gin + GORM 的云盘资源管理系统，支持 115网盘、夸克网盘等多个云盘平台的资源搜索和转存。

## 特性

- 🚀 **单一二进制部署** - 前端资源内嵌，无需额外配置
- ⚡ **高性能** - Go协程并发处理，内存占用低
- 🔐 **JWT认证** - 安全的用户认证机制
- 📦 **多云盘支持** - 115网盘、夸克网盘
- 🔍 **资源搜索** - Telegram频道资源聚合搜索
- 🎬 **豆瓣榜单** - 集成豆瓣电影/电视剧榜单
- 📱 **响应式设计** - 支持PC和移动端

## 技术栈

### 后端
- Go 1.22+
- Gin - Web框架
- GORM - ORM框架
- SQLite - 数据库
- JWT - 身份认证

### 前端
- Vue 3 + TypeScript
- Vite - 构建工具
- UnoCSS - 原子化CSS
- Pinia - 状态管理
- Vue Router - 路由

## 快速开始

### 前置要求

- Go 1.22+
- Node.js 18+
- pnpm

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/your/cloudsaver.git
cd cloudsaver

# 使用Makefile构建
make build

# 或使用脚本
chmod +x scripts/build.sh
./scripts/build.sh

# 运行
./bin/cloudsaver
```

### 使用预编译二进制

```bash
# 下载对应平台的二进制文件
# 创建配置文件
cp configs/config.example.yaml config.yaml

# 编辑配置文件（可选）
vim config.yaml

# 运行
./cloudsaver
```

## 配置说明

配置文件 `config.yaml`：

```yaml
server:
  port: 8080              # 服务端口
  mode: release           # 运行模式: debug/release

database:
  path: ./data/cloudsaver.db  # 数据库文件路径

jwt:
  secret: "your-secret-key"   # JWT密钥（生产环境请修改）
  expire_hours: 168           # Token过期时间（小时）
```

### 环境变量

也可以通过环境变量覆盖配置：

```bash
export CLOUDSAVER_SERVER_PORT=8080
export CLOUDSAVER_JWT_SECRET=your-secret-key
./cloudsaver
```

## API 文档

### 用户接口

#### 注册
```
POST /api/user/register
Body: {
  "username": "user",
  "password": "password",
  "registerCode": 9527
}
```

#### 登录
```
POST /api/user/login
Body: {
  "username": "user",
  "password": "password"
}
```

### 云盘接口

#### 获取115分享信息
```
GET /api/cloud115/share-info?shareCode=xxx&passcode=xxx
Header: Authorization: Bearer <token>
```

#### 获取夸克分享信息
```
GET /api/quark/share-info?shareCode=xxx&passcode=xxx
Header: Authorization: Bearer <token>
```

完整API文档请查看：[API Documentation](docs/api.md)

## 开发指南

### 项目结构

```
cloudsaver/
├── cmd/server/           # 程序入口
├── internal/             # 内部代码
│   ├── api/              # API层（Handler、Middleware、DTO）
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── model/            # 数据模型
│   ├── pkg/              # 内部共享包
│   └── config/           # 配置管理
├── web/                  # 前端源码
├── embed/                # 嵌入式资源
├── configs/              # 配置文件
├── migrations/           # 数据库迁移
└── scripts/              # 构建脚本
```

### 开发命令

```bash
# 安装依赖
go mod download

# 仅构建后端（开发模式）
go build -o bin/cloudsaver cmd/server/main.go

# 运行测试
go test -v ./...

# 清理构建产物
make clean
```

### 跨平台编译

```bash
# Linux
make build-linux

# Windows
make build-windows

# macOS (Intel)
make build-darwin-amd64

# macOS (Apple Silicon)
make build-darwin-arm64

# 编译所有平台
make build-all
```

## 部署

### Docker部署

```dockerfile
# Dockerfile已包含在项目中
docker build -t cloudsaver .
docker run -d -p 8080:8080 -v ./data:/app/data cloudsaver
```

### Systemd服务

```ini
[Unit]
Description=CloudSaver Service
After=network.target

[Service]
Type=simple
User=cloudsaver
WorkingDirectory=/opt/cloudsaver
ExecStart=/opt/cloudsaver/cloudsaver
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 注册码说明

系统支持两种注册码：
- `9527` - 普通用户注册码（默认）
- `230713` - 管理员注册码（默认）

可在配置文件中修改注册码。

## 注意事项

1. **首次运行**: 程序会自动创建数据库和必要的目录
2. **JWT密钥**: 生产环境务必修改默认密钥
3. **Cookie配置**: 需在设置页面配置对应云盘的Cookie才能使用转存功能
4. **代理设置**: 如需访问Telegram，请配置代理

## 许可证

MIT License

## 致谢

感谢原始项目的贡献者们。

---

**注意**: 本项目为教育学习目的，请遵守相关平台的服务条款。
