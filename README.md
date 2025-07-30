# 选课社区2.0后端 (jcourse_go)

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/SJTU-jCourse/jcourse_go)

选课社区2.0的后端服务，采用领域驱动设计（DDD）和清洁架构构建的Go语言课程评价系统。

## 🌟 项目特性

- **领域驱动设计**：采用DDD模式，清晰的分层架构
- **清洁架构**：依赖倒置，易于测试和维护
- **RESTful API**：基于Gin框架的高性能HTTP服务
- **认证授权**：完整的用户认证和权限管理系统
- **课程评价**：支持课程评价、评分、学期管理
- **积分系统**：用户积分跟踪和管理
- **审计追踪**：评价修改历史记录
- **高并发支持**：异步处理和限流机制

## 🛠 技术栈

- **语言**: Go 1.24
- **Web框架**: Gin
- **架构模式**: 领域驱动设计 (DDD)
- **数据库**: MySQL (可配置)
- **缓存**: Redis (可配置)
- **测试**: Testify
- **代码工具**: gofmt, goimports

## 📁 项目结构

```
cmd/                    # 应用程序入口点
  api/                 # HTTP 服务器
  worker/              # 后台工作进程
internal/
  app/                 # 依赖注入容器
  application/         # 应用服务层
    auth/              # 认证服务
    review/            # 评价服务
    point/             # 积分服务
  domain/              # 领域层
    auth/              # 认证领域
    review/            # 评价领域
    point/             # 积分领域
    common/            # 共享领域概念
    event/             # 领域事件
  config/              # 配置管理
  interface/           # 接口层
    web/               # HTTP 控制器
    middleware/        # 中间件
pkg/                   # 公共库
  apperror/            # 错误处理
  password/            # 密码工具
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+
- Redis 6.0+

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/SJTU-jCourse/jcourse_go.git
   cd jcourse_go
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置文件**
   在 `config/` 目录下创建配置文件 `config.yaml`：
   ```yaml
   db:
     dsn: "user:password@tcp(localhost:3306)/jcourse?charset=utf8mb4&parseTime=True&loc=Local"
   redis:
     addr: "localhost"
     port: 6379
     password: ""
     db: 0
   smtp:
     host: "smtp.example.com"
     port: 587
     username: "your-email@example.com"
     password: "your-password"
     sender: "noreply@example.com"
   ```

4. **运行项目**
   ```bash
   # 启动 API 服务
   go run cmd/api/main.go
   
   # 启动后台工作进程
   go run cmd/worker/main.go
   ```

### 开发工具

```bash
# 格式化代码
make lint

# 运行测试
go test ./...

# 运行特定测试
go test -v ./internal/application/auth/...
```

## 📖 API 文档

### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/verify` - 邮箱验证

### 课程评价
- `GET /api/courses` - 获取课程列表
- `GET /api/courses/:id` - 获取课程详情
- `POST /api/reviews` - 发布评价
- `PUT /api/reviews/:id` - 更新评价
- `DELETE /api/reviews/:id` - 删除评价

### 积分系统
- `GET /api/points` - 获取积分记录
- `POST /api/points/earn` - 获得积分

## 🤝 贡献指南

我们欢迎任何形式的贡献！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范

- 遵循 Go 语言标准代码风格
- 编写单元测试覆盖业务逻辑
- 使用 DDD 模式进行领域建模
- 保持清晰的分层架构

## 📝 许可证

本项目采用 APGLv3 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为选课社区2.0项目做出贡献的开发者。

## 📞 联系我们

- 项目地址: [https://github.com/SJTU-jCourse/jcourse_go](https://github.com/SJTU-jCourse/jcourse_go)
- 问题反馈: [Issues](https://github.com/SJTU-jCourse/jcourse_go/issues)

---

⭐ 如果这个项目对你有帮助，请给个 star！