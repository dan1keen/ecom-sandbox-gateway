# ---------- STAGE 1 ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/main.go

# ---------- STAGE 2 ----------
FROM alpine:latest

WORKDIR /app
RUN adduser -D appuser

COPY --from=builder /app/app .

USER appuser
EXPOSE 8080

CMD ["./app"]