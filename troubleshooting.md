# 🔧 Guia de Resolução de Problemas

## Problemas Comuns e Soluções

### 1. Containers não iniciam

**Sintomas:**
- `docker-compose up` falha
- Containers ficam em estado "Restarting"

**Verificar:**
```bash
# Ver status dos containers
docker-compose ps

# Ver logs de erro
docker-compose logs

# Ver logs específicos
docker-compose logs server
docker-compose logs broker
```

**Soluções:**
```bash
# Reconstruir sem cache
docker-compose build --no-cache

# Limpar e reiniciar
make clean
make rebuild

# Verificar portas em uso
sudo lsof -i :5555
sudo lsof -i :5557
sudo lsof -i :5558
```

### 2. Cliente não conecta ao servidor

**Sintomas:**
- Erro: "Erro ao conectar ao servidor"
- Timeout na conexão

**Verificar:**
```bash
# Servidor está rodando?
docker ps | grep messaging-server

# Servidor está escutando?
docker exec messaging-server netstat -an | grep 5555

# Cliente consegue pingar o servidor?
docker exec messaging-client ping -c 3 server
```

**Soluções:**
```bash
# Reiniciar servidor
docker-compose restart server

# Verificar variável de ambiente
docker exec messaging-client env | grep SERVER_URL

# Deve ser: SERVER_URL=tcp://server:5555
```

### 3. Mensagens não chegam (Pub/Sub)

**Sintomas:**
- Publicações são bem-sucedidas mas ninguém recebe
- Cliente inscrito não vê mensagens

**Verificar:**
```bash
# Broker está rodando?
docker ps | grep messaging-broker

# Broker está escutando nas portas corretas?
docker exec messaging-broker netstat -an | grep 555

# Ver logs do broker
docker-compose logs -f broker
```

**Checklist:**
1. ✅ Cliente fez login?
2. ✅ Cliente se inscreveu no canal? (opção 5)
3. ✅ Canal existe no servidor?
4. ✅ Broker está recebendo mensagens do servidor?

**Soluções:**
```bash
# Reiniciar broker
docker-compose restart broker

# Verificar conexão servidor -> broker
docker exec messaging-server ping -c 3 broker

# Ver se servidor está conectado ao broker
docker-compose logs server | grep "Socket PUB conectado"
```

### 4. Clientes automatizados não funcionam

**Sintomas:**
- Bots não aparecem nos logs
- Nenhuma mensagem automática

**Verificar:**
```bash
# Containers dos bots estão rodando?
docker ps | grep auto-client

# Ver logs detalhados
docker-compose logs -f auto-client-1
docker-compose logs -f auto-client-2
```

**Soluções:**
```bash
# Reiniciar bots
docker-compose restart auto-client-1 auto-client-2

# Verificar se conseguem conectar ao servidor
docker exec messaging-auto-client-1 ping -c 3 server

# Reconstruir imagem dos bots
docker-compose build auto-client-1
docker-compose up -d auto-client-1
```

### 5. Persistência não funciona

**Sintomas:**
- Dados não são salvos após restart
- Arquivo JSON não existe ou está vazio

**Verificar:**
```bash
# Volume existe?
docker volume ls | grep messaging

# Diretório /data existe no container?
docker exec messaging-server ls -la /data

# Arquivo de dados existe?
docker exec messaging-server ls -la /data/server_data.json
```

**Soluções:**
```bash
# Ver conteúdo do arquivo
make data

# Verificar permissões
docker exec messaging-server ls -la /data

# Recriar volume
docker-compose down -v
docker-compose up -d

# Verificar se server tem permissão de escrita
docker exec messaging-server touch /data/test.txt
docker exec messaging-server rm /data/test.txt
```

### 6. Erro "Address already in use"

**Sintomas:**
- Erro ao iniciar: "bind: address already in use"

**Identificar processo usando a porta:**
```bash
# Linux
sudo lsof -i :5555
sudo lsof -i :5557
sudo lsof -i :5558

# macOS
sudo lsof -i -P | grep 5555

# Windows (PowerShell)
netstat -ano | findstr 5555
```

**Soluções:**
```bash
# Matar processo específico
kill -9 <PID>

# Parar todos containers Docker
docker-compose down

# Ou usar portas diferentes no docker-compose.yml
ports:
  - "5565:5555"  # Mapear porta externa diferente
```

### 7. Cliente travado / não responde

**Sintomas:**
- Menu não aparece
- Input não funciona

**Verificar:**
```bash
# Container está rodando?
docker ps | grep messaging-client

# Ver logs
docker-compose logs client
```

**Soluções:**
```bash
# Acessar container e reiniciar cliente
docker exec -it messaging-client sh
npm start

# Ou criar novo cliente
make client-new

# Forçar restart do container
docker-compose restart client
```

### 8. Erro de build

**Sintomas:**
- `docker-compose build` falha
- Erro ao instalar dependências

**Soluções para Go (server):**
```bash
# Limpar cache de módulos
cd server
go clean -modcache
go mod download

# Ou reconstruir sem cache
docker-compose build --no-cache server
```

**Soluções para Node.js (client):**
```bash
# Limpar node_modules
cd client
rm -rf node_modules package-lock.json
npm install

# Ou reconstruir sem cache
docker-compose build --no-cache client
```

**Soluções para Python (broker/auto-client):**
```bash
# Reconstruir sem cache
docker-compose build --no-cache broker
docker-compose build --no-cache auto-client-1
```

### 9. Logs não aparecem

**Sintomas:**
- `docker-compose logs` não mostra nada
- Logs antigos não aparecem

**Soluções:**
```bash
# Ver logs em tempo real
docker-compose logs -f

# Ver logs com timestamp
docker-compose logs -t

# Ver últimas 100 linhas
docker-compose logs --tail=100

# Ver logs de container específico
docker logs messaging-server
docker logs messaging-broker
```

### 10. Network não funciona

**Sintomas:**
- Containers não se comunicam
- DNS não resolve nomes

**Verificar:**
```bash
# Network existe?
docker network ls | grep messaging

# Containers estão na network?
docker network inspect messaging-network

# DNS funciona?
docker exec messaging-client ping server
docker exec messaging-client ping broker
```

**Soluções:**
```bash
# Recriar network
docker-compose down
docker network rm messaging-network
docker-compose up -d

# Verificar configuração
docker network inspect messaging-network
```

## Comandos Úteis para Debug

### Monitoramento em Tempo Real

```bash
# CPU e memória de todos containers
docker stats

# Logs combinados
docker-compose logs -f | grep -E "ERROR|WARNING|✗|❌"

# Ver apenas erros
docker-compose logs 2>&1 | grep -i error
```

### Inspeção de Containers

```bash
# Detalhes do container
docker inspect messaging-server

# Variáveis de ambiente
docker exec messaging-server env

# Processos rodando
docker exec messaging-server ps aux

# Conexões de rede
docker exec messaging-server netstat -an
```

### Acesso Interativo

```bash
# Shell no servidor (Go - Alpine)
docker exec -it messaging-server sh

# Shell no cliente (Node - Alpine)
docker exec -it messaging-client sh

# Shell no broker (Python - Alpine)
docker exec -it messaging-broker sh

# Python interativo no auto-client
docker exec -it messaging-auto-client-1 python
```

### Limpeza Completa

```bash
# Parar tudo e remover volumes
docker-compose down -v

# Remover imagens também
docker-compose down -v --rmi all

# Limpar sistema Docker completo (CUIDADO!)
docker system prune -a --volumes
```

## Testes de Diagnóstico

### Teste 1: Conectividade Básica

```bash
# De dentro do cliente
docker exec messaging-client sh -c "
  echo 'Testando conectividade...'
  ping -c 2 server && echo '✓ Server OK' || echo '✗ Server falhou'
  ping -c 2 broker && echo '✓ Broker OK' || echo '✗ Broker falhou'
"
```

### Teste 2: Portas

```bash
# Verificar se portas estão abertas
for port in 5555 5557 5558; do
  nc -zv localhost $port && echo "✓ Porta $port OK" || echo "✗ Porta $port falhou"
done
```

### Teste 3: ZeroMQ

```bash
# No servidor, verificar sockets ZeroMQ
docker exec messaging-server sh -c "
  netstat -an | grep 5555 && echo '✓ REP socket OK'
"

# No broker
docker exec messaging-broker sh -c "
  netstat -an | grep 5557 && echo '✓ XSUB socket OK'
  netstat -an | grep 5558 && echo '✓ XPUB socket OK'
"
```

## Logs Importantes

### Inicialização Bem-Sucedida

**Server:**
```
🚀 Iniciando servidor...
📊 Dados carregados: X logins, Y canais, ...
📡 Socket REP escutando na porta 5555...
🔌 Socket PUB conectado ao broker em tcp://broker:5557
✅ Servidor pronto para receber requisições!
```

**Broker:**
```
🚀 Iniciando Broker Pub/Sub...
📥 XSUB vinculado na porta 5557
📤 XPUB vinculado na porta 5558
✅ Broker pronto para rotear mensagens!
```

**Cliente:**
```
🚀 Iniciando cliente de mensagens...
🔌 Conectando ao servidor: tcp://server:5555
🔌 Conectando ao broker: tcp://broker:5558
✅ Conectado ao servidor e broker com sucesso!
```

## Quando Pedir Ajuda

Se depois de seguir este guia o problema persistir, colete as seguintes informações:

```bash
# Sistema operacional
uname -a

# Versão do Docker
docker --version
docker-compose --version

# Status dos containers
docker-compose ps

# Todos os logs
docker-compose logs > logs.txt

# Configuração de rede
docker network inspect messaging-network > network.txt

# Dados persistidos
docker exec messaging-server cat /data/server_data.json > data.txt 2>&1
```

E descreva:
1. O que você está tentando fazer
2. O que acontece (erro exato)
3. O que você já tentou
4. Anexe os arquivos de log

## Recursos Adicionais

- [Docker Troubleshooting](https://docs.docker.com/config/containers/troubleshooting/)
- [ZeroMQ FAQ](https://zeromq.org/socket-api/#faq)
- [Docker Compose Networking](https://docs.docker.com/compose/networking/)
