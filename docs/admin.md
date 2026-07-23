# Kasseapparat: Admin Documentation

This documentation will give hints how to set up Kasseapparat on a server.

## Prerequisites

### Server

- Server with minimal specs (the staging environment is running smoothly on a VPS with 2 cores and 2 GB ram)
- Current Ubuntu (as we are using docker a different setup will probably work just as fine.)
- Docker including docker compose installed
- A mail account (will be used to send visitor arrival notifications)

### Client

Find a list below with hardware that was tried and tested at demoparties (feel free to extend this list). Other hardware will work as well.

- Apple iPad Air 10.9" Wi-Fi 64GB, 4th generation (add as [shortcut to the home screen](https://support.apple.com/en-gb/guide/shortcuts/apd735880972/ios))
- Barcode scanner:
  - [Tera 2500C](https://tera-digital.com/products/2500c-ccd-1d-usb-barcode-scanner-wholesale) (preferred to scan mobile displays)
  - [Tera 5100](https://tera-digital.com/products/5100-laser-1d-wireless-barcode-scanner-wholesale)

## Set up directories

- create directory /app/kasseapparat
- create directory /app/kasseapparat/data
- create directory /app/kasseapparat/backup
- create directory /app/kasseapparat/config

## Configuration

Kasseapparat uses a layered configuration approach. The primary source of truth is the `config.yaml` file, while the `.env` file is strictly reserved for sensitive secrets (like passwords).

### 1. Main Configuration (`config.yaml`)

Run the following command to generate the base configuration:

```bash
docker compose run --rm kasseapparat config create
```

This generates `config/config.yaml` with sensible defaults. You should edit this file to adjust the core behavior and business rules of the application. Key areas to configure include:

- **Authentication Mode:** Kasseapparat relies entirely on external identity providers. Set the mode to `proxy` (default, e.g., Traefik with ForwardAuth) or `oidc` (e.g., Dex, Keycloak). In proxy mode, login/logout UI elements are disabled, and the backend trusts the `Remote-User` header.
- **System URLs:** Define your frontend URL to ensure correct link generation within the application and notifications.
- **Localization:** Adjust locale, currency codes, and fraction digits to your local preferences.
- **Mail Settings:** Configure your sender address and subject prefixes for automated visitor arrival notifications.
- **Business Rules:** Define VAT rates and accepted payment methods.

_Optional:_ You can create a `config/config.local.yaml` for environment-specific settings that shouldn't be committed to version control.

### 2. Secrets & Credentials (`.env`)

While the YAML file handles the application structure, sensitive credentials must be kept in the `.env` file.

Copy the `.env.example` from the repository to your server and rename it to `.env`. Here you only need to fill in:

- Your SMTP login credentials (e.g., `MAIL_DSN`)
- Any required API keys

_(Note: While it is technically possible to override YAML settings via environment variables like `APP_AUTH_MODE`, sticking to the `config.yaml` is highly recommended for a clean setup)._

## Command Line Interface (CLI)

Kasseapparat comes with a powerful built-in CLI to manage the application.

### Global Flags

These flags can be appended to almost any command:

- `--log-level=info` (or debug, warn, error)
- `--log-format=json` (or text)
- `--db-file="kasseapparat"` (to set the SQLite filename)

### Available Commands

- `kasseapparat serve`: Starts the main web server and API.
- `kasseapparat config create`: Generates `config/config.yaml` with default values (use `--force` to overwrite).
- `kasseapparat config export`: Prints the final, merged configuration (config.yaml + config.local.yaml + .env + CLI flags) as a JSON tree. Sensitive data like secrets and API keys are automatically redacted for safety.
- `kasseapparat database migrate`: Creates or updates the database tables to the latest schema.
- `kasseapparat database seed`: Fills the database with dummy data (useful for development).
- `kasseapparat database reset`: Drops all tables and recreates them from scratch (WARNING: Deletes all data!).

### SENTRY

We are using [https://sentry.io/](https://sentry.io/) for fetching some bugs. Please ignore those settings in the config file.

## Create a /app/kasseapparat/docker-compose.yml

```yaml
services:
  traefik:
    image: traefik:v3.7
    restart: always
    ports:
      - 80:80
      - 443:443
    networks:
      - proxy
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - letsencrypt:/letsencrypt
    command:
      - --api.dashboard=true
      - --log.level=INFO
      - --accesslog=true
      - --providers.docker.network=proxy
      - --providers.docker.exposedByDefault=false
      - --entrypoints.web.address=:80
      - --entrypoints.web.http.redirections.entrypoint.to=websecure
      - --entrypoints.web.http.redirections.entrypoint.scheme=https
      - --entrypoints.websecure.address=:443
      - --entrypoints.websecure.http.tls.certresolver=myresolver
      - --certificatesresolvers.myresolver.acme.email=me@example.com
      - --certificatesresolvers.myresolver.acme.tlschallenge=true
      - --certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json
    labels:
      - traefik.enable=true
      - traefik.http.routers.mydashboard.rule=Host(`kasseapparat-traefik.example.com`)
      - traefik.http.routers.mydashboard.service=api@internal
      - traefik.http.routers.mydashboard.middlewares=myauth
      - traefik.http.middlewares.myauth.basicauth.users=dashboard:somepassword

  kasseapparat:
    image: ghcr.io/potibm/kasseapparat:latest
    restart: always
    volumes:
      - ./data:/app/data
      - ./config:/app/config:ro
    env_file: ".env"
    environment:
      - "APP_GIN_MODE=release"
      - "APP_CORS_ALLOW_ORIGINS=https://kasseapparat.example.com"
      - "APP_FRONTEND_URL=https://kasseapparat.example.com"
    labels:
      - traefik.enable=true
      - traefik.http.routers.kasseapparat.entrypoints=websecure
      - traefik.http.routers.kasseapparat.rule=Host(`kasseapparat.example.com`)
      - traefik.http.routers.kasseapparat.tls.certresolver=myresolver
      - traefik.http.routers.kasseapparat.tls=true
      - traefik.http.middlewares.mywwwredirect.redirectregex.regex=^https://www\.(.*)
      - traefik.http.middlewares.mywwwredirect.redirectregex.replacement=https://$${1}
      - traefik.http.routers.kasseapparat.middlewares=mywwwredirect
      # --- Uncomment the following lines if you are using proxy auth ---
      # - traefik.http.routers.kasseapparat.middlewares=mywwwredirect,forward-auth
      # - traefik.http.middlewares.forward-auth.forwardauth.address=http://authelia:9091/api/verify?rd=https://auth.example.com/
      # - traefik.http.middlewares.forward-auth.forwardauth.trustForwardHeader=true
      # - traefik.http.middlewares.forward-auth.forwardauth.authResponseHeaders=Remote-User,Remote-Groups
      # -----------------------------------------------------------------
    networks:
      - proxy

networks:
  proxy:
    name: proxy
    external: true

volumes:
  letsencrypt:
    name: letsencrypt
```

### Urls

Replace the urls above with the one that you will use.

### Password

According to https://doc.traefik.io/traefik/middlewares/http/basicauth/ you may generate a password with

```bash
echo $(htpasswd -nB user) | sed -e s/\\$/\\$\\$/g
```

## Create update.sh

```bash
#!/bin/bash

# Backup database
BASEDIR=$(dirname $0)
FILENAME=$BASEDIR/backup/data_`date +"%Y%m%d_%H%M%S"`.tar.gz
tar cfvz $FILENAME $BASEDIR/data/

# Pull the latest image
docker pull ghcr.io/potibm/kasseapparat:latest

# Stop and remove the existing container
docker compose stop kasseapparat
docker compose rm -f kasseapparat

# Start the container with the latest image
docker compose up -d kasseapparat

# Ensure the container is started successfully
if [ $(docker ps -q -f name=kasseapparat) ]; then
  echo "Container started successfully"
else
  echo "Failed to start container"
  exit 1
fi

# Execute any necessary commands inside the container
docker compose run --rm kasseapparat database migrate

# Optional: Clean up dangling images
docker image prune -f
```

Make the script executable.

## First start

### Correct UID/GUID for data directory

The data directory at /app/kasseapparat/data needs the correct UID/GID. That is probably 1000:1000, so

```bash
sudo chown 1000:1000 /app/kasseapparat/data
```

You can check those values by running

```bash
docker run --rm ghcr.io/potibm/kasseapparat:latest id appuser
```

### Create Database

```bash
docker compose run --rm kasseapparat database seed --test-data
```

## Startup

```bash
docker compose up -d
```

## Update

To update the docker image call the update.sh. A backup is performed and stored in the backup directory.

## OpenTelemetry & Monitoring

Kasseapparat natively supports **OpenTelemetry (OTel)**. You can enable the export of Traces, Logs, and Metrics by providing the OTLP endpoint:

`--otel-endpoint="localhost:3017"` (compatible with backends like **OpenObserve**, Grafana, or Jaeger).

### Key Metrics & Dashboards

Use the following metrics in your monitoring backend (e.g., OpenObserve) to track system health:

- **Go Runtime:**
  - `go_goroutine_count`: Number of active goroutines (essential to identify connection leaks).
  - `go_memory_count`: Current memory usage/statistics.
- **WebSockets:**
  - `ws_active_connections`: Current number of connected terminals/iPads (Gauge).
  - `ws_messages_sent_total`: Data throughput (use the `msg_type` attribute for filtering).
- **Database (SQL & GORM):**
  - `go.sql.connections_open`: Current number of established database connections.
  - `go.sql.connections_in_use`: Number of connections currently processing queries.
  - `go.sql.connections_wait_duration`: Total time blocked waiting for a new connection (indicates SQLite bottlenecks).
- **Performance:**
  - `http_server_request_duration`: Measures API latency from the Gin router.

### OpenObserve

For local development we provide [OpenObserve](https://github.com/openobserve/openobserve) at [https://observe.kasseapparat.test](https://observe.kasseapparat.test) by running `mise run infra:up`.

Sample dashboard templates can be found at `infra/openobserve`.
