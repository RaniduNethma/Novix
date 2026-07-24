.PHONY: up down build logs ps clean

# Start all services
up:
	docker compose up -d

# Start with fresh build
build:
	docker compose up -d --build

# Stop all services
down:
	docker compose down

# Stop and remove volumes
clean:
	docker compose down -v

# View logs for all services
logs:
	docker compose logs -f

# View logs for specific service
log-%:
	docker compose logs -f $*

# View running containers
ps:
	docker compose ps

# Restart specific service
restart-%:
	docker compose restart $*

# Build specific service only
build-%:
	docker compose build $*

# Reload Kong config without restart
kong-reload:
	docker-compose exec kong kong reload

# View Kong logs
kong-logs:
	docker-compose logs -f kong

# Test Kong is up
kong-status:
	curl -s http://localhost:8001/status | python3 -m json.tool