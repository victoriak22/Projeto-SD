# Dockerfile para o Servidor de Referência (Python)
FROM python:3.11-alpine

# Instalar dependências do sistema para ZeroMQ
RUN apk add --no-cache \
    zeromq-dev \
    gcc \
    musl-dev

WORKDIR /app

# Copiar requirements e instalar dependências
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copiar código fonte
COPY main.py .

# Criar diretório para dados persistentes
RUN mkdir -p /data

# Expor porta
EXPOSE 5559

# Executar servidor de referência
CMD ["python", "main.py"]