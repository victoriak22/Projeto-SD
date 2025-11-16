# Dockerfile para o Broker (Python)
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

# Expor portas do broker
EXPOSE 5557 5558

# Executar broker
CMD ["python", "main.py"]