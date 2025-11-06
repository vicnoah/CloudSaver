# CloudSaver Go 版本 - 快速启动指南

## 🚀 快速开始

### 1. 编译程序

```bash
cd /data/workspace/CloudSaver
go build -tags modernc -o bin/cloudsaver cmd/server/main.go
```

### 2. 运行程序

```bash
# 默认配置运行 (端口 8080)
./bin/cloudsaver

# 或使用自定义端口
CLOUDSAVER_SERVER_PORT=8888 ./bin/cloudsaver
```

### 3. 测试 API

#### 注册用户
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123",
    "registerCode": 9527
  }'
```

**默认注册码**: 
- 普通用户: `9527`
- 管理员: `230713`

#### 登录
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

登录成功后会返回 Token，保存它以供后续请求使用。

#### 获取用户设置
```bash
curl -X GET http://localhost:8080/api/setting/get \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### 保存云盘 Cookie
```bash
curl -X POST http://localhost:8080/api/setting/save \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "cloud115_cookie": "您的115网盘Cookie",
    "quark_cookie": "您的夸克网盘Cookie"
  }'
```

## 📋 API 接口列表

### 用户接口 (无需认证)
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录

### 设置接口 (需要认证)
- `GET /api/setting/get` - 获取用户设置
- `POST /api/setting/save` - 保存用户设置

### 115 网盘接口 (需要认证)
- `GET /api/cloud115/share-info` - 获取分享信息
  - 参数: `share_code`, `passcode`
- `GET /api/cloud115/folders` - 获取文件夹列表
  - 参数: `parent_cid` (可选，默认为根目录)
- `POST /api/cloud115/save` - 转存文件
  - 参数: `share_code`, `passcode`, `folder_id`, `file_ids`

### 夸克网盘接口 (需要认证)
- `GET /api/quark/share-info` - 获取分享信息
  - 参数: `share_code` (pwd_id), `passcode`
- `GET /api/quark/folders` - 获取文件夹列表
  - 参数: `parent_cid` (可选，默认为根目录)
- `POST /api/quark/save` - 转存文件
  - 参数: `pwd_id`, `stoken`, `folder_id`, `file_ids`, `file_tokens`

## 🔧 配置

### 环境变量
- `CLOUDSAVER_SERVER_PORT` - 服务器端口 (默认: 8080)
- `CLOUDSAVER_SERVER_MODE` - 运行模式 (默认: release)
- `CLOUDSAVER_DATABASE_PATH` - 数据库路径 (默认: ./data/cloudsaver.db)

### 配置文件 (可选)
创建 `config.yaml` 文件：

```yaml
server:
  port: 8080
  mode: release

database:
  path: ./data/cloudsaver.db

jwt:
  secret: your-secret-key-change-in-production
  expire_hours: 168  # 7天
```

## 📁 目录结构

```
CloudSaver/
├── bin/                    # 编译输出
│   └── cloudsaver         # 可执行文件
├── cmd/
│   └── server/
│       └── main.go        # 程序入口
├── configs/               # 配置文件目录
├── data/                  # 数据目录
│   └── cloudsaver.db     # SQLite 数据库
├── embed/                 # 前端资源嵌入
│   └── dist/             # 前端构建产物
├── internal/
│   ├── api/
│   │   ├── handler/      # HTTP 处理器
│   │   ├── middleware/   # 中间件
│   │   ├── request/      # 请求结构
│   │   └── response/     # 响应结构
│   ├── config/           # 配置管理
│   ├── model/            # 数据模型
│   ├── pkg/
│   │   └── cloud/        # 云盘服务
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
└── web/                  # Vue 前端项目

```

## 🔑 获取云盘 Cookie

### 115 网盘
1. 打开浏览器开发者工具 (F12)
2. 访问 https://115.com
3. 登录账号
4. 在 Network 标签页中找到任意请求
5. 复制 Request Headers 中的 Cookie 值

### 夸克网盘
1. 打开浏览器开发者工具 (F12)
2. 访问 https://pan.quark.cn
3. 登录账号
4. 在 Network 标签页中找到任意请求
5. 复制 Request Headers 中的 Cookie 值

## ⚠️ 注意事项

1. **Cookie 安全**: 请妥善保管您的云盘 Cookie，不要泄露给他人
2. **注册码**: 首次部署后请修改默认注册码
3. **JWT Secret**: 生产环境请修改默认的 JWT Secret
4. **端口冲突**: 如果 8080 端口被占用，请使用环境变量指定其他端口

## 🎯 下一步

1. 构建 Vue 前端项目
2. 将前端产物复制到 `embed/dist/` 目录
3. 重新编译包含前端的完整版本
4. 通过浏览器访问 Web 界面

## 📞 技术支持

如遇到问题，请检查：
1. Go 版本是否 >= 1.22
2. 数据库文件是否有写入权限
3. 端口是否被占用
4. Cookie 是否有效

---

**版本**: 1.0.0  
**构建时间**: 2025-11-06  
**Go 版本**: 1.22+
