FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o proxy-routes ./cmd/proxy/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/proxy-routes /app

EXPOSE 8080

CMD ["sh", "-c", "./proxy-routes"]