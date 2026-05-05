FROM alpine:latest
WORKDIR /app

# 复制二进制文件（运行时通过参数指定具体文件）
ARG BIN_FILE
COPY ${BIN_FILE} ./iconstation
COPY static /app/static
COPY icon.png /usr/share/icons/icon.png

# 一次性创建目录 + 装证书时区
RUN mkdir -p /app/UserData/icons /app/UserData/chunks \
    && apk add --no-cache ca-certificates tzdata

LABEL org.opencontainers.image.icon=/usr/share/icons/icon.png

EXPOSE 9168
CMD ["./iconstation"]