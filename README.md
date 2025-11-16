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

## 🚀 Como Usar

### Pré-requisitos

- Docker
- Docker Compose

### Método 1: Inicialização Rápida (Recomendado)

```bash
# 1. Parar tudo (se houver algo rodando)
docker-compose down

# 2. Subir serviços de backend em background
docker-compose up -d broker server-1 server-2 server-3 reference auto-client-1 auto-client-2

# 3. Aguardar inicialização (3-5 segundos)
sleep 5

# 4. Rodar cliente interativo
docker-compose run --rm client
```

### Método 2: Script Automatizado

Salve este conteúdo em `start-client.sh`:

```bash
#!/bin/bash
echo "🚀 Iniciando Sistema de Mensagens..."
docker-compose down
echo "⚙️  Iniciando serviços de backend..."
docker-compose up -d broker server-1 server-2 server-3 reference auto-client-1 auto-client-2
echo "⏳ Aguardando inicialização..."
sleep 5
echo "🎮 Iniciando cliente interativo..."
docker-compose run --rm client
```

Depois:
```bash
chmod +x start-client.sh
./start-client.sh
```

### Método 3: Inicialização Manual Passo a Passo

```bash
# Passo 1: Limpar ambiente
docker-compose down

# Passo 2: Construir imagens (apenas primeira vez ou após mudanças)
docker-compose build

# Passo 3: Iniciar serviços essenciais
docker-compose up -d reference
sleep 2

docker-compose up -d broker
sleep 2

docker-compose up -d server-1 server-2 server-3
sleep 3

# Passo 4: Iniciar clientes automatizados (opcional)
docker-compose up -d auto-client-1 auto-client-2

# Passo 5: Verificar status
docker-compose ps

# Passo 6: Iniciar cliente interativo
docker-compose run --rm client
```

---

## 🎮 Usando o Cliente Interativo

Após executar qualquer dos métodos acima, você verá o menu:

```
============================================================
📱 SISTEMA DE MENSAGENS - MENU PRINCIPAL
============================================================

Opções:
  1. Fazer login
  0. Sair
============================================================
```

### Fluxo Típico de Uso

#### 1. Fazer Login
```
Escolha uma opção: 1
Digite seu nome de usuário: alice
✅ Login realizado com sucesso!
```

#### 2. Criar Canal
```
Escolha uma opção: 3
Digite o nome do canal: geral
✅ Canal criado com sucesso!
```

#### 3. Ver Canais Disponíveis
```
Escolha uma opção: 4

📺 Canais disponíveis:
  - geral
  - tech
  - random
```

#### 4. Inscrever-se em Canal (para receber mensagens)
```
Escolha uma opção: 7
Digite o nome do canal: geral
✅ Inscrito no canal #geral
```

#### 5. Publicar Mensagem em Canal
```
Escolha uma opção: 5
Digite o nome do canal: geral
Digite sua mensagem: Olá pessoal!
✅ Mensagem publicada com sucesso!
```

#### 6. Ver Usuários Online
```
Escolha uma opção: 2

👥 Usuários cadastrados:
  - alice
  - bob
  - charlie
```

#### 7. Enviar Mensagem Direta
```
Escolha uma opção: 6
Digite o nome do usuário: bob
Digite sua mensagem: Oi Bob, tudo bem?
✅ Mensagem enviada com sucesso!
```

#### 8. Sair
```
Escolha uma opção: 0
👋 Encerrando cliente...
```

### Menu Completo

```
============================================================
📱 MENU DO USUÁRIO: alice
============================================================

Opções:
  1. Fazer login novamente
  2. Listar usuários
  3. Criar canal
  4. Listar canais
  5. Publicar em canal
  6. Enviar mensagem direta
  7. Inscrever-se em canal
  0. Sair
============================================================
```

---

## 🧪 Testando o Sistema

### Teste 1: Comunicação Básica (1 Cliente)

```bash
docker-compose run --rm client
```

1. Login como "alice"
2. Criar canal "geral"
3. Inscrever-se no canal "geral"
4. Publicar mensagem "Olá!"
5. Sair

### Teste 2: Múltiplos Clientes (Pub/Sub)

#### Terminal 1 - Alice
```bash
docker-compose run --rm client
```
1. Login: alice
2. Criar canal: geral
3. Inscrever-se: geral
4. Aguardar mensagens...

#### Terminal 2 - Bob
```bash
docker-compose run --rm client
```
1. Login: bob
2. Inscrever-se: geral
3. Publicar: "Oi Alice!"
4. Ver mensagem chegando no Terminal 1

#### Terminal 3 - Charlie
```bash
docker-compose run --rm client
```
1. Login: charlie
2. Enviar DM para alice: "Mensagem privada!"

### Teste 3: Replicação entre Servidores

#### Terminal 1 - Cliente no Server-1
```bash
docker-compose run --rm -e SERVER_URL=tcp://server-1:5555 client
```
1. Login: teste_replicacao
2. Criar canal: canal_teste

#### Terminal 2 - Verificar Server-2
```bash
docker exec messaging-server-2 cat /data/server_data.json | grep teste_replicacao
# Deve mostrar o usuário replicado!
```

#### Terminal 3 - Verificar Server-3
```bash
docker exec messaging-server-3 cat /data/server_data.json | grep canal_teste
# Deve mostrar o canal replicado!
```

### Teste 4: Tolerância a Falhas (Eleição Bully)

```bash
# 1. Ver coordenador atual
docker-compose logs server-3 | grep "Coordenador"

# 2. Parar o coordenador (server-3 com rank 4)
docker-compose stop server-3

# 3. Aguardar 30-40 segundos

# 4. Ver nova eleição nos logs
docker-compose logs server-2 | grep "eleição"
docker-compose logs server-1 | grep "eleição"

# 5. Verificar novo coordenador (deve ser server-2 com rank 2)
docker-compose logs server-2 | grep "coordenador"

# 6. Reiniciar server-3
docker-compose start server-3

# 7. Aguardar e verificar que server-3 volta como coordenador
docker-compose logs server-3 | tail -20
```

### Teste 5: Sincronização Berkeley

```bash
# 1. Fazer 10+ operações para forçar sincronização
docker-compose run --rm client
# Login, criar 3 canais, publicar 5 mensagens, etc.

# 2. Ver logs de sincronização
docker-compose logs server-3 | grep "Berkeley"

# Deve mostrar:
# 🎯 Iniciando sincronização Berkeley como COORDENADOR
# 📊 Coletando timestamps...
# ✅ Sincronização Berkeley concluída
```

### Teste 6: Clientes Automatizados

```bash
# Ver os auto-clients em ação
docker-compose logs -f auto-client-1 auto-client-2

# Deve mostrar:
# ✅ Login realizado: bot_1234
# ✅ Canal criado: geral
# 📤 Publicado em #geral: Olá pessoal!
```

---

## 📊 Monitoramento e Debug

### Ver Logs em Tempo Real

```bash
# Todos os serviços
docker-compose logs -f

# Apenas servidores
docker-compose logs -f server-1 server-2 server-3

# Apenas um servidor
docker-compose logs -f server-1

# Servidor de referência
docker-compose logs -f reference

# Broker
docker-compose logs -f broker

# Clientes automatizados
docker-compose logs -f auto-client-1 auto-client-2
```

### Verificar Status dos Containers

```bash
# Listar todos os containers
docker-compose ps

# Esperado:
# messaging-reference       running   5559/tcp
# messaging-broker          running   5557/tcp, 5558/tcp
# messaging-server-1        running
# messaging-server-2        running
# messaging-server-3        running
# messaging-auto-client-1   running
# messaging-auto-client-2   running
```

### Verificar Dados Persistidos

```bash
# Ver dados do Server-1
docker exec messaging-server-1 cat /data/server_data.json

# Ver dados do Server-2
docker exec messaging-server-2 cat /data/server_data.json

# Ver dados do Server-3
docker exec messaging-server-3 cat /data/server_data.json

# Ver dados do Reference
docker exec messaging-reference cat /data/reference_data.json

# Buscar usuário específico
docker exec messaging-server-1 cat /data/server_data.json | grep "alice"

# Contar logins
docker exec messaging-server-1 cat /data/server_data.json | jq '.logins | length'
```

### Verificar Relógios Lógicos

```bash
# Ver valores de clock nos logs
docker-compose logs server-1 | grep "clock:"

# Ver sincronização Berkeley
docker-compose logs server-3 | grep "Berkeley"

# Ver ajustes de tempo
docker-compose logs | grep "Relógio ajustado"
```

### Verificar Replicação

```bash
# Ver tentativas de replicação
docker-compose logs | grep "Replicando"

# Ver dados replicados recebidos
docker-compose logs | grep "replicado"

# Ver sincronização completa
docker-compose logs | grep "Sincronização"
```

---

## 🛑 Parando o Sistema

### Parar Apenas o Cliente

```bash
# No terminal do cliente, pressione Ctrl+C ou digite 0
```

### Parar Serviços de Backend

```bash
# Parar mantendo dados
docker-compose stop

# Parar e remover containers (mantém volumes)
docker-compose down

# Parar e remover TUDO incluindo dados
docker-compose down -v
```

### Restart de Serviços Específicos

```bash
# Reiniciar um servidor
docker-compose restart server-1

# Reiniciar o broker
docker-compose restart broker

# Reiniciar reference
docker-compose restart reference
```

---

## 🔧 Comandos Úteis

### Reconstruir Após Mudanças no Código

```bash
# Reconstruir tudo
docker-compose build

# Reconstruir sem cache
docker-compose build --no-cache

# Reconstruir apenas um serviço
docker-compose build client
docker-compose build server-1

# Reconstruir e reiniciar
docker-compose up -d --build server-1
```

### Limpar Completamente

```bash
# Parar tudo
docker-compose down -v

# Remover imagens órfãs
docker image prune -f

# Remover volumes órfãos
docker volume prune -f

# Reconstruir do zero
docker-compose build --no-cache
docker-compose up -d broker server-1 server-2 server-3 reference
```

### Acessar Shell de um Container

```bash
# Bash no servidor
docker exec -it messaging-server-1 /bin/sh

# Bash no reference
docker exec -it messaging-reference /bin/sh

# Bash no broker
docker exec -it messaging-broker /bin/sh
```

### Copiar Arquivos

```bash
# Copiar dados do servidor para host
docker cp messaging-server-1:/data/server_data.json ./server1_backup.json

# Copiar código do host para container
docker cp ./server/main.go messaging-server-1:/app/main.go
```

---

## 🐛 Troubleshooting

### Problema: Cliente não aceita entrada

**Sintoma:**
```bash
docker-compose up client
# Menu aparece mas não consigo digitar
```

**Solução:**
```bash
# Use docker-compose run ao invés de up
docker-compose run --rm client
```

### Problema: "Nome de usuário não pode ser vazio"

**Sintoma:**
```bash
auto-client: ❌ Erro no login: Nome de usuário não pode ser vazio
```

**Solução:**
Verifique se o `auto_client.py` usa campos em **minúsculo**:
```python
# ✅ Correto
request = {
    "service": "login",
    "data": {
        "user": username,              # minúsculo!
        "timestamp": int(time.time()),
        "clock": increment_clock()
    }
}
```

### Problema: Containers não conectam

**Sintoma:**
```bash
ERROR: Network messaging-network not found
```

**Solução:**
```bash
docker-compose down
docker network prune -f
docker-compose up -d broker server-1 server-2 server-3 reference
```

### Problema: Mensagens não chegam

**Checklist:**
1. ✅ Broker está rodando? `docker-compose ps broker`
2. ✅ Cliente está inscrito no canal? (opção 7)
3. ✅ Há publishers? `docker-compose logs broker`

**Solução:**
```bash
# Reiniciar broker
docker-compose restart broker

# Ver logs
docker-compose logs -f broker
```

### Problema: Replicação não funciona

**Checklist:**
1. ✅ 3 servidores rodando? `docker-compose ps | grep server`
2. ✅ Reference ativo? `docker-compose ps reference`
3. ✅ Coordenador definido? `docker-compose logs | grep Coordenador`

**Solução:**
```bash
# Ver logs de replicação
docker-compose logs | grep -i "replic"

# Forçar sincronização
docker-compose restart server-1 server-2 server-3
```

### Problema: Porta já em uso

**Sintoma:**
```bash
ERROR: port is already allocated
```

**Solução:**
```bash
# Ver o que está usando a porta
lsof -i :5559  # ou 5557, 5558

# Parar processo
kill -9 <PID>

# Ou mudar porta no docker-compose.yml
ports:
  - "5560:5559"  # usa 5560 externamente
```

### Problema: Dados não persistem

**Sintoma:**
Após reiniciar, todos os dados sumiram.

**Solução:**
```bash
# Verificar volumes
docker volume ls | grep messaging

# NÃO use -v ao parar
docker-compose down  # ✅ mantém volumes
docker-compose down -v  # ❌ APAGA volumes

# Backup manual
docker exec messaging-server-1 cat /data/server_data.json > backup.json
```

---

## 📋 Funcionalidades Implementadas

### Parte 1: Request-Reply ✅
- ✅ Login de usuários
- ✅ Listagem de usuários
- ✅ Criação de canais
- ✅ Listagem de canais
- ✅ Persistência de dados

### Parte 2: Publisher-Subscriber ✅
- ✅ Broker Pub/Sub
- ✅ Publicação em canais
- ✅ Mensagens diretas
- ✅ Inscrição em canais
- ✅ Cliente automatizado
- ✅ Persistência de mensagens

### Parte 3: MessagePack ✅
- ✅ Serialização eficiente
- ✅ Compatibilidade entre linguagens
- ✅ Redução de tamanho (~25%)
- ✅ Melhor performance

### Parte 4: Relógios ✅
- ✅ Relógio Lógico de Lamport
- ✅ Servidor de Referência
- ✅ Múltiplos Servidores (3x)
- ✅ Sincronização Berkeley
- ✅ Eleição Bully

### Parte 5: Consistência e Replicação ✅
- ✅ Replicação automática
- ✅ Sincronização periódica
- ✅ Primary-Backup adaptado
- ✅ Consistência eventual
- ✅ Tolerância a falhas

---

## 🏗️ Arquitetura

```
                    ┌─────────────┐
                    │  Reference  │ :5559
                    │  (Python)   │
                    └─────────────┘
                         ↕ REQ/REP
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
   Server-1         Server-2         Server-3
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

---

## 🛠️ Tecnologias

- **Server**: Go 1.21 + ZeroMQ + MessagePack
- **Client**: Node.js 20 + ZeroMQ + MessagePack
- **Broker**: Python 3.11 + ZeroMQ
- **Reference**: Python 3.11 + ZeroMQ + MessagePack
- **Auto-client**: Python 3.11 + ZeroMQ + MessagePack
- **Containerização**: Docker + Docker Compose

---

## 📁 Estrutura do Projeto

```
.
├── reference/          # Servidor de referência
├── broker/             # Broker Pub/Sub
├── server/             # Servidor (Go)
├── client/             # Cliente interativo + automatizado
├── docker-compose.yml  # Orquestração
└── README.md          # Este arquivo
```

---

## 🎉 Status

✅ **Parte 1**: Request-Reply - COMPLETA  
✅ **Parte 2**: Publisher-Subscriber - COMPLETA  
✅ **Parte 3**: MessagePack - COMPLETA  
✅ **Parte 4**: Relógios - COMPLETA  
✅ **Parte 5**: Replicação - COMPLETA  

---

## 📄 Licença

MIT
