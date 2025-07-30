# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

**项目名称：选课社区2.0后端 (jcourse_go)**

## 开发命令

### 构建和运行
- `go fmt ./...` - 格式化 Go 代码
- `goimports -local jcourse_go -w $(find . -type f -name '*.go')` - 格式化导入包
- `go mod tidy` - 清理 Go 模块
- `make lint` - 运行所有代码检查/格式化命令

### 测试
- `go test ./...` - 运行所有测试
- `go test -v ./internal/application/auth/...` - 运行特定包的测试

## 架构概览

这是一个采用领域驱动设计（DDD）模式和清洁架构的 Go 课程评价系统后端。

### 核心层级

**领域层** (`internal/domain/`):
- 包含业务逻辑、实体、值对象和领域服务
- 核心领域：`auth`（认证）、`review`（评价）、`point`（积分）、`permission`（权限）
- 领域实体强制执行业务规则和不变量
- 在领域层定义仓储接口

**应用层** (`internal/application/`):
- 包含用例和应用服务
- 命令用于写操作，查询用于读操作
- 协调领域对象和仓储
- 视图对象用于 API 响应

**基础设施层** (`internal/infrastructure/`):
- 数据库实现、外部服务
- 领域仓储的具体实现

**接口层** (`internal/interface/`):
- HTTP 控制器、中间件和路由
- Web 框架集成（Gin）

### 核心领域

**评价系统**:
- `Review` 实体：包含课程关联和用户所有权
- `Course` 和 `OfferedCourse` 实体：包含教师关系
- `ReviewRevision`：用于审计追踪
- 命令处理器处理写操作，查询处理器处理读操作
- 值对象：`Rating`（评分）、`Semester`（学期）、`Category`（分类）

**认证系统**:
- 用户管理和会话处理
- 邮箱验证码

**积分系统**:
- 用户积分跟踪和管理

### 常用模式

**依赖注入**:
- 服务容器模式在 `internal/app/container.go`
- 仓储接口注入到应用服务中

**错误处理**:
- 自定义错误类型在 `pkg/apperror/`
- 领域特定的错误代码和包装

**值对象**:
- 不可变类型带验证（如 `Rating`、`Semester`）
- 使用业务规则创建的工厂方法

**命令和查询**:
- 写操作（命令）和读操作（查询）分离的处理器
- 用于输入验证的命令 DTO

### 项目结构

```
cmd/                    # 应用程序入口点
  api/                 # HTTP 服务器
  worker/              # 后台工作进程
internal/
  app/                 # 依赖注入容器
  application/         # 用例和应用服务
    auth/              # 认证命令/查询
    review/            # 评价系统命令/查询
    point/             # 积分系统命令/查询
  domain/              # 业务逻辑和实体
    auth/              # 用户和会话实体
    review/            # 课程和评价实体
    point/             # 积分实体
    common/            # 共享领域概念
    event/             # 领域事件
  config/              # 配置结构
  interface/           # 外部接口
    web/               # HTTP 控制器和路由
    middleware/        # HTTP 中间件
pkg/                   # 共享库
  apperror/            # 错误处理工具
  password/            # 密码工具
```

### 配置

- 基于 YAML 的配置，包含数据库、Redis 和 SMTP 设置
- 环境特定配置文件位于 `config/` 目录
- 服务容器管理依赖注入

### 开发说明

- 需要 Go 1.24
- 使用 Gin Web 框架进行 HTTP 路由
- 使用 Testify 进行测试
- 标准 Go 项目布局，`internal/` 用于私有代码
- 空的 `main.go` 文件表明这是一个模板/启动项目