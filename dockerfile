# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app .
COPY .env .

EXPOSE 8080
CMD ["./rate-limiter"]
