# Candhis Wave Data Scraper & API

Scrapes buoy wave data from [Candhis](https://candhis.cerema.fr/) (Cerema) and makes it usable locally. Candhis publishes the tables on the web but has no public API.

Two scrapers do the work:

- **`sessionid_scraper`** — obtains the Candhis session cookie (via headless Chrome / chromedp) and stores it in **PostgreSQL**
- **`campaigns_scraper`** — uses that session to fetch the campaign HTML table, validates each row, and indexes the observations in **Elasticsearch**

There is also a small Go HTTP API (OpenAPI). Today it only exposes `/ping`; listing wave data from Elasticsearch is still on the roadmap.

Example source page: [Les Pierres Noires](https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==) (currently the only campaign wired in).

## Storage

| Store | What lives there |
| --- | --- |
| PostgreSQL | Candhis session ID (`candhis_session`) |
| Elasticsearch | Wave observations (e.g. index `les-pierres-noires`) |

Wave rows are **not** written to Postgres.

## Prerequisites

- Docker / Docker Compose
- Go 1.23+ (for local builds and tests)

## Local setup

Start Postgres (with migrations), Elasticsearch, headless Chrome, and the optional logs stack:

```bash
make run-infra
```

Run scrapers / API with a config file (examples under `conf/`):

```bash
go run ./cmd/sessionid_scraper -config conf/sessionid_scrapper.yml
go run ./cmd/campaigns_scraper -config conf/campaigns_scrapper.yml
go run ./cmd/api -config conf/api.yml
```

`make build` produces Linux binaries under `bin/` (used for deploy).

Useful make targets: `test-unit`, `test-integration`, `test-e2e`, `lint`, `stop`, `clean`.

### Migrations

Schema changes live in `infra/db/migrations` and are applied by the `migrate` Compose service (`make run-infra` / `make db`). To add a new migration with the [golang-migrate](https://github.com/golang-migrate/migrate) CLI:

```bash
migrate create -ext sql -dir infra/db/migrations -seq <migration_name>
```

## Deploy notes

Production deploy is handled with Ansible under `infra/ansible` (systemd units/timers for the scrapers). Keep host inventory and SSH details out of this README — see that folder if you need to deploy.

## Next steps

- Expose wave data over the API (read from Elasticsearch)
- Retry when scraping fails
- Support more campaigns than Les Pierres Noires
