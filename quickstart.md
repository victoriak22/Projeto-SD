# 🚀 Início Rápido - Partes 1 e 2

## ⚡ Comandos Essenciais

```bash
# 1. Iniciar o sistema (inclui broker e clientes automatizados)
make up

# 2. Abrir cliente (em 3 terminais diferentes)
make client        # Terminal 1
make client-new    # Terminal 2
make client-new    # Terminal 3

# 3. Ver logs
make logs          # Todos
make logs-server   # Apenas servidor
make logs-broker   # Apenas broker
make logs-auto     # Clientes automatizados

# 4. Ver dados salvos
make data

# 5. Parar sistema
make down
```

## 📝 Fluxo de Teste Rápido - Parte 1

### Terminal 1 - Alice
```
Escolha: 1
Nome: alice
Escolha: 3
Canal: geral
Escolha: 4 (ver canais)
```

### Terminal 2 - Bob
```
Escolha: 1
Nome: bob
Escolha: 2 (ver usuários - deve mostrar alice e bob)
Escolha: 3
Canal: tech
Escolha: 4 (ver canais - deve mostrar geral e tech)
```

## 📝 Fluxo de Teste Completo - Parte 2 (Pub/Sub)

### Terminal 1 - Alice
```
Escolha: 1
Nome: alice
Escolha: 5 (inscrever em canal)
Canal: geral
Escolha: 6 (publicar mensagem)
Canal: geral
Mensagem: Olá pessoal!
[Aguardar e ver mensagens dos bots chegando]
```

### Terminal 2 - Bob
```
Escolha: 1
Nome: bob
Escolha: 5 (inscrever em canal)
Canal: geral
[Deve ver mensagens de Alice e dos bots]
Escolha: 7 (enviar mensagem direta)
Destinatário: alice
Mensagem: Oi Alice, tudo bem?
```

### Terminal 3 - Ver Bots
```bash
# Em outro terminal
make logs-auto

# Você verá os bots:
# - bot_XXXX enviando mensagens automaticamente
# - auto_YYYY publicando nos canais
# - Mensagens sendo enviadas a cada 1-3 segundos
```

## ✅ Checklist de Validação

**Parte 1:**
- [ ] Sistema inicia sem erros
- [ ] Múltiplos clientes podem se conectar
- [ ] Login funciona e impede duplicatas
- [ ] Lista de usuários mostra todos cadastrados
- [ ] Criação de canais funciona
- [ ] Lista de canais mostra todos criados
- [ ] Dados persistem após restart

**Parte 2:**
- [ ] Broker inicia e conecta publisher/subscribers
- [ ] Clientes podem se inscrever em canais
- [ ] Publicações em canais chegam aos inscritos
- [ ] Mensagens diretas funcionam
- [ ] Clientes automatizados estão ativos
- [ ] Múltiplos clientes recebem mesma mensagem
- [ ] Mensagens são persistidas
- [ ] Timestamps corretos em todas mensagens

## 🎯 Testes Específicos

### Teste 1: Publicação em Canal
```bash
# Terminal 1
1. Login: alice
5. Inscrever canal: geral
6. Publicar: "Primeira mensagem!"

# Terminal 2
1. Login: bob
5. Inscrever canal: geral
# Bob deve ver: "alice: Primeira mensagem!"
```

### Teste 2: Mensagem Direta
```bash
# Terminal 1 (Alice)
7. Mensagem direta
   Destinatário: bob
   Mensagem: "Oi Bob!"

# Terminal 2 (Bob) - deve receber imediatamente
💬 [DM de alice]: Oi Bob!
```

### Teste 3: Múltiplas Inscrições
```bash
# Um cliente pode se inscrever em vários canais
1. Login: charlie
5. Inscrever: geral
5. Inscrever: tech
5. Inscrever: random
# Charlie receberá mensagens de todos esses canais
```

### Teste 4: Clientes Automatizados
```bash
# Verificar se os bots estão funcionando
make logs-auto

# Deve mostrar:
# ✅ Login realizado: bot_XXXX
# 📤 Publicado em #geral: ...
# 📤 Publicado em #random: ...
```

## 🐛 Resolução Rápida de Problemas

**Erro ao conectar ao broker?**
```bash
# Verificar se broker está rodando
docker ps | grep broker

# Ver logs
make logs-broker
```

**Não recebe mensagens?**
```bash
# Verificar se está inscrito no canal
# No cliente, opção 4 para ver canais
# Depois opção 5 para se inscrever
```

**Clientes automatizados não funcionam?**
```bash
# Verificar containers
docker ps | grep auto-client

# Reiniciar
docker-compose restart auto-client-1 auto-client-2
```

**Limpar tudo e recomeçar?**
```bash
make clean
make rebuild
```

## 📊 Estrutura das Pastas (Atualizada)

```
projeto/
├── server/              # Go - servidor REP + PUB
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── client/              # Node.js + Python
│   ├── main.js         # Cliente interativo (Node.js)
│   ├── auto_client.py  # Cliente automatizado (Python)
│   ├── package.json
│   ├── Dockerfile
│   └── Dockerfile.auto
├── broker/              # Python - proxy XSUB/XPUB
│   ├── main.py
│   ├── requirements.txt
│   └── Dockerfile
├── proxy/               # (Próximas partes)
└── docker-compose.yml
```

## 🎯 Objetivos Completos

**Parte 1:**
- ✅ Request-Reply com ZeroMQ
- ✅ Login de usuários
- ✅ Listagem de usuários
- ✅ Criação de canais
- ✅ Listagem de canais
- ✅ Persistência em JSON
- ✅ Go + JavaScript

**Parte 2:**
- ✅ Broker Pub/Sub (XSUB/XPUB)
- ✅ Publicação em canais
- ✅ Mensagens diretas
- ✅ Inscrição em canais
- ✅ Cliente automatizado
- ✅ Persistência de mensagens
- ✅ Go + JavaScript + Python (3 linguagens)

## 🔄 Workflow Git Sugerido

```bash
# Após validar Parte 2
git checkout -b part-2
git add .
git commit -m "feat: implementa Parte 2 - Publisher-Subscriber

- Broker Python com XSUB/XPUB
- Publicação em canais e mensagens diretas
- Cliente automatizado para testes
- Inscrição em canais
- Persistência de todas mensagens"

git checkout main
git merge part-2
git push origin main
```

## 📚 Dicas de Uso

1. **Para testar a comunicação**, abra vários terminais e veja as mensagens fluindo em tempo real

2. **Use os bots** para simular carga: eles criarão canais e enviarão mensagens automaticamente

3. **Monitore os logs** em tempo real para entender o fluxo:
   ```bash
   # Terminal 1: Broker
   make logs-broker
   
   # Terminal 2: Server
   make logs-server
   
   # Terminal 3: Bots
   make logs-auto
   ```

4. **Persistência**: Todas as mensagens ficam salvas em `/data/server_data.json` dentro do container do servidor

---

**Dica:** Use `make help` para ver todos os comandos disponíveis!