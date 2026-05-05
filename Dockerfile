FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY main.go ./
COPY static ./static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o iconstation main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/iconstation ./

RUN mkdir -p /app/UserData/icons /app/UserData/chunks \
    && apk add --no-cache ca-certificates tzdata

EXPOSE 9168
CMD ["./iconstation"]
