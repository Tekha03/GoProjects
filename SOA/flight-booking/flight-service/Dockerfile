FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o flight-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/flight-service /usr/local/bin/flight-service
COPY --from=builder /app/proto /proto
CMD ["flight-service"]
