FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/main .
COPY config ./config

CMD ["./main"]