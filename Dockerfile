# 构建阶段
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# 复制依赖文件（如果有 go.sum 也会一起复制）
COPY go.mod ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（静态链接，便于在 scratch/alpine 运行）
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o iconstation main.go

# 运行阶段
FROM alpine:latest

# 安装 CA 证书（用于 HTTPS 请求）和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/iconstation .

# 数据目录（容器内可执行文件所在目录会自动创建 UserData）
VOLUME ["/app/UserData"]

EXPOSE 9168

ENTRYPOINT ["./iconstation"]
