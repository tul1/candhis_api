COMPOSE_FILE=docker-compose.yml
# Define environment variables for database connection
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=user
DATABASE_PASSWORD=password
DATABASE_NAME=candhis_db
# Combine all environment variables into a single variable
DB_ENV_VARS=DATABASE_HOST=$(DATABASE_HOST) \
            DATABASE_PORT=$(DATABASE_PORT) \
            DATABASE_USER=$(DATABASE_USER) \
            DATABASE_PASSWORD=$(DATABASE_PASSWORD) \
            DATABASE_NAME=$(DATABASE_NAME)


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

.PHONY: logs_stack
logs_stack:
	@echo "Starting the Elasticsearch, fluentd and kibana for logs service..."
	docker-compose -f $(COMPOSE_FILE) up -d elasticsearch_logs fluentd kibana_logs

.PHONY: run-infra
run-infra: db migrate elasticsearch elasticsearch_logs
	@echo "Infrastructure services are up and running."

.PHONY: sessionID_scrapper
sessionid_scrapper:
	@echo "Starting the sessionid_scrapper service..."
	docker-compose -f $(COMPOSE_FILE) up --build sessionid_scrapper

.PHONY: campaigns_scrapper
campaigns_scrapper:
	@echo "Starting the campaigns_scrapper service..."
	docker-compose -f $(COMPOSE_FILE) up --build campaigns_scrapper

.PHONY: build-campaigns-scrapper
build-campaigns-scrapper:
	@echo "Building the campaigns_scrapper binary"
	@cd cmd/campaigns_scrapper && $(MAKE) build --no-print-directory

.PHONY: build-sessionid-scrapper
build-sessionid-scrapper:
	@echo "Building the sessionid_scrapper binary"
	@cd cmd/sessionid_scrapper && $(MAKE) build --no-print-directory

.PHONY: download
download:
	go mod download

.PHONY: build
build: build-sessionid-scrapper build-campaigns-scrapper

.PHONY: test-integration
test-integration:
	go clean -testcache
	$(DB_ENV_VARS) go test -timeout=15s -count=1 -p 1 ./test/integration/...

.PHONY: stop
stop:
	@echo "Stopping all services..."
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: clean
clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose -f $(COMPOSE_FILE) down -v
