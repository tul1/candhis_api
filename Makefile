COMPOSE_FILE=docker-compose.yml

.PHONY: db
db:
	@echo "Starting the database and applying migrations..."
	docker-compose -f $(COMPOSE_FILE) up -d postgres migrate

.PHONY: migrate
migrate: db
	@echo "Applying migrations..."
	docker-compose -f $(COMPOSE_FILE) up migrate

.PHONY: elasticsearch
elasticsearch:
	@echo "Starting the Elasticsearch service..."
	docker-compose -f $(COMPOSE_FILE) up -d elasticsearch

.PHONY: run-infra
run-infra: db migrate elasticsearch
	@echo "Infrastructure services are up and running."

.PHONY: sessionID_scrapper
sessionid_scrapper:
	@echo "Starting the sessionid_scrapper service..."
	docker-compose -f $(COMPOSE_FILE) up --build sessionid_scrapper

.PHONY: campaigns_scrapper
campaigns_scrapper:
	@echo "Starting the campaigns_scrapper service..."
	docker-compose -f $(COMPOSE_FILE) up --build campaigns_scrapper

.PHONY: stop
stop:
	@echo "Stopping all services..."
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: clean
clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose -f $(COMPOSE_FILE) down -v
