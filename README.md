# Sistema de Mensagens Instantâneas - Partes 1 e 2

Sistema distribuído de mensagens instantâneas inspirado em BBS/IRC, desenvolvido para a disciplina de Sistemas Distribuídos.

## 📋 Parte 1: Request-Reply ✅

Implementação do padrão Request-Reply para comunicação entre cliente e servidor.

### Funcionalidades Implementadas

- ✅ **Login de usuários**: Cadastro de novos usuários no sistema
- ✅ **Listagem de usuários**: Visualização de todos os usuários cadastrados
- ✅ **Criação de canais**: Criação de novos canais de comunicação
- ✅ **Listagem de canais**: Visualização de todos os canais disponíveis
- ✅ **Persistência de dados**: Armazenamento em disco de logins e canais

## 📋 Parte 2: Publisher-Subscriber ✅

Implementação do padrão Pub/Sub para troca de mensagens entre usuários.

### Funcionalidades Implementadas

- ✅ **Broker Pub/Sub**: Proxy XSUB/XPUB para distribuição de mensagens
- ✅ **Publicação em canais**: Usuários podem publicar mensagens em canais públicos
- ✅ **Mensagens diretas**: Envio de mensagens privadas entre usuários
- ✅ **Inscrição em canais**: Usuários podem se inscrever em canais para receber mensagens
- ✅ **Cliente automatizado**: Bots que geram mensagens aleatórias para testes
- ✅ **Persistência de mensagens**: Todas as mensagens são armazenadas em disco

### 🏗️ Arquitetura

**Parte 1 - Request-Reply:**
```
Cliente (Node.js) <---> Servidor (Go)
    REQ                    REP
```

**Parte 2 - Publisher-Subscriber:**
```
Clientes REQ ──┐
               ├──> Server (REP + PUB) ──> Broker (XSUB/XPUB) ──┐
Clientes REQ ──┘                                                  ├──> Clientes SUB
                                                                  │    (recebem msgs)
Clientes Auto ────────────────────────────────────────────────────┘
```

- **Cliente**: Interface CLI interativa em JavaScript/Node.js (REQ + SUB)
- **Servidor**: Backend em Go com ZeroMQ (REP + PUB) e persistência JSON
- **Broker**: Proxy Pub/Sub em Python (XSUB/XPUB)
- **Clientes Automatizados**: Bots em Python que geram mensagens

### 🛠️ Tecnologias

- **Server**: Go 1.21 + ZeroMQ (REP + PUB)
- **Client**: Node.js 20 + ZeroMQ (REQ + SUB)
- **Broker**: Python 3.11 + ZeroMQ (XSUB/XPUB)
- **Cliente Automatizado**: Python 3.11 + ZeroMQ (REQ)
- **Comunicação**: ZeroMQ (Request-Reply + Pub/Sub patterns)
- **Persistência**: JSON
- **Containerização**: Docker + Docker Compose

## 🚀 Como Executar

### Pré-requisitos

- Docker
- Docker Compose

### Executar o sistema completo

```bash
# Construir e iniciar todos os containers
docker-compose up --build

# Executar em background
docker-compose up -d --build
```

### Interagir com o cliente

```bash
# Acessar o container do cliente
docker exec -it messaging-client sh

# Dentro do container, iniciar o cliente
npm start
```

Ou diretamente:

```bash
docker-compose exec client npm start
```

### Testar múltiplos clientes

Para simular múltiplos usuários, você pode iniciar vários clientes:

```bash
# Terminal 1
docker-compose run --rm client npm start

# Terminal 2
docker-compose run --rm client npm start

# Terminal 3
docker-compose run --rm client npm start
```

### Parar o sistema

```bash
docker-compose down

# Para limpar volumes (apaga dados persistentes)
docker-compose down -v
```

## 📁 Estrutura do Projeto

```
.
├── broker/
│   └── main.py              # (Parte 2)
├── client/
│   ├── Dockerfile
│   ├── main.js              # Cliente Node.js
│   └── package.json
├── proxy/
│   └── main.py              # (Parte 2)
├── server/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go              # Servidor Go
├── docker-compose.yml
├── .gitignore
└── README.md
```

## 🔌 Formato das Mensagens

### Login
**Request:**
```json
{
  "service": "login",
  "data": {
    "user": "nome_usuario",
    "timestamp": 1234567890
  }
}
```

### Publicar em Canal (Parte 2)
**Request:**
```json
{
  "service": "publish",
  "data": {
    "user": "alice",
    "channel": "geral",
    "message": "Olá pessoal!",
    "timestamp": 1234567890
  }
}
```

**Response:**
```json
{
  "service": "publish",
  "data": {
    "status": "OK",
    "timestamp": 1234567890
  }
}
```

### Mensagem Direta (Parte 2)
**Request:**
```json
{
  "service": "message",
  "data": {
    "src": "alice",
    "dst": "bob",
    "message": "Oi Bob, tudo bem?",
    "timestamp": 1234567890
  }
}
```

**Response:**
```json
{
  "service": "message",
  "data": {
    "status": "OK",
    "timestamp": 1234567890
  }
}
```

### Mensagem Publicada no Broker
**Tópico**: Nome do canal ou usuário  
**Payload (canal)**:
```json
{
  "user": "alice",
  "message": "Olá pessoal!",
  "timestamp": 1234567890
}
```

**Payload (mensagem direta)**:
```json
{
  "from": "alice",
  "message": "Oi Bob!",
  "timestamp": 1234567890
}
```

**Response:**
```json
{
  "service": "login",
  "data": {
    "status": "sucesso",
    "timestamp": 1234567890
  }
}
```

### Listar Usuários
**Request:**
```json
{
  "service": "users",
  "data": {
    "timestamp": 1234567890
  }
}
```

**Response:**
```json
{
  "service": "users",
  "data": {
    "timestamp": 1234567890,
    "users": ["alice", "bob", "charlie"]
  }
}
```

### Criar Canal
**Request:**
```json
{
  "service": "channel",
  "data": {
    "channel": "geral",
    "timestamp": 1234567890
  }
}
```

**Response:**
```json
{
  "service": "channel",
  "data": {
    "status": "sucesso",
    "timestamp": 1234567890
  }
}
```

### Listar Canais
**Request:**
```json
{
  "service": "channels",
  "data": {
    "timestamp": 1234567890
  }
}
```

**Response:**
```json
{
  "service": "channels",
  "data": {
    "timestamp": 1234567890,
    "channels": ["geral", "random", "tech"]
  }
}
```

## 🧪 Testando

### Cenário Parte 1 - Request-Reply

1. **Iniciar o sistema**
   ```bash
   docker-compose up --build
   ```

2. **Abrir 3 terminais para 3 clientes diferentes**

3. **Terminal 1 - Usuário Alice**
   ```
   1. Fazer login como "alice"
   2. Criar canal "geral"
   3. Listar canais (deve ver "geral")
   ```

4. **Terminal 2 - Usuário Bob**
   ```
   1. Fazer login como "bob"
   2. Listar usuários (deve ver "alice" e "bob")
   3. Criar canal "tech"
   4. Listar canais (deve ver "geral" e "tech")
   ```

5. **Terminal 3 - Usuário Charlie**
   ```
   1. Fazer login como "charlie"
   2. Listar usuários (deve ver "alice", "bob" e "charlie")
   3. Listar canais (deve ver "geral" e "tech")
   ```

### Cenário Parte 2 - Pub/Sub

1. **Iniciar o sistema com clientes automatizados**
   ```bash
   docker-compose up --build
   ```

2. **Terminal 1 - Alice**
   ```
   1. Login como "alice"
   2. Inscrever no canal "geral" (opção 5)
   3. Publicar mensagem no "geral" (opção 6)
   4. Aguardar e ver mensagens dos bots
   ```

3. **Terminal 2 - Bob**
   ```
   1. Login como "bob"
   2. Inscrever no canal "geral" (opção 5)
   3. Ver mensagens de Alice e dos bots
   4. Enviar mensagem direta para Alice (opção 7)
   ```

4. **Terminal 3 - Charlie**
   ```
   1. Login como "charlie"
   2. Listar canais (opção 4)
   3. Inscrever em múltiplos canais
   4. Ver mensagens de todos os canais inscritos
   ```

5. **Verificar clientes automatizados**
   ```bash
   # Ver logs dos bots
   make logs-auto
   
   # Os bots devem estar enviando mensagens automaticamente
   ```

6. **Verificar persistência**
   ```bash
   # Parar containers
   docker-compose down
   
   # Reiniciar
   docker-compose up
   
   # Os dados devem persistir!
   ```

## 📊 Logs e Debug

Ver logs do servidor:
```bash
docker-compose logs -f server
```

Ver logs do cliente:
```bash
docker-compose logs -f client
```

Ver dados persistidos:
```bash
docker exec -it messaging-server cat /data/server_data.json
```

## 🐛 Troubleshooting

### Cliente não conecta ao servidor
- Verifique se o servidor está rodando: `docker-compose ps`
- Veja os logs: `docker-compose logs server`

### Erro ao buildar
- Limpe containers antigos: `docker-compose down -v`
- Reconstrua: `docker-compose build --no-cache`

### Dados não persistem
- Verifique se o volume está criado: `docker volume ls`
- Veja o conteúdo: `docker exec -it messaging-server ls -la /data`

## 📝 Próximas Partes

- **Parte 2**: Publisher-Subscriber (Broker e troca de mensagens)
- **Parte 3**: MessagePack (Serialização eficiente)
- **Parte 4**: Relógios (Lamport, vetoriais)
- **Parte 5**: Consistência e Replicação

## 👥 Desenvolvimento

Este projeto foi desenvolvido como parte da disciplina de Sistemas Distribuídos, utilizando 3+ linguagens de programação:
- Go (Server)
- JavaScript/Node.js (Client)
- Python (Broker e Proxy - Parte 2)

## 📄 Licença

MIT