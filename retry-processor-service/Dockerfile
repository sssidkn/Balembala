FROM golang:1.23.3-bullseye AS builder

RUN apt-get update && \
    apt-get install -y \
    gcc \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /retry-processor-service ./cmd/main.go

FROM debian:bullseye-slim
RUN apt-get update && \
    apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /retry-processor-service /retry-processor-service
ENTRYPOINT ["/retry-processor-service"]
