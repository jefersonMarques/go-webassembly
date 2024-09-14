# Etapa de construção do backend gRPC
FROM golang:1.20 AS builder

WORKDIR /app

# Copiar os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Não baixar dependências ou compilar automaticamente
# Vamos fazer isso manualmente dentro do container

# Copiar o restante do código do aplicativo
COPY . .

# Finalização
FROM ubuntu:20.04

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 50051

# Manter o container rodando indefinidamente para você acessar e resolver manualmente
CMD ["sleep", "infinity"]
