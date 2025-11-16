# 🏗️ Arquitetura do Sistema

## Visão Geral

O sistema é composto por 5 componentes principais que se comunicam usando ZeroMQ:

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Cliente   │◄───────►│   Servidor  │────────►│   Broker    │
│  (Node.js)  │  REQ    │    (Go)     │   PUB   │  (Python)   │
│  REQ + SUB  │  REP    │  REP + PUB  │  XSUB   │ XSUB + XPUB │
└─────────────┘         └─────────────┘         └─────────────┘
      ▲                                                 │
      │                                                 │ XPUB
      │                 ┌─────────────┐                │
      └─────────────────┤ Auto Client │◄───────────────┘
            SUB         │  (Python)   │
                        │   REQ + SUB  │
                        └─────────────┘
```

## Componentes

### 1. Servidor (Go)

**Responsabilidades:**
- Processar requisições de login, cadastro de canais
- Receber requisições de publicação e mensagens diretas
- Publicar no broker (socket PUB)
- Persistir todos os dados em JSON

**Sockets:**
- `REP` na porta 5555 - Recebe requisições dos clientes
- `PUB` conectado ao broker:5557 - Publica mensagens

**Serviços:**
- `login` - Cadastro de usuários
- `users` - Listagem de usuários
- `channel` - Criação de canais
- `channels` - Listagem de canais
- `publish` - Publicação em canal
- `message` - Mensagem direta

### 2. Broker (Python)

**Responsabilidades:**
- Atuar como proxy entre publishers e subscribers
- Rotear mensagens baseado em tópicos
- Gerenciar inscrições de clientes

**Sockets:**
- `XSUB` na porta 5557 - Recebe de publishers (servidor)
- `XPUB` na porta 5558 - Distribui para subscribers (clientes)

**Funcionamento:**
```python
# Proxy simples que conecta os dois sockets
zmq.proxy(xsub, xpub)

# XSUB recebe:
# - Mensagens dos publishers (servidor)
# - Inscrições dos subscribers (via XPUB)

# XPUB envia:
# - Mensagens para subscribers
# - Notificações de inscrição para XSUB
```

### 3. Cliente Interativo (Node.js)

**Responsabilidades:**
- Interface CLI para o usuário
- Enviar requisições ao servidor
- Receber mensagens publicadas no broker

**Sockets:**
- `REQ` conectado ao server:5555 - Envia requisições
- `SUB` conectado ao broker:5558 - Recebe mensagens

**Tópicos de Inscrição:**
- Nome do próprio usuário (mensagens diretas)
- Nomes dos canais que o usuário escolheu

### 4. Cliente Automatizado (Python)

**Responsabilidades:**
- Gerar mensagens aleatórias para testes
- Criar carga no sistema
- Validar funcionamento do Pub/Sub

**Comportamento:**
1. Gera username aleatório e faz login
2. Cria canais iniciais se não existirem
3. Loop infinito:
   - Escolhe canal aleatório
   - Envia 10 mensagens
   - Pausa 5-10 segundos
   - Repete

### 5. Proxy (Futuro)

**Status:** Placeholder para próximas partes
**Possíveis funcionalidades:**
- Cache de mensagens
- Balanceamento de carga
- Roteamento inteligente

## Fluxos de Comunicação

### Fluxo 1: Login

```
Cliente ─────REQ────►  Servidor
        {service: "login", data: {user, timestamp}}

Cliente ◄────REP─────  Servidor
        {service: "login", data: {status, timestamp}}

[Servidor persiste login em JSON]
```

### Fluxo 2: Publicação em Canal

```
1. Cliente envia requisição
Cliente ─────REQ────►  Servidor
        {service: "publish", data: {user, channel, message, timestamp}}

2. Servidor valida e publica no broker
Servidor ────PUB────►  Broker (tópico = channel)
        {user, message, timestamp}

3. Broker distribui para subscribers
Broker ──────XPUB───►  Clientes SUB (inscritos no canal)

4. Servidor responde ao cliente original
Cliente ◄────REP─────  Servidor
        {service: "publish", data: {status: "OK", timestamp}}

5. Servidor persiste mensagem
[JSON: channel_messages array]
```

### Fluxo 3: Mensagem Direta

```
1. Alice envia para Bob
Alice ───────REQ────►  Servidor
        {service: "message", data: {src: "alice", dst: "bob", message}}

2. Servidor publica no tópico do Bob
Servidor ────PUB────►  Broker (tópico = "bob")
        {from: "alice", message, timestamp}

3. Bob recebe (se estiver inscrito)
Broker ──────XPUB───►  Bob (SUB no tópico "bob")

4. Servidor confirma para Alice
Alice ◄──────REP─────  Servidor
        {service: "message", data: {status: "OK"}}

5. Servidor persiste
[JSON: user_messages array]
```

### Fluxo 4: Inscrição em Canal

```
1. Cliente se inscreve localmente
cliente.subSocket.subscribe("geral")

2. ZeroMQ envia mensagem de inscrição
Cliente ─────SUB────►  Broker
        [mensagem de controle do ZeroMQ]

3. Broker roteia para o servidor (via XSUB)
Broker ──────────────►  Servidor
        [inscrição propagada automaticamente]

4. A partir deste momento, cliente recebe mensagens do canal
```

## Padrões de Mensagens

**Nota**: A partir da Parte 3, todas as mensagens são serializadas em **MessagePack** (formato binário) ao invés de JSON.

### Tópicos no Broker

O broker usa dois tipos de tópicos:

1. **Canais públicos**: Nome do canal
   - Exemplo: `"geral"`, `"tech"`, `"random"`
   - Qualquer cliente inscrito recebe

2. **Mensagens diretas**: Nome do usuário de destino
   - Exemplo: `"alice"`, `"bob"`
   - Apenas o usuário específico recebe

### Formato das Publicações

**Formato lógico** (serializadas em MessagePack):

**Publicação em Canal:**
```json
{
  "user": "alice",
  "message": "Olá pessoal!",
  "timestamp": 1234567890
}
```

**Mensagem Direta:**
```json
{
  "from": "alice",
  "message": "Oi Bob!",
  "timestamp": 1234567890
}
```

## Persistência

### Arquivo: `/data/server_data.json`

```json
{
  "logins": [
    {
      "username": "alice",
      "timestamp": 1234567890
    }
  ],
  "channels": ["geral", "tech", "random"],
  "channel_messages": [
    {
      "user": "alice",
      "channel": "geral",
      "message": "Olá!",
      "timestamp": 1234567890
    }
  ],
  "user_messages": [
    {
      "src": "alice",
      "dst": "bob",
      "message": "Oi Bob!",
      "timestamp": 1234567890
    }
  ]
}
```

## Portas

| Serviço | Porta | Tipo | Descrição |
|---------|-------|------|-----------|
| Server | 5555 | REP | Requisições dos clientes |
| Broker | 5557 | XSUB | Recebe de publishers |
| Broker | 5558 | XPUB | Distribui para subscribers |

## Vantagens da Arquitetura

1. **Desacoplamento**: Clientes não precisam saber uns dos outros
2. **Escalabilidade**: Múltiplos publishers e subscribers
3. **Confiabilidade**: Broker centralizado gerencia distribuição
4. **Flexibilidade**: Fácil adicionar novos clientes
5. **Persistência**: Todas as mensagens são armazenadas

## Limitações Atuais

1. **Single Point of Failure**: Broker único
2. **Sem Garantia de Entrega**: Pub/Sub é best-effort
3. **Sem Ordenação Global**: Apenas timestamps locais
4. **Sem Autenticação**: Usuários não precisam de senha

## Próximas Melhorias

- **Parte 3**: MessagePack para serialização mais eficiente
- **Parte 4**: Relógios lógicos para ordenação correta
- **Parte 5**: Replicação do broker para alta disponibilidade

## Debugging

### Ver fluxo completo de uma mensagem

```bash
# Terminal 1: Broker
docker-compose logs -f broker

# Terminal 2: Server
docker-compose logs -f server

# Terminal 3: Cliente
docker-compose logs -f client

# Terminal 4: Enviar mensagem
# No cliente interativo, publicar em canal
```

### Monitorar tráfego ZeroMQ

```bash
# Ver todas conexões ativas
docker exec messaging-server netstat -an | grep 555

# Ver processos ZeroMQ
docker exec messaging-broker ps aux | grep python

# Ver uso de CPU/memória
docker stats
```

## Documentação Adicional

- Ver **README.md** para instruções de instalação e uso
- Consultar a documentação oficial do [ZeroMQ](https://zguide.zeromq.org/)
- Referências sobre [MessagePack](https://msgpack.org/)