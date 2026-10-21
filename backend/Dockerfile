# Stage 1: Build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o roccia-api ./cmd/api/main.go

# Stage 2: Produzione
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/roccia-api .
EXPOSE 8080
CMD ["./roccia-api"]