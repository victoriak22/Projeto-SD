#!/bin/bash

# Script de teste para as Partes 1 e 2 do projeto

set -e

echo "🧪 Iniciando testes das Partes 1 e 2"
echo "========================================================"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para testar se container está rodando
check_container() {
    container_name=$1
    if docker ps | grep -q $container_name; then
        echo -e "${GREEN}✓${NC} Container $container_name está rodando"
        return 0
    else
        echo -e "${RED}✗${NC} Container $container_name NÃO está rodando"
        return 1
    fi
}

# Função para verificar logs
check_logs() {
    container_name=$1
    search_text=$2
    if docker logs $container_name 2>&1 | grep -q "$search_text"; then
        echo -e "${GREEN}✓${NC} Log encontrado em $container_name: $search_text"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} Log não encontrado em $container_name: $search_text"
        return 1
    fi
}

# Função para verificar portas
check_port() {
    port=$1
    service=$2
    if docker ps | grep -q "0.0.0.0:$port"; then
        echo -e "${GREEN}✓${NC} Porta $port ($service) está exposta"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} Porta $port ($service) não está exposta"
        return 1
    fi
}

echo ""
echo "1️⃣  Limpando ambiente..."
docker-compose down -v > /dev/null 2>&1 || true
sleep 2

echo ""
echo "2️⃣  Construindo containers..."
docker-compose build --quiet
echo -e "${GREEN}✓${NC} Build concluído"

echo ""
echo "3️⃣  Iniciando sistema..."
docker-compose up -d
echo "Aguardando inicialização completa..."
sleep 8

echo ""
echo "4️⃣  Verificando containers..."
all_containers_ok=true
check_container "messaging-server" || all_containers_ok=false
check_container "messaging-broker" || all_containers_ok=false
check_container "messaging-client" || all_containers_ok=false
check_container "messaging-auto-client-1" || all_containers_ok=false
check_container "messaging-auto-client-2" || all_containers_ok=false

echo ""
echo "5️⃣  Verificando portas..."
check_port "5555" "Server REQ-REP"
check_port "5557" "Broker XSUB"
check_port "5558" "Broker XPUB"

echo ""
echo "6️⃣  Verificando logs do servidor..."
check_logs "messaging-server" "Servidor pronto"
check_logs "messaging-server" "Socket REP escutando"
check_logs "messaging-server" "Socket PUB conectado"

echo ""
echo "7️⃣  Verificando logs do broker..."
check_logs "messaging-broker" "Broker pronto"
check_logs "messaging-broker" "XSUB vinculado"
check_logs "messaging-broker" "XPUB vinculado"

echo ""
echo "8️⃣  Verificando clientes automatizados..."
sleep 5  # Aguardar bots iniciarem
check_logs "messaging-auto-client-1" "Login realizado" || echo -e "${YELLOW}⚠${NC} Bot 1 ainda não fez login"
check_logs "messaging-auto-client-2" "Login realizado" || echo -e "${YELLOW}⚠${NC} Bot 2 ainda não fez login"

echo ""
echo "9️⃣  Verificando conectividade..."
if docker exec messaging-client ping -c 1 server > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Cliente consegue alcançar o servidor"
else
    echo -e "${RED}✗${NC} Cliente NÃO consegue alcançar o servidor"
fi

if docker exec messaging-client ping -c 1 broker > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Cliente consegue alcançar o broker"
else
    echo -e "${RED}✗${NC} Cliente NÃO consegue alcançar o broker"
fi

echo ""
echo "🔟 Verificando persistência..."
if docker exec messaging-server test -d /data; then
    echo -e "${GREEN}✓${NC} Diretório de dados existe"
    
    # Verificar se arquivo de dados foi criado
    if docker exec messaging-server test -f /data/server_data.json; then
        echo -e "${GREEN}✓${NC} Arquivo de dados criado"
        
        # Mostrar estrutura dos dados
        echo ""
        echo -e "${BLUE}📊 Estrutura dos dados:${NC}"
        docker exec messaging-server cat /data/server_data.json 2>/dev/null | head -20
    else
        echo -e "${YELLOW}⚠${NC} Arquivo de dados ainda não criado"
    fi
else
    echo -e "${RED}✗${NC} Diretório de dados NÃO existe"
fi

echo ""
echo "1️⃣1️⃣  Aguardando atividade dos bots (15 segundos)..."
sleep 15

echo ""
echo "1️⃣2️⃣  Verificando se bots estão publicando mensagens..."
if docker logs messaging-auto-client-1 2>&1 | grep -q "Publicado"; then
    echo -e "${GREEN}✓${NC} Bot 1 está publicando mensagens"
    bot1_msg_count=$(docker logs messaging-auto-client-1 2>&1 | grep -c "Publicado" || echo "0")
    echo -e "   ${BLUE}→${NC} Mensagens enviadas: $bot1_msg_count"
else
    echo -e "${YELLOW}⚠${NC} Bot 1 ainda não publicou mensagens"
fi

if docker logs messaging-auto-client-2 2>&1 | grep -q "Publicado"; then
    echo -e "${GREEN}✓${NC} Bot 2 está publicando mensagens"
    bot2_msg_count=$(docker logs messaging-auto-client-2 2>&1 | grep -c "Publicado" || echo "0")
    echo -e "   ${BLUE}→${NC} Mensagens enviadas: $bot2_msg_count"
else
    echo -e "${YELLOW}⚠${NC} Bot 2 ainda não publicou mensagens"
fi

echo ""
echo "1️⃣3️⃣  Verificando se servidor está processando publicações..."
if docker logs messaging-server 2>&1 | grep -q "Publicação no canal"; then
    echo -e "${GREEN}✓${NC} Servidor está processando publicações"
    pub_count=$(docker logs messaging-server 2>&1 | grep -c "Publicação no canal" || echo "0")
    echo -e "   ${BLUE}→${NC} Publicações processadas: $pub_count"
else
    echo -e "${YELLOW}⚠${NC} Nenhuma publicação processada ainda"
fi

echo ""
echo "1️⃣4️⃣  Resumo da arquitetura:"
echo ""
echo "   📦 Containers ativos:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep messaging

echo ""
echo "   🌐 Rede:"
docker network inspect messaging-network --format '{{range .Containers}}   • {{.Name}} ({{.IPv4Address}})
{{end}}'

echo ""
echo "   📁 Volumes:"
docker volume ls | grep messaging

echo ""
echo "========================================================"
echo "✅ Testes básicos concluídos!"
echo ""
echo -e "${BLUE}🎯 Próximos passos para teste manual:${NC}"
echo ""
echo "1. Testar cliente interativo (3 terminais):"
echo "   Terminal 1: make client"
echo "   Terminal 2: make client-new"
echo "   Terminal 3: make client-new"
echo ""
echo "2. Fluxo de teste sugerido:"
echo "   • Alice: login → criar canal 'geral' → inscrever → publicar"
echo "   • Bob: login → inscrever 'geral' → ver mensagens → enviar DM"
echo "   • Charlie: login → listar usuários → inscrever canais → observar"
echo ""
echo "3. Monitorar atividade:"
echo "   $ make logs-broker    # Ver distribuição de mensagens"
echo "   $ make logs-server    # Ver processamento"
echo "   $ make logs-auto      # Ver atividade dos bots"
echo ""
echo "4. Verificar dados salvos:"
echo "   $ make data           # Ver JSON com todas as mensagens"
echo ""
echo "5. Parar sistema:"
echo "   $ make down"
echo ""
echo "========================================================"

# Status final
if [ "$all_containers_ok" = true ]; then
    echo -e "${GREEN}🎉 Sistema funcionando corretamente!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Alguns containers apresentaram problemas${NC}"
    exit 1
fi