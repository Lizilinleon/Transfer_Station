# AI 智能对话平台 - 开发规则

## 技术栈

### 前端
- **框架**：React JS（不使用 TypeScript）
- **构建工具**：Vite
- **样式**：普通 CSS（不使用 Tailwind）

### 后端
- **语言**：Go
- **Web 框架**：Gin

### 数据库
- **第一阶段**：SQLite

## 开发规则

### 1. 功能实现
- ❌ 不要一次性实现所有功能
- ✅ 每次任务只做明确的一小步
- ✅ 每个功能点单独测试验证

### 2. 提交和文档
- 每次修改后需要说明：
  - 改了什么（具体改动内容）
  - 如何运行（相关命令）
  - 如何验证（测试方式）

### 3. 敏感信息管理
- ❌ 不要把 LLM API Key 写死在代码里
- ✅ 只能通过 `.env` 文件读取
- ✅ `.env` 文件添加到 `.gitignore`

### 4. 版本控制
不要提交以下文件/目录：
- `node_modules`
- `dist`
- 数据库文件（`*.db`）
- `.env`
- `backend/tmp`
- `frontend/.vite`

## 项目结构

```
.
├── frontend/          # React + Vite 项目
├── backend/           # Go + Gin 项目
├── docs/              # 文档
├── README.md
├── AGENTS.md
├── .gitignore
└── .env.example       # .env 示例文件
```

## 命名规范

### 前端
- 组件文件：PascalCase（如 `ChatBox.js`）
- 样式文件：kebab-case（如 `chat-box.css`）
- 工具函数：camelCase（如 `formatMessage()`)

### 后端
- 包名：lowercase
- 函数名：CamelCase（exported）
- 变量名：camelCase

## CI/CD 规范

- 提交前确保代码能运行
- 更新相关文档
- 清晰的 Git commit message
