# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.16-alpine as builder
# 启用go module
ENV GO111MODULE=on\
    GOPROXY=https://goproxy.cn,direct
WORKDIR /wm-motor.com/Infra/httpserver
COPY * .
# 指定OS等，并go build
RUN GOOS=linux GOARCH=amd64 go build . -o httpserver
# 运行阶段指定scratch作为基础镜像
FROM alpine
WORKDIR /app
COPY --from=0 /wm-motor.com/Infra/httpserver/httpserver .
EXPOSE 8090
ENTRYPOINT  ["./httpserver"]
