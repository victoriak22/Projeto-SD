# Sistema de Mensagens Instantâneas - Projeto Completo

Sistema distribuído de mensagens instantâneas inspirado em BBS/IRC, desenvolvido para a disciplina de Sistemas Distribuídos.

## 📋 Visão Geral

Este projeto implementa um sistema completo de mensagens distribuídas com:
- ✅ 5 Partes implementadas (Request-Reply, Pub/Sub, MessagePack, Relógios, Replicação)
- ✅ 3 Linguagens de programação (Go, JavaScript/Node.js, Python)
- ✅ Múltiplos padrões de comunicação (REQ-REP, PUB-SUB)
- ✅ Sincronização de relógios (Lamport + Berkeley)
- ✅ Eleição de coordenador (Bully)
- ✅ Replicação de dados (Primary-Backup)
- ✅ Alta disponibilidade e tolerância a falhas

---

## 📋 Parte 1: Request-Reply ✅

Implementação do padrão Request-Reply para comunicação entre cliente e servidor.

### Funcionalidades Implementadas

- ✅ **Login de usuários**: Cadastro de novos usuários no sistema
- ✅ **Listagem de usuários**: Visualização de todos os usuários cadastrados
- ✅ **Criação de canais**: Criação de novos canais de comunicação
- ✅ **Listagem de canais**: Visualização de todos os canais disponíveis
- ✅ **Persistência de dados**: Armazenamento em disco de logins e canais

---

## 📋 Parte 2: Publisher-Subscriber ✅

Implementação do padrão Pub/Sub para troca de mensagens entre usuários.

### Funcionalidades Implementadas

- ✅ **Broker Pub/Sub**: Proxy XSUB/XPUB para distribuição de mensagens
- ✅ **Publicação em canais**: Usuários podem publicar mensagens em canais públicos
- ✅ **Mensagens diretas**: Envio de mensagens privadas entre usuários
- ✅ **Inscrição em canais**: Usuários podem se inscrever em canais para receber mensagens
- ✅ **Cliente automatizado**: Bots que geram mensagens aleatórias para testes
- ✅ **Persistência de mensagens**: Todas as mensagens são armazenadas em disco

---

## 📋 Parte 3: MessagePack ✅

Otimização da serialização de mensagens usando MessagePack ao invés de JSON.

### Funcionalidades Implementadas

- ✅ **Serialização eficiente**: Mensagens em formato binário (MessagePack)
- ✅ **Compatibilidade entre linguagens**: Go, JavaScript e Python usando MessagePack
- ✅ **Redução de tamanho**: Mensagens ~25% menores que JSON
- ✅ **Melhor performance**: Serialização/deserialização mais rápida
- ✅ **Transparente**: Mesma funcionalidade, formato diferente

### Comparação de Tamanho

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

---

## 📋 Parte 4: Relógios ✅

Implementação de relógios lógicos e físicos para sincronização em sistemas distribuídos.

### Etapa 1: Relógio Lógico de Lamport ✅

Implementação de relógios lógicos para ordenação de eventos distribuídos.

#### Funcionalidades

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

**Exemplo:**
```
Cliente envia login (clock=1) → 
Servidor recebe (atualiza para clock=2) →
Servidor responde (clock=3) →
Cliente recebe (atualiza para clock=4)
```

### Etapa 2: Servidor de Referência ✅

Novo componente para gerenciar registro e descoberta de servidores.

#### Funcionalidades

- ✅ **Novo container `reference`**: Servidor de referência em Python (porta 5559)
- ✅ **Serviço `rank`**: Atribui rank único a cada servidor
- ✅ **Serviço `list`**: Retorna lista de servidores ativos
- ✅ **Serviço `heartbeat`**: Mantém lista de servidores atualizada
- ✅ **Persistência**: Salva lista de servidores em disco
- ✅ **Cleanup automático**: Remove servidores inativos (timeout 60s)
- ✅ **Heartbeat periódico**: A cada 10 segundos

#### Como Funciona

**Registro:**
1. Servidor inicia e conecta ao reference (porta 5559)
2. Envia requisição `rank` com seu nome
3. Reference atribui rank único (1, 2, 3, ...)
4. Servidor armazena seu rank

**Heartbeat:**
- Servidor envia heartbeat a cada 10 segundos
- Reference atualiza timestamp
- Servidores sem heartbeat por 60s são removidos

### Etapa 3: Múltiplos Servidores ✅

Configuração de 3 réplicas do servidor para alta disponibilidade.

#### Funcionalidades

- ✅ **3 servidores independentes**: server-1, server-2, server-3
- ✅ **Ranks únicos**: 1, 2, 3
- ✅ **Dados independentes**: Cada servidor tem seu próprio volume
- ✅ **Portas diferentes**: 5555, 5556, 5557 (externamente)
- ✅ **Todos registrados**: Conectados ao reference
- ✅ **Heartbeats simultâneos**: Todos enviam heartbeat periódico

**Configuração:**
```
server-1: porta 5555, rank 1, volume server-1-data
server-2: porta 5556, rank 2, volume server-2-data
server-3: porta 5557, rank 3, volume server-3-data
```

### Etapa 4: Sincronização Berkeley ✅

Implementação do Algoritmo de Berkeley para sincronização de relógios físicos.

#### Funcionalidades

- ✅ **Coordenador eleito**: Servidor com maior rank
- ✅ **Coleta de timestamps**: Coordenador pede tempo de todos
- ✅ **Cálculo de média**: `média = soma(timestamps) / N`
- ✅ **Distribuição de ajustes**: Envia ajuste individual para cada servidor
- ✅ **Aplicação de ajustes**: Servidores ajustam seus relógios
- ✅ **Sincronização periódica**: A cada 10 mensagens processadas
- ✅ **Offset de tempo**: Mantém ajuste sem modificar relógio do sistema

#### Como Funciona

**Algoritmo:**
1. Coordenador coleta timestamps: `T1=100, T2=110, T3=105`
2. Calcula média: `média = (100+110+105)/3 = 105`
3. Calcula ajustes: `A1=+5, A2=-5, A3=0`
4. Distribui ajustes para cada servidor
5. Todos sincronizados: `T1'=105, T2'=105, T3'=105`

**Logs esperados:**
```
🎯 Iniciando sincronização Berkeley como COORDENADOR
📊 Coletando timestamps de 3 servidores...
   📥 server-1: 100
   📥 server-2: 110
📊 Tempo médio calculado: 105
   📤 Enviado ajuste de +5s para server-1
   📤 Enviado ajuste de -5s para server-2
✅ Sincronização Berkeley concluída
```

### Etapa 5: Eleição Bully ✅

Implementação do Algoritmo Bully para eleição automática de coordenador.

#### Funcionalidades

- ✅ **Detecção de falha**: Verifica coordenador a cada 30 segundos
- ✅ **Algoritmo Bully**: Eleição baseada em rank (maior vence)
- ✅ **Mensagens de eleição**: Envia `election` para ranks maiores
- ✅ **Resposta OK**: Servidores maiores respondem e iniciam própria eleição
- ✅ **Anúncio de coordenador**: Publicado no tópico `servers`
- ✅ **Subscrição ao tópico**: Todos os servidores recebem anúncios
- ✅ **Atualização automática**: Todos atualizam coordenador atual
- ✅ **Coordenador inicial**: Determinado ao iniciar (maior rank)

#### Como Funciona

**Algoritmo Bully:**
1. Servidor detecta falha do coordenador
2. Envia `election` para todos com rank maior
3. Se alguém responde "OK": aguarda novo coordenador
4. Se ninguém responde: torna-se coordenador
5. Publica no tópico `servers`
6. Todos recebem e atualizam

**Exemplo:**
```
Estado inicial: server-3 (rank 3) é coordenador

[server-3 falha]

server-2 detecta → envia election → timeout → 
se torna coordenador → publica no tópico 'servers'

server-1 recebe anúncio → atualiza coordenador = server-2
```

**Logs esperados:**
```
⚠️ Coordenador server-3 não respondeu - iniciando eleição
🗳️ Iniciando eleição Bully...
👑 Ninguém respondeu. Me tornando coordenador!
📢 Anúncio de coordenador publicado no tópico 'servers'
```

---

## 📋 Parte 5: Consistência e Replicação ✅

Implementação de replicação de dados para garantir que todos os servidores tenham cópia completa dos dados.

### Problema

O broker distribui clientes entre servidores (load balancing). Consequentemente:
- ❌ Cada servidor possui apenas parte das mensagens
- ❌ Se um servidor falha, dados são perdidos
- ❌ Clientes recebem histórico incompleto ao consultar um servidor específico

### Solução Implementada

**Método escolhido: Primary-Backup com Propagação Assíncrona**

Adaptação do modelo Primary-Backup com as seguintes características:

#### Características do Método

1. **Primary (Coordenador)**: 
   - Servidor com maior rank atua como primary
   - Determinado pelo algoritmo Bully
   - Responsável por coordenar sincronização

2. **Backups**: 
   - Todos os outros servidores são backups
   - Recebem replicações do primary e de outros servidores
   - Podem promover-se a primary via eleição

3. **Propagação Assíncrona**:
   - Replicação não bloqueia operações do usuário
   - Executada em goroutines/threads separadas
   - Melhor performance mas janela de inconsistência temporária

4. **Sincronização Periódica**:
   - A cada 60 segundos, backups sincronizam com coordenador
   - Garante convergência para consistência eventual
   - Resolve inconsistências e preenche lacunas

5. **Tolerância a Falhas**:
   - Eleição Bully garante novo primary automaticamente
   - Replicação continua após eleição
   - Dados não são perdidos

### Funcionalidades Implementadas

- ✅ **Replicação automática**: Dados replicados para todos os servidores ao salvar
- ✅ **Sincronização periódica**: A cada 60s, servidores solicitam sincronização completa
- ✅ **Sincronização sob demanda**: Serviço `sync` para sincronização manual
- ✅ **Replicação assíncrona**: Não bloqueia operações do usuário
- ✅ **Thread-safe**: Mutex protege acesso aos dados compartilhados
- ✅ **Consistência eventual**: Todos os servidores convergem para o mesmo estado
- ✅ **Merge inteligente**: Previne duplicatas usando timestamps

### Tipos de Dados Replicados

1. **Logins** (`login`): Novos usuários cadastrados
2. **Canais** (`channel`): Novos canais criados
3. **Mensagens de Canal** (`channel_message`): Publicações em canais
4. **Mensagens Diretas** (`user_message`): Mensagens entre usuários

### Fluxo de Replicação

**Operação Normal:**
```
1. Cliente faz login no server-1
2. server-1 salva localmente
3. server-1 replica assincronamente para server-2 e server-3
4. server-2 e server-3 recebem e salvam
5. Todos os servidores têm o login
```

**Sincronização Periódica:**
```
A cada 60 segundos:
1. Backups (server-1, server-2) solicitam sync do coordenador (server-3)
2. Coordenador envia todos os dados: logins, canais, mensagens
3. Backups fazem merge com dados locais
4. Duplicatas são ignoradas (usando timestamp + username/channel)
5. Sistema converge para consistência
```

### Formato das Mensagens

**Replicação:**
```json
{
  "service": "replicate",
  "data": {
    "type": "login",
    "content": {
      "username": "alice",
      "timestamp": 1234567890
    },
    "timestamp": 1234567890,
    "clock": 42
  }
}
```

**Sincronização Completa:**

Request:
```json
{
  "service": "sync",
  "data": {
    "last_sync": 1234567000,
    "timestamp": 1234567890,
    "clock": 50
  }
}
```

Response:
```json
{
  "service": "sync",
  "data": {
    "logins": [
      {"username": "alice", "timestamp": 1234567890},
      {"username": "bob", "timestamp": 1234567895}
    ],
    "channels": ["geral", "tech"],
    "channel_messages": [...],
    "user_messages": [...],
    "timestamp": 1234567890,
    "clock": 51
  }
}
```

### Modificações no Método Primary-Backup Tradicional

**Diferenças do Primary-Backup clássico:**

1. **Replicação Multi-Direcional**:
   - Clássico: Apenas primary replica para backups
   - **Nossa implementação**: Qualquer servidor pode replicar para outros
   - Vantagem: Mesmo sem ser primary, servidor pode garantir dados replicados

2. **Sincronização Periódica Adicional**:
   - Clássico: Apenas replicação sob demanda
   - **Nossa implementação**: Sync periódica a cada 60s
   - Vantagem: Autocorreção de inconsistências

3. **Eleição Automática de Primary**:
   - Clássico: Primary fixo ou manual
   - **Nossa implementação**: Algoritmo Bully elege automaticamente
   - Vantagem: Tolerância a falhas sem intervenção

4. **Assíncrono com Consistência Eventual**:
   - Clássico: Geralmente síncrono (bloqueante)
   - **Nossa implementação**: Assíncrono para performance
   - Trade-off: Janela de inconsistência aceitável

### Vantagens

✅ **Performance**: Replicação assíncrona não bloqueia cliente  
✅ **Simplicidade**: Coordenador centraliza lógica de sincronização  
✅ **Tolerância a Falhas**: Eleição automática + múltiplos backups  
✅ **Consistência Eventual**: Sistema converge automaticamente  
✅ **Escalabilidade**: Fácil adicionar novos servidores  
✅ **Autocorreção**: Sincronização periódica corrige inconsistências  

### Desvantagens e Trade-offs

⚠️ **Janela de Inconsistência**: Breve período (< 60s) onde dados podem não estar em todos  
⚠️ **Overhead de Rede**: Cada operação gera N-1 replicações  
⚠️ **Duplicatas Possíveis**: Sync pode criar duplicatas temporárias (aceitáveis)  
⚠️ **Não é ACID forte**: Consistência eventual, não imediata  

### Garantias Fornecidas

✅ **Disponibilidade**: Sistema continua funcionando com falhas  
✅ **Partição**: Tolera partições de rede temporárias  
✅ **Consistência Eventual**: Todos convergem para mesmo estado  
✅ **Durabilidade**: Dados persistidos em múltiplos servidores  

### Logs Esperados

```
Server-1 (recebe login):
✅ Novo usuário cadastrado: alice (clock: 15)
🔄 Replicando login para 2 servidores...
   ✅ Replicado para server-2
   ✅ Replicado para server-3

Server-2 (recebe replicação):
🔄 Login replicado: alice

Server-3 (recebe replicação):
🔄 Login replicado: alice

[60 segundos depois]

Server-1 (sincronização periódica):
🔄 Solicitando sincronização completa de server-3...
✅ Sincronização recebida: 5 logins, 3 canais, 10 msgs canal, 5 msgs diretas
🔄 Merge de dados locais com dados sincronizados
✅ Sincronização completa concluída
```

---

## 🏗️ Arquitetura Completa

```
                    ┌─────────────┐
                    │  Reference  │ :5559
                    │  (Python)   │
                    └─────────────┘
                         ↕ REQ/REP
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
   Server-1         Server-2         Server-3 (Primary/Coordenador)
   rank=1           rank=2           rank=3
   :5555            :5556            :5557
        ↓                ↓                ↓
    [Replicação entre servidores]
    [Sincronização Berkeley]
    [Eleição Bully]
        ↓                ↓                ↓
        └────────────────┼────────────────┘
                         ↓
                    ┌─────────┐
                    │  Broker │ :5557/:5558
                    │(XSUB/XPUB)│
                    └─────────┘
                         ↓
                  ┌──────┴──────┐
                  ↓             ↓
              Client        Auto-Clients
           (Node.js)        (Python)
```

### Componentes

- **Reference**: Registro de servidores, ranks, heartbeats
- **Servers (3x)**: Request-Reply + Publisher + Replicação
- **Broker**: Pub/Sub proxy (XSUB/XPUB)
- **Client**: Interface CLI interativa
- **Auto-clients**: Bots geradores de carga

---

## 🛠️ Tecnologias

- **Server**: Go 1.21 + ZeroMQ (REP + PUB) + MessagePack
- **Client**: Node.js 20 + ZeroMQ (REQ + SUB) + MessagePack
- **Broker**: Python 3.11 + ZeroMQ (XSUB/XPUB)
- **Reference**: Python 3.11 + ZeroMQ (REP) + MessagePack
- **Auto-client**: Python 3.11 + ZeroMQ (REQ) + MessagePack
- **Comunicação**: ZeroMQ (Request-Reply + Pub/Sub)
- **Serialização**: MessagePack (binário, eficiente)
- **Persistência**: JSON (legível)
- **Containerização**: Docker + Docker Compose

### Bibliotecas MessagePack

- **Go**: `github.com/vmihailenco/msgpack/v5`
- **JavaScript**: `@msgpack/msgpack`
- **Python**: `msgpack`

---

## 🚀 Como Executar

### Pré-requisitos

- Docker
- Docker Compose

### Iniciar o Sistema Completo

```bash
# Construir e iniciar todos os containers
docker-compose up --build

# Executar em background
docker-compose up -d --build
```

### Interagir com o Cliente

```bash
# Acessar cliente interativo
docker exec -it messaging-client npm start

# Ou criar novo cliente
docker-compose run --rm client npm start
```

### Testar Múltiplos Clientes

```bash
# Terminal 1 - Alice
docker-compose run --rm client npm start

# Terminal 2 - Bob
docker-compose run --rm client npm start

# Terminal 3 - Charlie
docker-compose run --rm client npm start
```

### Ver Logs

```bash
# Todos os serviços
docker-compose logs -f

# Servidores
docker-compose logs -f server-1 server-2 server-3

# Reference
docker-compose logs -f reference

# Broker
docker-compose logs -f broker

# Clientes automatizados
docker-compose logs -f auto-client-1 auto-client-2
```

### Parar o Sistema

```bash
# Parar containers
docker-compose down

# Limpar volumes (apaga dados)
docker-compose down -v
```

---

## 📁 Estrutura do Projeto

```
.
├── reference/
│   ├── main.py              # Servidor de referência (Python)
│   ├── requirements.txt
│   └── Dockerfile
├── broker/
│   ├── main.py              # Broker Pub/Sub (Python)
│   ├── requirements.txt
│   └── Dockerfile
├── server/
│   ├── main.go              # Servidor (Go)
│   ├── go.mod
│   └── Dockerfile
├── client/
│   ├── main.js              # Cliente interativo (Node.js)
│   ├── auto_client.py       # Cliente automatizado (Python)
│   ├── package.json
│   ├── Dockerfile
│   └── Dockerfile.auto
├── proxy/
│   └── main.py              # Placeholder
├── docker-compose.yml       # Orquestração (6 containers)
├── .gitignore
├── README.md
└── ARCHITECTURE.md
```

---

## 🧪 Testes Completos

### Teste 1: Request-Reply (Parte 1)

```bash
docker-compose up --build
docker-compose run --rm client npm start
```

1. Login como "alice"
2. Criar canal "geral"
3. Listar canais
4. Listar usuários

### Teste 2: Pub/Sub (Parte 2)

```bash
# Terminal 1 - Alice
docker-compose run --rm client npm start
# Login → Inscrever canal "geral" → Publicar mensagem

# Terminal 2 - Bob
docker-compose run --rm client npm start
# Login → Inscrever canal "geral" → Ver mensagens → Enviar DM para Alice
```

### Teste 3: Relógios (Parte 4)

```bash
# Ver logs com clocks
docker-compose logs server-1 | grep "clock:"

# Fazer 10 operações para forçar sincronização Berkeley
# Ver logs de sincronização
docker-compose logs server-3 | grep "Berkeley"
```

### Teste 4: Eleição (Parte 4)

```bash
# Parar coordenador
docker-compose stop server-3

# Aguardar 30s e ver eleição
docker-compose logs server-2 | grep "eleição"

# Deve mostrar: server-2 se torna coordenador
```

### Teste 5: Replicação (Parte 5)

```bash
# 1. Fazer login em server-1
docker-compose run --rm -e SERVER_URL=tcp://server-1:5555 client npm start
# Login como "teste_replicacao"

# 2. Verificar replicação nos logs
docker-compose logs server-1 | grep "Replicando"
docker-compose logs server-2 | grep "replicado"
docker-compose logs server-3 | grep "replicado"

# 3. Verificar dados em todos os servidores
docker exec messaging-server-1 cat /data/server_data.json | grep "teste_replicacao"
docker exec messaging-server-2 cat /data/server_data.json | grep "teste_replicacao"
docker exec messaging-server-3 cat /data/server_data.json | grep "teste_replicacao"

# Todos devem ter o usuário!
```

---

## 📊 Logs e Debug

```bash
# Ver dados persistidos
docker exec messaging-server-1 cat /data/server_data.json
docker exec messaging-server-2 cat /data/server_data.json
docker exec messaging-server-3 cat /data/server_data.json

# Ver dados do reference
docker exec messaging-reference cat /data/reference_data.json

# Status dos containers
docker-compose ps

# Logs específicos
docker-compose logs -f server-1
docker-compose logs -f reference
docker-compose logs -f broker
```

---

## 🐛 Troubleshooting

### Cliente não conecta
- Verifique containers: `docker-compose ps`
- Veja logs: `docker-compose logs server-1`
- Reinicie: `docker-compose restart server-1`

### Mensagens não chegam
- Verifique broker: `docker-compose ps broker`
- Cliente inscrito no canal? (opção 5)
- Logs do broker: `docker-compose logs broker`

### Replicação não funciona
- Verifique se 3 servidores estão ativos
- Veja logs: `docker-compose logs | grep "Replicando"`
- Verifique coordenador: `docker-compose logs | grep "Coordenador"`

### Erro ao buildar
- Limpe: `docker-compose down -v`
- Rebuild: `docker-compose build --no-cache`

### Dados não persistem
- Verifique volumes: `docker volume ls | grep messaging`
- Veja conteúdo: `docker exec messaging-server-1 ls -la /data`

---

## 👥 Desenvolvimento

**Linguagens utilizadas:**
- **Go** (Server) - Request-Reply, Pub, Relógios, Replicação
- **JavaScript/Node.js** (Client) - CLI interativo
- **Python** (Broker, Reference, Auto-client)

**Padrões implementados:**
- Request-Reply (REQ-REP)
- Publisher-Subscriber (PUB-SUB, XSUB-XPUB)
- Relógio Lógico de Lamport
- Sincronização de Berkeley
- Eleição Bully
- Primary-Backup com Replicação Assíncrona

---

## 📄 Licença

MIT

---

## 🎉 Status do Projeto

✅ **Parte 1**: Request-Reply - COMPLETA  
✅ **Parte 2**: Publisher-Subscriber - COMPLETA  
✅ **Parte 3**: MessagePack - COMPLETA  
✅ **Parte 4**: Relógios (5 etapas) - COMPLETA  
✅ **Parte 5**: Consistência e Replicação - COMPLETA  

**Projeto 100% Concluído! 🎊**