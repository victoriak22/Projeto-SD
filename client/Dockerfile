# Dockerfile para o Cliente Automatizado (Python)
FROM python:3.11-alpine

# Instalar dependências do sistema para ZeroMQ
RUN apk add --no-cache \
    zeromq-dev \
    gcc \
    musl-dev

WORKDIR /app

# Copiar requirements reais, agora usando os que você criou
COPY requirements.txt .

# Instalar dependências Python
RUN pip install --no-cache-dir -r requirements.txt

# Copiar script do cliente automatizado
COPY auto_client.py .

# Executar cliente automatizado
CMD ["python", "auto_client.py"]
