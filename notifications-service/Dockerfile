FROM golang:1.23.3-bullseye AS builder

RUN sed -i 's|http://deb.debian.org|http://mirror.docker.ru|g' /etc/apt/sources.list \
    && sed -i 's|http://security.debian.org|http://mirror.docker.ru/debian-security|g' /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y \
    gcc \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /notifications-service ./cmd/main.go

FROM debian:bullseye-slim
RUN sed -i 's|http://deb.debian.org|http://mirror.docker.ru|g' /etc/apt/sources.list \
    && sed -i 's|http://security.debian.org|http://mirror.docker.ru/debian-security|g' /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /notifications-service /notifications-service
EXPOSE 9010 8083
ENTRYPOINT ["/notifications-service"]
