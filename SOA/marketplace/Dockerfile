# syntax=docker/dockerfile:1

# 1. Build stage
FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git bash

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Генерация OpenAPI (если используется oapi-codegen)
# RUN go generate ./...

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o marketplace ./cmd/app

# 2. Run stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates bash

WORKDIR /app

# Копируем бинарник из билд-стейджа
COPY --from=builder /app/marketplace .

# Копируем SQL миграции
COPY ./migrations ./migrations

# ENV example, можно переопределять через docker-compose
ENV APP_ENV=production
ENV PORT=8080

EXPOSE 8080

CMD ["./marketplace"]
