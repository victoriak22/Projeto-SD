#!/usr/bin/env python3
"""
Servidor de Referência - Parte 4 Etapa 2
Gerencia registro, ranks e heartbeats dos servidores
"""

import zmq
import msgpack
import time
import json
import logging
import sys
import os
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

# Estrutura de dados dos servidores
class ServerRegistry:
    def __init__(self, data_file='/data/reference_data.json'):
        self.data_file = data_file
        self.servers = {}  # {name: {rank, last_heartbeat}}
        self.next_rank = 1
        self.load_data()
    
    def load_data(self):
        """Carrega dados persistidos"""
        try:
            if os.path.exists(self.data_file):
                with open(self.data_file, 'r') as f:
                    data = json.load(f)
                    self.servers = data.get('servers', {})
                    self.next_rank = data.get('next_rank', 1)
                logger.info(f"📊 Dados carregados: {len(self.servers)} servidores")
            else:
                logger.info("📊 Iniciando com registro vazio")
        except Exception as e:
            logger.error(f"❌ Erro ao carregar dados: {e}")
    
    def save_data(self):
        """Salva dados em disco"""
        try:
            os.makedirs(os.path.dirname(self.data_file), exist_ok=True)
            data = {
                'servers': self.servers,
                'next_rank': self.next_rank
            }
            with open(self.data_file, 'w') as f:
                json.dump(data, f, indent=2)
        except Exception as e:
            logger.error(f"❌ Erro ao salvar dados: {e}")
    
    def register_server(self, name):
        """Registra um novo servidor ou retorna rank existente"""
        if name in self.servers:
            # Servidor já existe, retorna rank existente
            rank = self.servers[name]['rank']
            logger.info(f"🔄 Servidor {name} já registrado com rank {rank}")
        else:
            # Novo servidor, atribui próximo rank
            rank = self.next_rank
            self.servers[name] = {
                'rank': rank,
                'last_heartbeat': time.time()
            }
            self.next_rank += 1
            self.save_data()
            logger.info(f"✅ Novo servidor registrado: {name} com rank {rank}")
        
        # Atualizar heartbeat
        self.servers[name]['last_heartbeat'] = time.time()
        return rank
    
    def update_heartbeat(self, name):
        """Atualiza o heartbeat de um servidor"""
        if name in self.servers:
            self.servers[name]['last_heartbeat'] = time.time()
            logger.debug(f"💓 Heartbeat recebido de {name}")
            return True
        else:
            logger.warning(f"⚠️  Heartbeat de servidor não registrado: {name}")
            return False
    
    def get_server_list(self):
        """Retorna lista de servidores ativos"""
        current_time = time.time()
        timeout = 30  # 30 segundos
        
        # Filtrar servidores ativos (receberam heartbeat recentemente)
        active_servers = []
        for name, info in self.servers.items():
            if current_time - info['last_heartbeat'] < timeout:
                active_servers.append({
                    'name': name,
                    'rank': info['rank']
                })
        
        return active_servers
    
    def cleanup_inactive_servers(self):
        """Remove servidores inativos (sem heartbeat há muito tempo)"""
        current_time = time.time()
        timeout = 60  # 60 segundos
        
        inactive = []
        for name, info in self.servers.items():
            if current_time - info['last_heartbeat'] > timeout:
                inactive.append(name)
        
        for name in inactive:
            logger.warning(f"🗑️  Removendo servidor inativo: {name}")
            del self.servers[name]
            self.save_data()


# Instância global do registro
registry = ServerRegistry()


def handle_rank(request):
    """Handler para serviço 'rank' - atribui rank ao servidor"""
    try:
        server_name = request['data']['user']
        
        # Registrar servidor e obter rank
        rank = registry.register_server(server_name)
        
        response = {
            'service': 'rank',
            'data': {
                'rank': rank,
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }
        
        return response
    except Exception as e:
        logger.error(f"❌ Erro em handle_rank: {e}")
        return {
            'service': 'rank',
            'data': {
                'error': str(e),
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }


def handle_list(request):
    """Handler para serviço 'list' - retorna lista de servidores"""
    try:
        server_list = registry.get_server_list()
        
        response = {
            'service': 'list',
            'data': {
                'list': server_list,
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }
        
        logger.debug(f"📋 Lista de servidores: {len(server_list)} ativos")
        
        return response
    except Exception as e:
        logger.error(f"❌ Erro em handle_list: {e}")
        return {
            'service': 'list',
            'data': {
                'error': str(e),
                'list': [],
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }


def handle_heartbeat(request):
    """Handler para serviço 'heartbeat' - atualiza status do servidor"""
    try:
        server_name = request['data']['user']
        
        # Atualizar heartbeat
        registry.update_heartbeat(server_name)
        
        response = {
            'service': 'heartbeat',
            'data': {
                'status': 'OK',
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }
        
        return response
    except Exception as e:
        logger.error(f"❌ Erro em handle_heartbeat: {e}")
        return {
            'service': 'heartbeat',
            'data': {
                'error': str(e),
                'timestamp': int(time.time()),
                'clock': increment_clock()
            }
        }


def main():
    """Função principal do servidor de referência"""
    logger.info("🚀 Iniciando Servidor de Referência...")
    
    # Configurar ZeroMQ
    context = zmq.Context()
    socket = context.socket(zmq.REP)
    
    port = 5559
    socket.bind(f"tcp://*:{port}")
    logger.info(f"📡 Servidor de Referência escutando na porta {port}")
    logger.info("=" * 60)
    
    # Contador de mensagens para limpeza periódica
    message_count = 0
    
    try:
        while True:
            # Receber requisição
            message = socket.recv()
            request = msgpack.unpackb(message, raw=False)
            
            # Atualizar relógio lógico
            if 'data' in request and 'clock' in request['data']:
                update_clock(request['data']['clock'])
            
            service = request.get('service')
            logger.debug(f"📥 Requisição recebida: {service}")
            
            # Processar requisição
            if service == 'rank':
                response = handle_rank(request)
            elif service == 'list':
                response = handle_list(request)
            elif service == 'heartbeat':
                response = handle_heartbeat(request)
            else:
                logger.warning(f"⚠️  Serviço desconhecido: {service}")
                response = {
                    'service': service,
                    'data': {
                        'error': f'Serviço desconhecido: {service}',
                        'timestamp': int(time.time()),
                        'clock': increment_clock()
                    }
                }
            
            # Enviar resposta
            response_packed = msgpack.packb(response)
            socket.send(response_packed)
            
            # Limpeza periódica de servidores inativos
            message_count += 1
            if message_count % 100 == 0:
                registry.cleanup_inactive_servers()
    
    except KeyboardInterrupt:
        logger.info("\n👋 Encerrando Servidor de Referência...")
    except Exception as e:
        logger.error(f"❌ Erro no servidor: {e}")
    finally:
        socket.close()
        context.term()
        logger.info("✅ Servidor de Referência encerrado")


if __name__ == "__main__":
    main()