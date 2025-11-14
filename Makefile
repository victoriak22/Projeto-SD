.PHONY: help build up down restart logs clean client test

help: ## Mostrar ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Construir containers
	docker-compose build

up: ## Iniciar sistema
	docker-compose up -d
	@echo "Sistema iniciado!"
	@echo "Para acessar o cliente: make client"

down: ## Parar sistema
	docker-compose down

restart: down up ## Reiniciar sistema

logs: ## Ver logs de todos os serviços
	docker-compose logs -f

logs-server: ## Ver logs do servidor
	docker-compose logs -f server

logs-client: ## Ver logs do cliente
	docker-compose logs -f client

client: ## Conectar ao cliente
	docker-compose exec client npm start

client-new: ## Iniciar novo cliente
	docker-compose run --rm client npm start

logs-broker: ## Ver logs do broker
	docker-compose logs -f broker

logs-auto: ## Ver logs dos clientes automatizados
	docker-compose logs -f auto-client-1 auto-client-2

scale-auto: ## Aumentar clientes automatizados (use: make scale-auto N=5)
	docker-compose up -d --scale auto-client-1=$(N)

clean: ## Limpar tudo (incluindo volumes)
	docker-compose down -v
	@echo "Todos os dados foram removidos!"

data: ## Mostrar dados persistidos
	docker exec -it messaging-server cat /data/server_data.json 2>/dev/null || echo "Nenhum dado encontrado"

test: ## Testar sistema completo
	@echo "Iniciando teste..."
	docker-compose up -d
	@sleep 3
	@echo "Sistema rodando. Execute 'make client' em múltiplas janelas para testar"

ps: ## Mostrar status dos containers
	docker-compose ps

rebuild: clean build up ## Reconstruir do zero