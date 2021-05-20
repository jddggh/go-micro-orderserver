# 使用Go micro + Gin + Consul的简单微服务demo

- 主要展示服务注册与服务发现功能。

## consul启动命令

```bash
./consul agent -dev -ui -node=consul-dev -client=0.0.0.0
```

## 编译

```bash
go build -mod=vendor -o orderserver main.go
```
