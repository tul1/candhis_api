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

.PHONY: elasticsearch
elasticsearch:
	@echo "Starting the Elasticsearch service..."
	docker-compose up -d elasticsearch

.PHONY: logs_stack
logs_stack:
	@echo "Starting the Elasticsearch, fluentd, metricbeat and kibana for logs service..."
	docker-compose up -d elasticsearch_logs fluentd metricbeat kibana_logs

.PHONY: run_app_infra
run-infra: db elasticsearch logs_stack
	@echo "Infrastructure services are up and running."

# Building apps #

.PHONY: build-campaigns-scrapper
build-campaigns-scrapper:
	@echo "Building the campaigns_scrapper binary"
	@cd cmd/campaigns_scrapper && $(MAKE) build --no-print-directory

.PHONY: build-sessionid-scrapper
build-sessionid-scrapper:
	@echo "Building the sessionid_scrapper binary"
	@cd cmd/sessionid_scrapper && $(MAKE) build --no-print-directory

.PHONY: build
build: build-sessionid-scrapper build-campaigns-scrapper

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
	go test -count=1 ./cmd/... ./internal/... -coverprofile cover.out

.PHONY: test-integration
test-integration:
	go clean -testcache
	$(DB_ENV_VARS) go test -timeout=15s -count=1 -p 1 ./test/integration/...

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
