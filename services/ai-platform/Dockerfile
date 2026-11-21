# BRLIA BRLIA AI Platform Dockerfile
# Multi-stage build para imagem de producao otimizada

# Build stage
FROM golang:1.21-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o /brlia-ai-platform main.go

# Production stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=America/Sao_Paulo

RUN adduser -D -u 1000 brlia

COPY --from=builder /brlia-ai-platform /usr/local/bin/brlia-ai-platform

RUN mkdir -p /etc/brlia

COPY config/default.yaml /etc/brlia/ai-platform.yaml

EXPOSE 8090

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8090/health || exit 1

USER brlia
CMD ["brlia-ai-platform"]
