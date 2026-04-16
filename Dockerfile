# ITFeeds Dockerfile
# 多阶段构建：前端 → 后端 → 运行镜像

# 阶段1：前端构建
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

COPY web/package.json web/package-lock.json* ./

RUN npm install

COPY web/ ./

RUN npm run build

# 阶段2：后端构建
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git gcc g++

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./resource/public

RUN go build -o main main.go

# 阶段3：最终运行镜像
FROM alpine

LABEL maintainer="ITFeeds"

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=backend-builder /app/main ./
COPY --from=backend-builder /app/manifest/config/config.yaml ./manifest/config/
COPY --from=backend-builder /app/resource/public ./resource/public
COPY --from=backend-builder /app/resource/sql ./resource/sql

RUN chmod +x ./main

EXPOSE 8000

CMD ["./main"]
