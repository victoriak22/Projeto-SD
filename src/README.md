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

## 📋 Parte 3: MessagePack ✅

Otimização da serialização de mensagens usando MessagePack ao invés de JSON.

### Funcionalidades Implementadas

- ✅ **Serialização eficiente**: Mensagens em formato binário (MessagePack)
- ✅ **Compatibilidade entre linguagens**: Go, JavaScript e Python usando MessagePack
- ✅ **Redução de tamanho**: Mensagens ~25% menores que JSON
- ✅ **Melhor performance**: Serialização/deserialização mais rápida
- ✅ **Transparente**: Mesma funcionalidade, formato diferente

## 📋 Parte 4: Relógios ⏳ (Em andamento)

### Etapa 1: Relógio Lógico de Lamport ✅

Implementação de relógios lógicos para ordenação de eventos distribuídos.

#### Funcionalidades Implementadas

- ✅ **Relógio lógico em todos os processos**: Server, Client e Auto-client
- ✅ **Incremento antes de enviar**: `clock++` antes de cada envio
- ✅ **Atualização ao receber**: `clock = max(local, recebido) + 1`
- ✅ **Campo clock em todas as mensagens**: Incluído em requests e responses
- ✅ **Logs com clock**: Todas as operações mostram o valor do relógio lógico

#### Como Funciona

**Algoritmo de Lamport:**
1. Cada processo mantém um contador (`logicalClock`)
2. Antes de enviar mensagem: incrementa o contador
3. Ao receber mensagem: `clock = max(clock_local, clock_recebido) + 1`

**Exemplo de uso:**
```
Cliente envia login (clock=1) → 
Servidor recebe (atualiza para clock=2) →
Servidor responde (clock=3) →
Cliente recebe (atualiza para clock=4)
```

### Etapa 2: Servidor de Referência ✅

Novo componente para gerenciar registro e descoberta de servidores.

#### Funcionalidades Implementadas

- ✅ **Novo container `reference`**: Servidor de referência em Python
- ✅ **Serviço `rank`**: Atribui rank único a cada servidor
- ✅ **Serviço `list`**: Retorna lista de servidores ativos
- ✅ **Serviço `heartbeat`**: Mantém lista de servidores atualizada
- ✅ **Persistência**: Salva lista de servidores em disco
- ✅ **Cleanup automático**: Remove servidores inativos (timeout)
- ✅ **Servidor se registra ao iniciar**: Obtém rank automaticamente
- ✅ **Heartbeat periódico**: A cada 10 segundos

#### Como Funciona

**Registro de Servidores:**
1. Servidor inicia e conecta ao reference (porta 5559)
2. Envia requisição `rank` com seu nome
3. Reference atribui rank único ou retorna existente
4. Servidor armazena seu rank

**Heartbeat:**
1. Servidor envia heartbeat a cada 10 segundos
2. Reference atualiza timestamp do servidor
3. Servidores sem heartbeat por 60s são removidos

**Exemplo:**
```
server-1 → registra → rank 1
server-2 → registra → rank 2
server-3 → registra → rank 3

Heartbeats mantêm servidores na lista ativa
```

### Etapa 3: Múltiplos Servidores ✅

Configuração de 3 réplicas do servidor para alta disponibilidade.

#### Funcionalidades Implementadas

- ✅ **3 servidores independentes**: server-1, server-2, server-3
- ✅ **Ranks únicos**: Cada servidor tem rank diferente (1, 2, 3)
- ✅ **Dados independentes**: Cada servidor tem seu próprio volume
- ✅ **Portas diferentes**: 5555, 5556, 5557 (externamente)
- ✅ **Todos registrados**: Todos se conectam ao reference
- ✅ **Heartbeats simultâneos**: Todos enviam heartbeat periódico
- ✅ **Clientes distribuídos**: Cada auto-client conecta a servidor diferente

#### Como Funciona

**Configuração:**
```
server-1: porta 5555, rank 1, volume server-1-data
server-2: porta 5556, rank 2, volume server-2-data
server-3: porta 5557, rank 3, volume server-3-data
```

**Distribuição de Clientes:**
- client → server-1
- auto-client-1 → server-1
- auto-client-2 → server-2

**Exemplo de Logs:**
```
Reference:
✅ Novo servidor registrado: server-1 com rank 1
✅ Novo servidor registrado: server-2 com rank 2
✅ Novo servidor registrado: server-3 com rank 3
💓 Heartbeat recebido de server-1
💓 Heartbeat recebido de server-2
💓 Heartbeat recebido de server-3
```

### Etapa 4: Sincronização Berkeley ✅

Implementação do Algoritmo de Berkeley para sincronização de relógios físicos.

#### Funcionalidades Implementadas

- ✅ **Coordenador eleito**: Servidor com maior rank (server-3)
- ✅ **Coleta de timestamps**: Coordenador pede tempo de todos
- ✅ **Cálculo de média**: Calcula tempo médio de todos os servidores
- ✅ **Distribuição de ajustes**: Envia ajuste individual para cada servidor
- ✅ **Aplicação de ajustes**: Servidores ajustam seus relógios
- ✅ **Sincronização periódica**: A cada 10 mensagens processadas
- ✅ **Offset de tempo**: Mantém ajuste sem modificar relógio do sistema

#### Como Funciona

**Algoritmo de Berkeley:**
1. Coordenador (maior rank) coleta timestamps de todos os servidores
2. Calcula tempo médio: `média = soma(timestamps) / N`
3. Para cada servidor, calcula ajuste: `ajuste = média - tempo_servidor`
4. Distribui ajustes individuais
5. Servidores aplicam: `tempo_ajustado = tempo_real + offset`

**Exemplo:**
```
server-1: tempo 100, ajuste +5 → tempo_ajustado 105
server-2: tempo 110, ajuste -5 → tempo_ajustado 105
server-3: tempo 105, ajuste  0 → tempo_ajustado 105

Todos sincronizados em 105!
```

**Logs esperados:**
```
Server-3 (coordenador):
🎯 Iniciando sincronização Berkeley como COORDENADOR
📊 Coletando timestamps de 3 servidores...
   📥 server-1: 1700000100
   📥 server-2: 1700000110
📊 Tempo médio calculado: 1700000105
   📤 Enviado ajuste de +5s para server-1
   📤 Enviado ajuste de -5s para server-2
✅ Sincronização Berkeley concluída

Server-1:
⏰ Relógio ajustado em +5s (offset total: +5s)

Server-2:
⏰ Relógio ajustado em -5s (offset total: -5s)
```

### Etapa 5: Eleição Bully ✅

Implementação do Algoritmo Bully para eleição de coordenador.

#### Funcionalidades Implementadas

- ✅ **Detecção de falha**: Verifica coordenador a cada 30 segundos
- ✅ **Algoritmo Bully**: Eleição baseada em rank
- ✅ **Mensagens de eleição**: Envia `election` para ranks maiores
- ✅ **Resposta OK**: Servidores com rank maior respondem e iniciam própria eleição
- ✅ **Anúncio de coordenador**: Publicado no tópico `servers`
- ✅ **Subscrição ao tópico**: Todos os servidores recebem anúncios
- ✅ **Atualização automática**: Todos atualizam coordenador atual
- ✅ **Coordenador inicial**: Determinado ao iniciar (maior rank)

#### Como Funciona

**Algoritmo Bully:**
1. Servidor detecta que coordenador não responde
2. Envia `election` para todos com rank maior
3. Se alguém responde "OK", aguarda novo coordenador
4. Se ninguém responde, se torna coordenador
5. Publica no tópico `servers`
6. Todos recebem e atualizam

**Exemplo com 3 servidores:**
```
Estado inicial:
- server-3 (rank 3) é coordenador

[server-3 falha ou para de responder]

server-2 detecta:
- Envia election para server-3
- Timeout (sem resposta)
- Ninguém com rank maior respondeu
- Se torna coordenador
- Publica no tópico 'servers'

server-1 recebe anúncio:
- Atualiza: coordenador = server-2
```

**Logs esperados:**
```
Server-2:
⚠️ Coordenador server-3 não respondeu - iniciando eleição
🗳️ Iniciando eleição Bully...
📤 Enviando eleição para 1 servidores com rank maior
⚠️ server-3 não respondeu (pode estar offline)
👑 Ninguém respondeu. Me tornando coordenador!
👑 Agora sou o COORDENADOR (rank 2)
📢 Anúncio de coordenador publicado no tópico 'servers'

Server-1:
📢 Novo coordenador anunciado: server-2
```

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

- **Server**: Go 1.21 + ZeroMQ (REP + PUB) + MessagePack
- **Client**: Node.js 20 + ZeroMQ (REQ + SUB) + MessagePack
- **Broker**: Python 3.11 + ZeroMQ (XSUB/XPUB)
- **Cliente Automatizado**: Python 3.11 + ZeroMQ (REQ) + MessagePack
- **Comunicação**: ZeroMQ (Request-Reply + Pub/Sub patterns)
- **Serialização**: MessagePack (binário, eficiente)
- **Persistência**: JSON (legível para humanos)
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
# Acessar o container do cliente interativo
docker exec -it messaging-client npm start

# Ou criar novo cliente
docker-compose run --rm client npm start
```

### Testar múltiplos clientes

Para simular múltiplos usuários, abra vários terminais:

```bash
# Terminal 1 - Alice
docker-compose run --rm client npm start

# Terminal 2 - Bob
docker-compose run --rm client npm start

# Terminal 3 - Charlie
docker-compose run --rm client npm start
```

### Ver logs

```bash
# Todos os serviços
docker-compose logs -f

# Servidor apenas
docker-compose logs -f server

# Broker apenas
docker-compose logs -f broker

# Clientes automatizados
docker-compose logs -f auto-client-1 auto-client-2
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
│   ├── main.py              # Broker Pub/Sub (Python)
│   ├── requirements.txt
│   └── Dockerfile
├── client/
│   ├── main.js              # Cliente interativo (Node.js)
│   ├── auto_client.py       # Cliente automatizado (Python)
│   ├── package.json
│   ├── Dockerfile
│   └── Dockerfile.auto
├── proxy/
│   └── main.py              # Placeholder (próximas partes)
├── server/
│   ├── main.go              # Servidor (Go)
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml
├── .gitignore
├── README.md
└── ARCHITECTURE.md
```

## 🔌 Formato das Mensagens

**Nota**: A partir da Parte 3, todas as mensagens são serializadas em **MessagePack** (binário) ao invés de JSON. Os exemplos abaixo mostram o formato lógico das mensagens.

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

Ver logs do broker:
```bash
docker-compose logs -f broker
```

Ver logs dos clientes automatizados:
```bash
docker-compose logs -f auto-client-1 auto-client-2
```

Ver dados persistidos:
```bash
docker exec -it messaging-server cat /data/server_data.json
```

Ver status dos containers:
```bash
docker-compose ps
```

## 🐛 Troubleshooting

### Cliente não conecta ao servidor
- Verifique se o servidor está rodando: `docker-compose ps`
- Veja os logs: `docker-compose logs server`
- Reinicie: `docker-compose restart server`

### Mensagens não chegam (Pub/Sub)
- Verifique se o broker está rodando: `docker-compose ps`
- Certifique-se de que o cliente se inscreveu no canal (opção 5)
- Veja os logs do broker: `docker-compose logs broker`

### Erro ao buildar
- Limpe containers antigos: `docker-compose down -v`
- Reconstrua: `docker-compose build --no-cache`

### Dados não persistem
- Verifique o volume: `docker volume ls | grep messaging`
- Veja o conteúdo: `docker exec -it messaging-server ls -la /data`

## 📝 Próximas Partes

- **Parte 4**: Relógios (Lamport, vetoriais)
- **Parte 5**: Consistência e Replicação

## 👥 Desenvolvimento

Este projeto foi desenvolvido como parte da disciplina de Sistemas Distribuídos, utilizando 3 linguagens de programação:
- **Go** (Server) com MessagePack
- **JavaScript/Node.js** (Client interativo) com MessagePack
- **Python** (Broker e Cliente automatizado) com MessagePack

## 🎯 Bibliotecas MessagePack Utilizadas

- **Go**: `github.com/vmihailenco/msgpack/v5` - Serialização eficiente para Go
- **JavaScript**: `@msgpack/msgpack` - Implementação oficial para Node.js
- **Python**: `msgpack` - Biblioteca padrão para Python

### Comparação de Tamanho das Mensagens

**Exemplo: Login Request**

JSON (60 bytes):
```json
{"service":"login","data":{"user":"alice","timestamp":1234567890}}
```

MessagePack (~45 bytes - 25% menor):
```
\x82\xa7service\xa5login\xa4data\x82\xa4user\xa5alice\xa9timestamp\xce\x49\x96\x02\xd2
```

**Vantagens do MessagePack**:
- 📉 Mensagens menores (15-30% de redução)
- ⚡ Serialização/deserialização mais rápida
- 🔄 Compatível entre diferentes linguagens
- 💾 Menos uso de banda e memória

## 📄 Licença

MIT