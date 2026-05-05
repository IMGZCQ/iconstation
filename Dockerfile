FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o iconstation main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/iconstation ./
COPY static /app/static
COPY icon.png /usr/share/icons/icon.png

RUN mkdir -p /app/UserData/icons /app/UserData/chunks \
    && apk add --no-cache ca-certificates tzdata

LABEL org.opencontainers.image.icon=/usr/share/icons/icon.png

EXPOSE 9168
CMD ["./iconstation"]