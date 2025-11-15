#!/usr/bin/env python3
"""
Broker - Proxy para Publisher-Subscriber
Parte 2: Implementa XSUB/XPUB para distribuição de mensagens
"""

import zmq
import logging
import sys

# Configuração de logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger(__name__)


def main():
    """
    Broker que conecta publishers (servidor) com subscribers (clientes)
    XSUB: recebe mensagens dos publishers (porta 5557)
    XPUB: distribui mensagens para subscribers (porta 5558)
    """
    logger.info("🚀 Iniciando Broker Pub/Sub...")
    
    context = zmq.Context()
    
    # Socket XSUB: recebe de publishers (servidor)
    xsub = context.socket(zmq.XSUB)
    xsub.bind("tcp://*:5557")
    logger.info("📥 XSUB vinculado na porta 5557 (recebe de publishers)")
    
    # Socket XPUB: envia para subscribers (clientes)
    xpub = context.socket(zmq.XPUB)
    xpub.bind("tcp://*:5558")
    logger.info("📤 XPUB vinculado na porta 5558 (envia para subscribers)")
    
    logger.info("✅ Broker pronto para rotear mensagens!")
    logger.info("=" * 60)
    
    try:
        # Proxy: conecta XSUB <-> XPUB
        # Todas as mensagens recebidas no XSUB são enviadas ao XPUB
        # Todas as inscrições recebidas no XPUB são enviadas ao XSUB
        zmq.proxy(xsub, xpub)
    except KeyboardInterrupt:
        logger.info("\n👋 Recebido sinal de interrupção. Encerrando...")
    except Exception as e:
        logger.error(f"❌ Erro no broker: {e}")
    finally:
        logger.info("🔌 Fechando sockets...")
        xsub.close()
        xpub.close()
        context.term()
        logger.info("✅ Broker encerrado com sucesso")


if __name__ == "__main__":
    main()