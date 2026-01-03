FROM gcr.io/distroless/base-debian12:nonroot

# 应用目录
WORKDIR /app

# 拷贝 CI 构建产物
COPY ci/dist/go-demo /app/go-demo

# 端口（仅文档作用）
EXPOSE 8080

# 非 root 用户运行
USER nonroot:nonroot

# 启动
ENTRYPOINT ["/app/go-demo"]
