FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o product-service ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/product-service .
EXPOSE 8080

CMD ["./product-service"]
