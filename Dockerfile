# ---- 构建阶段 ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用缓存（依赖没变就不重新下载）
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码
COPY . .

# 编译，关闭CGO保证在alpine里能运行
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---- 运行阶段 ----
FROM alpine:latest

WORKDIR /app

# 只从构建阶段拷贝编译好的二进制文件
COPY --from=builder /app/main .
COPY config.yaml .

CMD ["./main"]