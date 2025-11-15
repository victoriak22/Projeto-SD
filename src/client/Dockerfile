# Dockerfile para o Client (Node.js)
FROM node:20-alpine

# Instalar dependências do sistema necessárias para ZeroMQ
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    zeromq-dev

WORKDIR /app

# Copiar arquivos de dependências
COPY package*.json ./

# Instalar dependências
RUN npm install

# Copiar código fonte
COPY main.js ./

# Variável de ambiente padrão para o servidor
ENV SERVER_URL=tcp://server:5555

# Executar cliente
CMD ["npm", "start"]