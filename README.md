# 插件授权中心

一个基于Go后端和Vue前端的插件授权申请和发布系统。

## 技术栈

- **后端**: Go + Gin + GORM + MySQL
- **前端**: Vue 3 + Ant Design Vue + Vue Router + Axios

## 项目结构

```
lpcenter/
├── server/                 # 后端Go项目
│   ├── config/            # 配置文件
│   ├── models/            # 数据模型
│   ├── database/          # 数据库连接
│   ├── handlers/          # API处理器
│   ├── routes/            # 路由配置
│   ├── main.go            # 主入口文件
│   └── go.mod             # Go模块依赖
├── web/                   # 前端Vue项目
│   ├── src/
│   │   ├── api/           # API调用
│   │   ├── router/        # 路由配置
│   │   ├── views/         # 页面组件
│   │   ├── App.vue        # 主组件
│   │   └── main.js        # 入口文件
│   └── package.json       # npm依赖
├── database.sql           # 数据库初始化脚本
├── start-server.sh        # 启动后端脚本
└── start-web.sh           # 启动前端脚本
```

## 功能特性

### 后端API

- **插件管理**
  - 创建插件
  - 获取插件列表
  - 获取单个插件详情
  - 更新插件信息
  - 删除插件

- **授权管理**
  - 创建授权申请
  - 获取授权申请列表
  - 获取单个授权申请详情
  - 批准授权申请
  - 拒绝授权申请
  - 获取用户的授权列表

### 前端页面

- **插件列表**: 展示所有可用插件，支持发布新插件、查看详情、下载和删除
- **插件详情**: 显示插件的详细信息，支持下载和申请授权
- **申请授权**: 提交插件授权申请
- **授权管理**: 管理所有授权申请，支持批准和拒绝操作

## 快速开始

### 1. 数据库设置

首先创建MySQL数据库并导入表结构：

```bash
mysql -u root -p < database.sql
```

或者手动执行SQL脚本中的内容。

### 2. 配置后端

复制环境变量配置文件：

```bash
cd server
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接信息：

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=lpcenter
SERVER_PORT=8080
```

### 3. 启动后端

```bash
chmod +x start-server.sh
./start-server.sh
```

或者手动启动：

```bash
cd server
export GOPROXY=https://proxy.golang.org,direct
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 4. 启动前端

```bash
chmod +x start-web.sh
./start-web.sh
```

或者手动启动：

```bash
cd web
npm install
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

## API文档

### 插件API

- `POST /api/plugins` - 创建插件
- `GET /api/plugins` - 获取插件列表
- `GET /api/plugins/:id` - 获取插件详情
- `PUT /api/plugins/:id` - 更新插件
- `DELETE /api/plugins/:id` - 删除插件

### 授权API

- `POST /api/licenses` - 创建授权申请
- `GET /api/licenses` - 获取授权申请列表
- `GET /api/licenses/:id` - 获取授权申请详情
- `PUT /api/licenses/:id/approve` - 批准授权
- `PUT /api/licenses/:id/reject` - 拒绝授权
- `GET /api/licenses/user/:userId` - 获取用户的授权列表

## 数据库表结构

### users表
- id: 用户ID
- username: 用户名
- email: 邮箱
- password: 密码
- created_at: 创建时间
- updated_at: 更新时间

### plugins表
- id: 插件ID
- name: 插件名称
- version: 版本号
- description: 描述
- author: 作者
- file_path: 文件路径
- download_url: 下载链接
- created_at: 创建时间
- updated_at: 更新时间

### licenses表
- id: 授权ID
- user_id: 用户ID
- plugin_id: 插件ID
- status: 状态 (pending/approved/rejected)
- reason: 申请理由
- created_at: 创建时间
- updated_at: 更新时间

## 开发说明

### 后端开发

- 后端使用Gin框架提供RESTful API
- 使用GORM进行数据库操作
- 支持CORS跨域请求
- 数据库表会自动创建和更新

### 前端开发

- 使用Vue 3 Composition API
- 使用Ant Design Vue组件库
- 使用Vue Router进行路由管理
- 使用Axios进行HTTP请求

## 注意事项

1. 确保MySQL服务已启动并可访问
2. 确保数据库用户有足够的权限
3. 前端默认连接到 `http://localhost:8080/api`
4. 如需修改后端地址，请编辑 `web/src/api/index.js` 中的 `baseURL`
