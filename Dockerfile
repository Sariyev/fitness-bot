FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy

# Build bot binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

# Build admin binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin ./cmd/admin

# Build webapp binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /webapp ./cmd/webapp

# Frontend build stage
FROM node:20-alpine AS frontend
WORKDIR /web
COPY web/package.json web/package-lock.json* ./
RUN npm install
COPY web/ .
RUN npm run build

# Bot image
FROM alpine:latest AS bot
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bot .
COPY --from=builder /app/migrations ./migrations
CMD ["./bot"]

# Admin image
FROM alpine:latest AS admin
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /admin .
COPY --from=builder /app/migrations ./migrations
CMD ["./admin"]

# WebApp image
FROM alpine:latest AS webapp
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /webapp .
COPY --from=builder /app/migrations ./migrations
COPY --from=frontend /web/dist ./static
CMD ["./webapp"]
