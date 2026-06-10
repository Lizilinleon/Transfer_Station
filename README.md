# AI 智能对话平台

一个分阶段开发的 AI 平台项目，当前已完成：

- 前端静态控制台页面
- Go + Gin 后端基础骨架
- SQLite 数据库初始化与默认数据写入
- 模型管理 CRUD
- 密钥生成与密钥管理 CRUD
- 系统配置项 CRUD

当前版本优先完成“管理后台基础能力”，暂未接入真实 LLM 对话 API。

## 技术栈

### 前端
- React JS
- Vite
- 普通 CSS

### 后端
- Go
- Gin
- GORM
- SQLite

## 当前目录结构

```text
.
├── frontend/          # React + Vite 前端
├── backend/           # Go + Gin 后端
├── docs/              # 项目文档
├── .env.example       # 环境变量示例
├── README.md
└── AGENTS.md
```

## 环境变量

在项目根目录创建 `.env` 文件，可参考 `.env.example`：

```env
APP_NAME=AI Smart Chat Platform
APP_ENV=development
BACKEND_PORT=8080
SQLITE_PATH=./data/app.db
CORS_ALLOW_ORIGIN=http://localhost:5173

DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_PASSWORD=admin123456

DEFAULT_MODEL_NAME=gpt-5.5
DEFAULT_MODEL_PROVIDER=openai-compatible
DEFAULT_MODEL_GROUP=codex专用分组
```

## 如何运行

### 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认访问：

```text
http://localhost:5173/app/keys
```

### 启动后端

```bash
cd backend
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=off
go mod tidy
go run .
```

如果你在 PowerShell 中运行，使用：

```powershell
$env:GOPROXY="https://goproxy.cn,direct"
$env:GOSUMDB="off"
go mod tidy
go run .
```

后端默认地址：

```text
http://localhost:8080
```

## 当前后端接口

### 健康检查
- `GET /health`

### 模型管理
- `GET /api/models`
- `GET /api/models/:id`
- `POST /api/models`
- `PUT /api/models/:id`
- `DELETE /api/models/:id`

### 密钥管理
- `GET /api/keys`
- `GET /api/keys/:id`
- `POST /api/keys`
- `PUT /api/keys/:id`
- `DELETE /api/keys/:id`

### 系统配置
- `GET /api/options`
- `POST /api/options`
- `PUT /api/options/:id`
- `DELETE /api/options/:id`

## 默认初始化内容

项目首次启动后端时会自动创建：

- SQLite 数据库文件
- 默认管理员账户
- 默认模型 `gpt-5.5`
- 默认分组 `codex专用分组`
- 默认系统配置项

说明：当前版本还没有登录接口，默认管理员主要用于后续扩展。

## 文档

- [开发规则](./AGENTS.md)
- [需求文档](./docs/PRD.md)
- [后端接口文档](./docs/backend-api.md)
