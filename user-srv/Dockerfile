# 基准镜像
FROM golang:1.18.2-alpine AS builder

# 配置
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码
COPY . /go/src/user-service

# 编译
WORKDIR /go/src/user-service
RUN go install main.go

# 生产镜像
FROM alpine:3.15
COPY --from=builder /go/bin/main /bin/main

# 声明暴露端口 
EXPOSE 8080

# 服务入口
ENTRYPOINT ["/bin/main"]

