# Define environment variables for database connection
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=user
DATABASE_PASSWORD=password
DATABASE_NAME=candhis_db
DB_ENV_VARS=DATABASE_HOST=$(DATABASE_HOST) \
            DATABASE_PORT=$(DATABASE_PORT) \
            DATABASE_USER=$(DATABASE_USER) \
            DATABASE_PASSWORD=$(DATABASE_PASSWORD) \
            DATABASE_NAME=$(DATABASE_NAME)

# Define environment variables for Elasticsearch connection
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_URL=ELASTICSEARCH_URL=http://$(ELASTICSEARCH_HOST):$(ELASTICSEARCH_PORT)


# Define environment variables for Chrome Headless Scraper
CHROME_URL=0.0.0.0:9222
TARGET_WEB=https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==
SCRAPER_ENV_VARS=CHROME_URL=$(CHROME_URL) \
                 TARGET_WEB=$(TARGET_WEB)


# Downloding dependencies #

.PHONY: download
download:
	go mod download

.PHONY: deps_test
deps_test:
	go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: deps_lint
deps_lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

# Infrastructure components #

.PHONY: db
db:
	@echo "Starting the database and applying migrations..."
	docker-compose up -d postgres migrate

.PHONY: chrome-headless
chrome-headless:
	@echo "Starting the chrome-headless to use scraper..."
	docker run -d --rm --cap-add=SYS_ADMIN -p 9222:9222 justinribeiro/chrome-headless:latest

.PHONY: elasticsearch
elasticsearch:
	@echo "Starting the Elasticsearch service..."
	docker-compose up -d elasticsearch

.PHONY: logs_stack
logs_stack:
	@echo "Starting the Elasticsearch, fluentd, metricbeat and kibana for logs service..."
	docker-compose up -d elasticsearch_logs fluentd metricbeat kibana_logs

.PHONY: run_app_infra
run-infra: db elasticsearch chrome-headless logs_stack
	@echo "Infrastructure services are up and running."

# Building apps #

.PHONY: build-campaigns-scraper
build-campaigns-scraper:
	@echo "Building the campaigns_scraper binary"
	@cd cmd/campaigns_scraper && $(MAKE) build --no-print-directory

.PHONY: build-sessionid-scraper
build-sessionid-scraper:
	@echo "Building the sessionid_scraper binary"
	@cd cmd/sessionid_scraper && $(MAKE) build --no-print-directory

.PHONY: build
build: build-sessionid-scraper build-campaigns-scraper

# Testing #

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run --timeout 15m0s --config .golangci.yml

.PHONY: test-unit
test-unit:
	go clean -testcache
	go test -count=1 ./internal/... -coverprofile cover.out

.PHONY: test-integration
test-integration:
	go clean -testcache
	$(DB_ENV_VARS) $(ELASTICSEARCH_URL) $(SCRAPER_ENV_VARS) go test -timeout=15s -count=1 -p 1 ./test/integration/...

# Cleaning #

.PHONY: format
format:
	golangci-lint run --config .golangci.yml --fix

.PHONY: stop
stop:
	@echo "Stopping all services..."
	docker-compose down

.PHONY: clean
clean:
	rm -rf ./bin
	rm -f ./cover.out
	find . -type d -name '*mock' -exec rm -rf {} +
