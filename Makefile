COMPOSE_FILE=./deploy/docker-compose.yml
COMPOSE=docker-compose -f $(COMPOSE_FILE)

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down -v

restart: down up

logs:
	$(COMPOSE) logs -f

build:
	$(COMPOSE) build
