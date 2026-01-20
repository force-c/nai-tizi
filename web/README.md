# NAI-TIZI 管理后台

基于 Vue 3 + TypeScript + Ant Design Vue 构建的现代化企业级管理后台系统。

---

## 📋 目录

- [技术栈](#技术栈)
- [核心特性](#核心特性)
- [快速开始](#快速开始)
- [环境配置](#环境配置)
- [开发指南](#开发指南)
- [构建部署](#构建部署)
- [项目结构](#项目结构)
- [常见问题](#常见问题)

---

## 🛠 技术栈

- **框架：** Vue 3.4+ (Composition API)
- **语言：** TypeScript 5.0+
- **构建工具：** Vite 5.0+
- **UI 组件库：** Ant Design Vue 4.x
- **状态管理：** Pinia + 持久化插件
- **路由：** Vue Router 4.x (动态路由)
- **HTTP 客户端：** Axios
- **代码规范：** ESLint + Prettier

---

## ✨ 核心特性

- ✅ **动态路由系统** - 支持后台配置菜单生成路由
- ✅ **权限控制** - 路由级 + 按钮级权限控制
- ✅ **用户认证** - 密码登录、Token 自动刷新
- ✅ **完整的业务页面** - 用户、角色、菜单、存储、文件管理
- ✅ **响应式布局** - 适配不同屏幕尺寸
- ✅ **TypeScript** - 完整的类型定义

---

## 🚀 快速开始

### 环境要求

- **Node.js:** >= 18.0.0
- **包管理器:** pnpm >= 8.0.0 (推荐) 或 npm >= 9.0.0

### 安装依赖

```bash
# 使用 pnpm (推荐)
pnpm install

# 或使用 npm
npm install

# 或使用 yarn
yarn install
```

### 启动开发服务器

```bash
# 使用 pnpm
pnpm dev

# 或使用 npm
npm run dev

# 或使用 yarn
yarn dev
```

启动成功后，访问：**http://localhost:3000**

### 默认登录账号

确保后端服务已启动（默认端口 8080），使用以下账号登录：

- **用户名：** admin
- **密码：** admin123

---

## ⚙️ 环境配置

### 环境变量文件

项目包含三个环境配置文件：

- `.env` - 所有环境的默认配置
- `.env.development` - 开发环境配置
- `.env.production` - 生产环境配置

### 配置说明

编辑 `.env.development` 文件：

```env
# 应用标题
VITE_APP_TITLE=NAI-TIZI 管理后台 (开发)

# API 基础地址（确保与后端服务地址一致）
VITE_API_BASE_URL=http://localhost:8080

# 客户端认证信息（需要与后端配置一致）
VITE_CLIENT_KEY=web_admin
VITE_CLIENT_SECRET=web_admin_secret

# 是否开启 Mock
VITE_USE_MOCK=false
```

### 后端配置

确保后端数据库中已创建对应的客户端配置：

```sql
INSERT INTO s_auth_client (
  client_id, 
  client_key, 
  client_secret, 
  grant_types, 
  device_type, 
  timeout, 
  active_timeout, 
  share_token, 
  status
)
VALUES (
  'web_admin', 
  'web_admin', 
  'web_admin_secret', 
  'password,email,refresh', 
  'web', 
  7200, 
  1800, 
  false, 
  '0'
);
```

---

## 💻 开发指南

### 代码检查

```bash
# ESLint 检查
pnpm lint

# 代码格式化
pnpm format
```

### 开发规范

- **文件命名：** kebab-case（user-list.vue）
- **组件命名：** PascalCase（UserList）
- **变量命名：** camelCase（userName）
- **常量命名：** UPPER_SNAKE_CASE（API_BASE_URL）

### 新增页面

基于动态路由系统，新增页面只需 3 步：

1. **后端配置菜单**（SQL）
2. **创建 Vue 组件**
3. **重新登录测试**

详见：[动态路由使用示例](../docs/07-前端开发/动态路由使用示例.md)

---

## 📦 构建部署

### 构建生产版本

```bash
# 使用 pnpm
pnpm build

# 或使用 npm
npm run build

# 或使用 yarn
yarn build
```

构建产物在 `dist/` 目录下。

### 预览生产构建

```bash
pnpm preview
```

### 部署方式

#### 方式一：Nginx 部署

1. **构建项目**

```bash
pnpm build
```

2. **配置 Nginx**

创建 Nginx 配置文件 `/etc/nginx/conf.d/nai-tizi.conf`：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    root /var/www/nai-tizi/dist;
    index index.html;
    
    # 前端路由配置（SPA）
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

3. **部署文件**

```bash
# 上传构建产物到服务器
scp -r dist/* user@server:/var/www/nai-tizi/dist/

# 重启 Nginx
sudo nginx -t
sudo systemctl reload nginx
```

#### 方式二：Docker 部署

1. **创建 Dockerfile**

在 `web/` 目录下创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM node:18-alpine as builder

WORKDIR /app

# 复制 package 文件
COPY package*.json ./
COPY pnpm-lock.yaml ./

# 安装 pnpm
RUN npm install -g pnpm

# 安装依赖
RUN pnpm install

# 复制源代码
COPY . .

# 构建
RUN pnpm build

# 生产阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

2. **创建 nginx.conf**

在 `web/` 目录下创建 `nginx.conf`：

```nginx
server {
    listen 80;
    server_name localhost;
    
    root /usr/share/nginx/html;
    index index.html;
    
    # SPA 路由配置
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
}
```

3. **构建和运行**

```bash
# 构建镜像
docker build -t nai-tizi-web .

# 运行容器
docker run -d -p 3000:80 --name nai-tizi-web nai-tizi-web
```

#### 方式三：Docker Compose 部署

创建 `docker-compose.yml`（在项目根目录）：

```yaml
version: '3.8'

services:
  # 后端服务
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
    depends_on:
      - postgres
    networks:
      - nai-tizi-network

  # 前端服务
  frontend:
    build: ./web
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - nai-tizi-network

  # 数据库
  postgres:
    image: postgres:14-alpine
    environment:
      - POSTGRES_DB=nai_tizi
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=your_password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - nai-tizi-network

networks:
  nai-tizi-network:
    driver: bridge

volumes:
  postgres-data:
```

启动所有服务：

```bash
docker-compose up -d
```

---

## 📁 项目结构

```
web/
├── public/                # 静态资源
├── src/
│   ├── api/              # API 接口定义
│   │   ├── auth/         # 认证 API
│   │   ├── user/         # 用户管理 API
│   │   ├── role/         # 角色管理 API
│   │   ├── menu/         # 菜单管理 API
│   │   ├── organization/ # 组织管理 API
│   │   └── storage/      # 存储管理 API
│   ├── assets/           # 资源文件
│   │   └── styles/       # 全局样式
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   │   ├── default/      # 默认布局
│   │   └── blank/        # 空白布局
│   ├── router/           # 路由配置
│   │   └── guards/       # 路由守卫
│   ├── stores/           # 状态管理
│   │   ├── auth.ts       # 认证状态
│   │   ├── user.ts       # 用户状态
│   │   ├── app.ts        # 应用状态
│   │   └── permission.ts # 权限路由
│   ├── types/            # TypeScript 类型
│   │   ├── api.d.ts      # API 类型
│   │   ├── menu.d.ts     # 菜单类型
│   │   ├── system.d.ts   # 系统类型
│   │   └── global.d.ts   # 全局类型
│   ├── utils/            # 工具函数
│   │   └── request.ts    # HTTP 请求封装
│   ├── views/            # 页面组件
│   │   ├── auth/         # 认证页面
│   │   ├── dashboard/    # 仪表盘
│   │   ├── system/       # 系统管理
│   │   │   ├── user/     # 用户管理
│   │   │   ├── role/     # 角色管理
│   │   │   └── menu/     # 菜单管理
│   │   ├── storage/      # 存储管理
│   │   │   ├── env/      # 存储环境
│   │   │   └── file/     # 文件管理
│   │   └── error/        # 错误页面
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── .env                  # 环境变量
├── .env.development      # 开发环境
├── .env.production       # 生产环境
├── index.html            # HTML 模板
├── package.json          # 项目依赖
├── vite.config.ts        # Vite 配置
├── tsconfig.json         # TypeScript 配置
└── README.md             # 项目说明
```

---

## ❓ 常见问题

### Q1: 依赖安装失败？

**A:** 尝试以下方法：

```bash
# 删除依赖和锁文件
rm -rf node_modules package-lock.json

# 清除缓存
npm cache clean --force

# 重新安装
npm install
```

### Q2: 启动后无法访问后端 API？

**A:** 检查以下配置：

1. 后端服务是否启动（默认端口 8080）
2. `.env.development` 中的 `VITE_API_BASE_URL` 是否正确
3. 浏览器控制台是否有 CORS 错误
4. 后端是否配置了 CORS 允许前端域名

### Q3: 登录后 Token 过期？

**A:** 

- Token 会自动刷新，如果刷新失败会跳转到登录页
- 检查后端 RefreshToken 配置是否正确
- 检查客户端认证信息是否匹配

### Q4: 页面样式错乱？

**A:**

- 确保已正确导入 Ant Design Vue 样式
- 检查 `src/main.ts` 中是否导入了 `ant-design-vue/dist/reset.css`
- 清除浏览器缓存后重试

### Q5: 构建后部署访问 404？

**A:**

- 确保 Nginx 配置了 SPA 路由规则：`try_files $uri $uri/ /index.html;`
- 检查 `vite.config.ts` 中的 `base` 配置
- 确保静态资源路径正确

### Q6: 如何调试？

**A:**

1. 使用 Vue DevTools 浏览器扩展
2. 在代码中使用 `console.log()` 或 `debugger`
3. 查看浏览器控制台的网络请求
4. 使用 Vite 的 HMR 功能进行热更新调试

---

## 📚 相关文档

- [前端快速开始](../docs/07-前端开发/前端快速开始.md)
- [前端工程设计方案](../docs/02-架构设计/前端工程设计方案.md)
- [前端动态路由设计方案](../docs/02-架构设计/前端动态路由设计方案.md)
- [动态路由使用示例](../docs/07-前端开发/动态路由使用示例.md)
- [页面开发指南](../docs/07-前端开发/页面开发指南.md)
- [前端页面清单](../docs/07-前端开发/前端页面清单.md)

---

## 🔗 相关链接

- [Vue 3 文档](https://cn.vuejs.org/)
- [Ant Design Vue 文档](https://antdv.com/)
- [Vite 文档](https://cn.vitejs.dev/)
- [TypeScript 文档](https://www.typescriptlang.org/)
- [Pinia 文档](https://pinia.vuejs.org/zh/)

---

## 📄 License

MIT

---

**最后更新：** 2024-12-21  
**项目状态：** ✅ 已完成，可直接使用
