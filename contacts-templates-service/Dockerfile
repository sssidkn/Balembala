FROM golang:1.23.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /contacts-templates-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /contacts-templates-service .

EXPOSE 8082 9080
CMD ["./contacts-templates-service"]
