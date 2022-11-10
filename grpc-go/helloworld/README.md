# helloworld

## 操作步骤

```
// 安装依赖
$ GO111MODULE=on go mod vendor

// 运行 grpc 服务端
$ go run greeter_server/main.go

// 打开另一个 terminal，运行 grpc 客户端，等待 3 秒左右返回结果
$ cd grpc-go/helloworld
$ go run greeter_client/main.go
```

