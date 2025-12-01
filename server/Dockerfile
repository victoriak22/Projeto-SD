# Dockerfile para o Server (Go)
FROM golang:1.21-alpine AS builder

# Instalar dependências do sistema necessárias para ZeroMQ
RUN apk add --no-cache \
    gcc \
    musl-dev \
    zeromq-dev

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum* ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar aplicação
RUN go build -o server main.go

# Imagem final
FROM alpine:latest

# Instalar runtime do ZeroMQ
RUN apk add --no-cache zeromq

WORKDIR /app

# Copiar binário compilado
COPY --from=builder /app/server .

# Criar diretório para dados persistentes
RUN mkdir -p /data

# Expor porta
EXPOSE 5555

# Executar servidor
CMD ["./server"]
