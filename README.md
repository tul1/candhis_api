[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=coverage)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=bugs)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=tul1_candhis_api&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=tul1_candhis_api)
# Candhis Wave Data Scraper & API

This project is a Golang-based service that scrapes wave data from the [Candhis website](https://candhis.cerema.fr/), specifically from tables containing buoy data. The data is then stored in a PostgreSQL database and pushed to an Elasticsearch instance for further analysis and retrieval. The project consists of two main components: `sessionID_scrapper` and `campaigns_scrapper`.

## Project Overview

Candhis (Centre d'Archivage National de Données de Houle In-Situ) provides public access to wave data from buoys around the coast of France. However, no public API is currently available for direct data access. This project aims to fill that gap by:

- Scraping the wave data table available on the [Candhis Campaigns page](https://candhis.cerema.fr/_public_/campagne.php).
- Storing the data in a PostgreSQL database for persistence.
- Pushing the data to an Elasticsearch instance for indexing and retrieval.

### Example Data

The URL [Candhis buoy data for Les Pierres Noires](https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==) is an example of the data provided by Candhis, showing wave data collected by the buoy located at Les Pierres Noires. The scraper extracts this data and stores it in a structured format for further use.

## Features

- **Data Scraper (`campaigns_scrapper`)**: Extracts wave data (including buoys, timestamps, wave heights, etc.) from the Candhis web page, validates the data, and stores it in a structured format.
- **Session Management (`sessionID_scrapper`)**: Manages the session ID required to access the Candhis data by periodically updating the session ID stored in the database.
- **Data Storage**: The scraped data is stored in a PostgreSQL database.
- **Data Indexing**: The data is indexed in Elasticsearch, allowing for efficient search and retrieval.
- **Automation**: The scrapers can be scheduled to run periodically to fetch and update data as needed.

## Prerequisites

Before setting up the project, ensure you have the following installed:

- Docker
- Docker Compose
- Golang (for local development)

## Setup and Running the Project

### 1. Configure Environment Variables

Ensure that the necessary environment variables for database and Elasticsearch connections are configured. These can be set in the `docker-compose.yml` file under each service.

### 2. Build and Run the Infrastructure

Use the `Makefile` to set up and run the necessary infrastructure components (PostgreSQL, migrations, Elasticsearch):

```bash
make run-infra
```

This command will:

- Start the PostgreSQL database.
- Apply the necessary database migrations.
- Start the Elasticsearch service.

### 3. Run the Scrapers

After the infrastructure is up and running, you can run the scrapers:

- **Run the `sessionID_scrapper`:**

  ```bash
  make sessionID_scrapper
  ```

  This will update the session ID in the database.

- **Run the `campaigns_scrapper`:**

  ```bash
  make campaigns_scrapper
  ```

  This will scrape the wave data from the Candhis website, store it in the database, and index it in Elasticsearch.

## Database Versioning with `golang-migrate`

This project uses `golang-migrate` for managing database schema migrations.

### How to Add a New Migration

#### **Install `golang-migrate` CLI:**

First, ensure you have `golang-migrate` installed:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Alternatively, you can use Docker to run migrations without installing the CLI locally.

#### **Create a New Migration File:**

```bash
migrate create -ext sql -dir infra/db/migrations -seq <migration_name>
```

This command will generate a new migration file in the `infra/db/migrations` directory.

#### **Run the Migration:**

```bash
docker-compose run migrate up
```

This command applies all pending migrations.

#### **Rollback a Migration (Optional):**

To roll back the last migration:

```bash
docker-compose run migrate down 1
```

### 5. Stopping and Cleaning Up

- **Stop All Services:**

  ```bash
  make stop
  ```

- **Clean Up Containers and Volumes:**

  ```bash
  make clean
  ```

## Accessing the Host and Managing Cron Jobs

### Connecting to the Host

To connect to your Vultr server, use the following SSH command:

```bash
ssh root@95.179.209.34
```

If you created a non-root user, use that username instead:

```bash
ssh astraydev@95.179.209.34
```

### Checking Cron Jobs

Once connected to the server, you can check the cron jobs that are set up for the `campaigns_scrapper` and `sessionID_scrapper` by editing or viewing the crontab:

1. **View the Crontab**:
   ```bash
   crontab -l
   ```

   This command lists all cron jobs currently set up for your user.

2. **Edit the Crontab**:
   ```bash
   crontab -e
   ```

   This command opens the crontab file in your default editor, allowing you to add, remove, or modify cron jobs.

### Managing Docker Logs

To monitor or troubleshoot the Docker containers, you can view the logs using the following commands:

1. **View All Running Containers**:
   ```bash
   docker ps
   ```

   This command lists all running Docker containers.

2. **View Logs for a Specific Container**:

   - **`campaigns_scrapper` Logs**:
     ```bash
     docker logs -f campaigns_scrapper
     ```

   - **`sessionID_scrapper` Logs**:
     ```bash
     docker logs -f sessionID_scrapper
     ```

   The `-f` option will "follow" the logs, meaning it will show real-time log updates. Press `Ctrl + C` to exit the log view.

## Next Steps

- **API Interface**: Add an API interface to allow external access to the data stored in Elasticsearch.
- **Support Multiple Campaigns**: Extend the `campaigns_scrapper` to support scraping data from multiple campaigns beyond just Les Pierres Noires.
- **Job Scheduling**: Implement an infrastructure to run the scrapers as cron jobs, with the `sessionID_scrapper` running every 12 hours and the `campaigns_scrapper` running every 30 minutes.
