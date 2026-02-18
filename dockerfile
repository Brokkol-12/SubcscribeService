# -------- BUILD --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# -------- RUNTIME --------
FROM alpine:latest

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/app .

USER appuser

EXPOSE 8081

CMD ["./app"]
