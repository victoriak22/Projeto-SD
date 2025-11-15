#!/usr/bin/env python3
"""
Cliente Automatizado - Parte 2
Gera mensagens aleatórias em canais para testes
"""

import zmq
import msgpack
import time
import random
import logging
import sys
from datetime import datetime

# Configuração de logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger(__name__)

# Relógio lógico
logical_clock = 0

def increment_clock():
    """Incrementa e retorna o relógio lógico"""
    global logical_clock
    logical_clock += 1
    return logical_clock

def update_clock(received_clock):
    """Atualiza o relógio lógico ao receber uma mensagem"""
    global logical_clock
    if received_clock > logical_clock:
        logical_clock = received_clock
    logical_clock += 1
    return logical_clock

# Lista de mensagens para enviar
MESSAGES = [
    "Olá pessoal!",
    "Como estão todos?",
    "Alguém por aqui?",
    "Ótimo dia para programar!",
    "Sistemas distribuídos são incríveis!",
    "ZeroMQ funciona muito bem!",
    "Este projeto está ficando legal",
    "Python + Go + Node.js = ❤️",
    "Adorei este sistema!",
    "Mensagem de teste automática"
]

# Nomes de usuários aleatórios
USER_PREFIXES = ["bot", "auto", "test", "user", "client"]


def generate_username():
    """Gera um nome de usuário aleatório"""
    prefix = random.choice(USER_PREFIXES)
    number = random.randint(1000, 9999)
    return f"{prefix}_{number}"


def send_request(socket, request):
    """Envia requisição e recebe resposta"""
    try:
        # Adicionar clock antes de enviar
        if 'data' in request:
            request['data']['clock'] = increment_clock()
        
        packed = msgpack.packb(request)
        socket.send(packed)
        response_data = socket.recv()
        response = msgpack.unpackb(response_data, raw=False)
        
        # Atualizar clock ao receber resposta
        if 'data' in response and 'clock' in response['data']:
            update_clock(response['data']['clock'])
        
        return response
    except Exception as e:
        logger.error(f"Erro na comunicação: {e}")
        return None


def login(socket, username):
    """Faz login no sistema"""
    request = {
        "service": "login",
        "data": {
            "user": username,
            "timestamp": int(time.time())
        }
    }
    
    logger.info(f"🔐 Tentando login como: {username}")
    response = send_request(socket, request)
    
    if response and response.get("data", {}).get("status") == "sucesso":
        logger.info(f"✅ Login realizado: {username}")
        return True
    else:
        logger.warning(f"❌ Erro no login: {response}")
        return False


def get_channels(socket):
    """Obtém lista de canais disponíveis"""
    request = {
        "service": "channels",
        "data": {
            "timestamp": int(time.time())
        }
    }
    
    response = send_request(socket, request)
    
    if response and "data" in response and "channels" in response["data"]:
        channels = response["data"]["channels"]
        logger.info(f"📺 Canais disponíveis: {channels}")
        return channels
    else:
        logger.warning("⚠️  Nenhum canal encontrado")
        return []


def create_channel(socket, channel_name):
    """Cria um novo canal"""
    request = {
        "service": "channel",
        "data": {
            "channel": channel_name,
            "timestamp": int(time.time())
        }
    }
    
    logger.info(f"🆕 Tentando criar canal: {channel_name}")
    response = send_request(socket, request)
    
    if response and response.get("data", {}).get("status") == "sucesso":
        logger.info(f"✅ Canal criado: {channel_name}")
        return True
    else:
        logger.debug(f"Canal {channel_name} já existe ou erro")
        return False


def publish_message(socket, username, channel, message):
    """Publica mensagem em um canal"""
    request = {
        "service": "publish",
        "data": {
            "user": username,
            "channel": channel,
            "message": message,
            "timestamp": int(time.time())
        }
    }
    
    response = send_request(socket, request)
    
    if response and response.get("data", {}).get("status") == "OK":
        logger.info(f"📤 Publicado em #{channel}: {message[:30]}...")
        return True
    else:
        logger.error(f"❌ Erro ao publicar: {response}")
        return False


def main():
    """Função principal do cliente automatizado"""
    logger.info("🤖 Iniciando cliente automatizado...")
    
    # Configuração
    server_url = os.getenv("SERVER_URL", "tcp://server-1:5555")
    username = generate_username()
    
    # Conectar ao servidor
    context = zmq.Context()
    socket = context.socket(zmq.REQ)
    
    logger.info(f"🔌 Conectando ao servidor: {server_url}")
    socket.connect(server_url)
    
    # Aguardar conexão
    time.sleep(2)
    
    # Fazer login
    if not login(socket, username):
        logger.error("❌ Falha no login. Encerrando...")
        return
    
    # Aguardar um pouco
    time.sleep(1)
    
    # Criar alguns canais iniciais se não existirem
    initial_channels = ["geral", "random", "tech", "bots"]
    for channel in initial_channels:
        create_channel(socket, channel)
        time.sleep(0.5)
    
    logger.info("🔄 Iniciando loop de mensagens...")
    
    # Loop infinito de envio de mensagens
    message_count = 0
    while True:
        try:
            # Obter canais disponíveis
            channels = get_channels(socket)
            
            if not channels:
                logger.warning("⚠️  Nenhum canal disponível. Aguardando...")
                time.sleep(5)
                continue
            
            # Escolher canal aleatório
            channel = random.choice(channels)
            
            # Enviar 10 mensagens
            for i in range(10):
                message = random.choice(MESSAGES)
                
                if publish_message(socket, username, channel, message):
                    message_count += 1
                    logger.info(f"📊 Total: {message_count} msgs | Clock: {logical_clock}")
                
                # Intervalo entre mensagens
                time.sleep(random.uniform(1, 3))
            
            # Pausa maior entre ciclos
            logger.info("⏸️  Pausa entre ciclos...")
            time.sleep(random.uniform(5, 10))
            
        except KeyboardInterrupt:
            logger.info("\n👋 Encerrando cliente automatizado...")
            break
        except Exception as e:
            logger.error(f"❌ Erro no loop: {e}")
            time.sleep(5)
    
    # Fechar socket
    socket.close()
    context.term()
    logger.info("✅ Cliente automatizado encerrado")


if __name__ == "__main__":
    main()