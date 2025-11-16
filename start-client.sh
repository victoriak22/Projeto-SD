#!/bin/bash

echo "🚀 Iniciando Sistema de Mensagens..."

# Para tudo primeiro
echo "🛑 Parando containers existentes..."
docker-compose down

# Sobe os serviços de backend
echo "⚙️  Iniciando serviços de backend..."
docker-compose up -d broker server-1 server-2 server-3 reference auto-client-1 auto-client-2

# Aguarda inicialização
echo "⏳ Aguardando inicialização (5 segundos)..."
sleep 5

# Mostra status
echo ""
echo "✅ Serviços iniciados!"
echo ""
docker-compose ps

echo ""
echo "🎮 Iniciando cliente interativo..."
echo "   (Para sair: Ctrl+C ou digite 0 no menu)"
echo ""

# Roda o cliente interativo
docker-compose run --rm client

echo ""
echo "👋 Cliente encerrado!"
