# BRLIA Waterfall Engine Dockerfile
# Multi-stage build para imagem de producao otimizada

# Build stage
FROM golang:1.21-alpine AS builder

# Configurar ambiente
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct

# Instalar dependencias
WORKDIR /app

# Copiar arquivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar codigo fonte
COPY . .

# Buildar aplicacao
RUN go build -a -installsuffix cgo -o /brlia-waterfall main.go

# Production stage
FROM alpine:3.18

# Instalar dependencias de runtime
RUN apk add --no-cache ca-certificates tzdata

# Configurar timezone
ENV TZ=America/Sao_Paulo

# Criar usuario
RUN adduser -D -u 1000 brlia

# Copiar binario
COPY --from=builder /brlia-waterfall /usr/local/bin/brlia-waterfall

# Criar diretorio de configuracao
RUN mkdir -p /etc/brlia

# Copiar configuracao padrao
COPY config/default.yaml /etc/brlia/waterfall.yaml

# Expor porta
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Entrypoint
USER brlia
CMD ["brlia-waterfall"]
