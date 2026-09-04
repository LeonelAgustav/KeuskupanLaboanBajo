# Stage 1: Build
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Jika butuh file env sebagai fallback
COPY --from=builder /app/.env.example .env

EXPOSE 8080
CMD ["./main"]