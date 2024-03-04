# 使用官方Go镜像作为构建环境
FROM golang:1.22 as builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件并下载依赖信息
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY src/ .

# 编译应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demo .

# 使用scratch作为运行环境
FROM scratch

# 从构建器中复制编译好的应用程序
COPY --from=builder /app/demo /demo

# 运行应用程序
CMD ["/demo"]
