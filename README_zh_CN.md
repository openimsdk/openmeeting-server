
# 快速实现 API/RPC 教程

本教程将指导您以 `user` 为例，逐步解释初始化到 API 和 RPC 的实现过程进行。

## 1. 初始化模块

首先，初始化您的 Go 模块：

```bash
go mod init github.com/your-account/your-project
```

## 2. 定义 Protobuf

在 `pkg/protocol` 目录下定义您的协议。通常，我们按照一个 RPC 一个目录的原则操作。创建 `user` 目录及 `user.proto` 协议文件：

```plaintext
pkg/protocol/user/user.proto
```

## 3. 生成 Protobuf 代码

在项目根目录执行以下命令，以生成 `user.pb.go` 代码：

```bash
mage protocol
```

**注意：** 如果是首次使用 `mage`，在 Linux/Mac 平台下执行 `bash bootstrap.sh`，在 Windows 执行 `bootstrap.bat`。

## 4. 实现存储逻辑

存储逻辑位于 `pkg/common/storage`，其中包含 `cache`, `controller`, `database`, `model` 目录：

- **Model**: 在 `model` 目录定义结构体。
- **Controller**: 在 `controller` 目录定义接口。
- **Cache/Database**: 在 `cache` 和 `database` 目录定义接口，并在 `redis` 和 `mgo` 目录下实现这些接口，以完成具体的业务逻辑。

## 5. 实现 RPC

在 `internal/rpc/user/user.go` 中实现 `Start` 函数及具体的 RPC 函数。

## 6. 实现 API

- 在 `internal/api/router.go` 增加路由。对于 whitelist 的路由不验证 token，其他 API 请求会进行 token 验证。
- 在 `internal/api/user.go` 完成 API 调用 RPC 的逻辑。

## 7. 入口函数

- **API 入口**: 在 `pkg/cmd/api.go`，具体的 `start` 函数在 `internal/api/init.go`。
- **RPC 入口**: 在 `pkg/cmd/user.go`，具体的 `start` 函数在 `pkg/startrpc/start.go`。

## 8. 关于鉴权、日志、OperationID、Token 和 Context

- **OperationID**: 在 RPC 中从 `context` 获取 `OperationID` 使用 `mcontext.GetOperationID(ctx)`。
- **UserID 和 PlatformID**: 在 RPC 中从 `context` 获取登录用户 `userID` 和平台 `platformID` 使用 `mcontext.GetOpUserID(ctx)` 和 `mcontext.GetOpUserPlatform(ctx)`。
- **API 调用者**: 需要将 token 和 operationID 设置在 HTTP header 中。
- **日志打印**: 统一使用"github.com/openimsdk/tools/log" 





