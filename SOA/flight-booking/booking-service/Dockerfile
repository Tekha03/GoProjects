FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o booking-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/booking-service /usr/local/bin/booking-service
COPY --from=builder /app/proto /proto
CMD ["booking-service"]
