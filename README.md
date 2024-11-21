# protoc-gen-gme

    protobuf协议代码自动生成工具，使用protoc-gen-gme生成框架相关的协议代码，来减少手写的代码量

## 安装 

    go get protoc-gen-gme仓库代码

同时需要安装: 

- [protoc](https://github.com/google/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf)

## 使用说明

定义protobuf协议，如下 `greeter.proto`

```
syntax = "proto3";

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

生成代码命令：

```
protoc --proto_path=$GOPATH/src:. --gme_out=. --go_out=. greeter.proto
```

得到的文件列表如下:

```
./
    greeter.proto	# 原始的protobuf协议文件
    greeter.pb.go	# protoc-gen-go 自动生成的标准代码
    greeter.gme.go	# protoc-gen-gme 自动生成的框架代码
```
# protoc-gen-gme
