# Kasseapparat: Admin Documentation

This documentation will give hints how to set up Kasseapparat on a server.

## Prequisites

- Server with minimal specs (the staging environment is running smoothly on a VPS with 2 cores and 2 GB ram)
- Current Ubuntu (as we are using docker a different setup will probably work just as fine.)
- Docker including docker compose installed
- A mail account (will be used to send password reset mails)

## Set up directories

- create directory /app/kassepparat
- create directory /app/kassepparat/data
- create directory /app/kassepparat/backup

## create /app/kasseapparat/.env

```
JWT_SECRET=1234
SENTRY_DSN=""
SENTRY_TRACE_SAMPLE_RATE="0.1"
SENTRY_REPLAY_SESSION_SAMPLE_RATE="0.1"
SENTRY_REPLAY_ERROR_SAMPLE_RATE="1"
LOCALE="da-DE"
CURRENCY_CODE="EUR"
FRACTION_DIGITS_MIN="0"
FRACTION_DIGITS_MAX="2"
FRONTEND_URL="https://kasseapparat.example.cp,"
MAIL_DSN="smtp://username:password@smtp.example.com:587"
MAIL_FROM="kasseapparat<kasseapparat@example.com>"
MAIL_SUBJECT_PREFIX="[Kasseapparat] "
```

### JWT_SECRET

Generate a random JWT secret e.g. by calling

```
< /dev/urandom tr -dc 'A-Za-z0-9!@#$%^&*()_+=' | head -c 16
```

or

```
openssl rand -base64 32
```

Using the default value is a major security risk. You should really add some entrophy here.

### SENTRY

We are using https://sentry.io/ for fetching some bugs. Please ignore those settings.

### LOCALE

You can set LOCALE, CURRENCY_CODE, FRACTION_DIGITS_MIN and FRACTION_DIGITS_MAX to your local preferences.

### FRONTEND_URL

For generating correct urls (within mails e.g.) set the URL here.

# MAIL

We will need the SMTP login for the mail account in MAIL_DSN.

Modify MAIL_FROM accordingly. Editing MAIL_SUBJECT_PREFIX is optional.

## Create a /app/kasseapparat/docker-compose.yml

```
services:
  traefik:
    image: traefik:v3.0
    restart: always
    ports:
      - 80:80
      - 443:443
    networks:
      - proxy
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
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
      - ./data:/app/kasseapparat/data
    env_file: ".env"
    environment:
      - "GIN_MODE=release"
      - "JWT_REALM=Kasseapparat"
      - "JWT_SECRET=${JWT_SECRET}"
      - "JWT_TIMEOUT=10"
      - "CORS_ALLOW_ORIGINS=https://kasseapparat.example.com"
    labels:
      - traefik.enable=true
      - traefik.http.routers.kasseapparat.entrypoints=websecure
      - traefik.http.routers.kasseapparat.rule=Host(`kasseapparat.example.com`)
      - traefik.http.routers.kasseapparat.tls.certresolver=myresolver
      - traefik.http.routers.kasseapparat.tls=true
      - traefik.http.middlewares.mywwwredirect.redirectregex.regex=^https://www\.(.*)
      - traefik.http.middlewares.mywwwredirect.redirectregex.replacement=https://$${1}
      - traefik.http.routers.kasseapparat.middlewares=mywwwredirect
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

```
echo $(htpasswd -nB user) | sed -e s/\\$/\\$\\$/g
```

## Create update.sh

```
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
docker compose exec kasseapparat /app/kasseapparat-tool

# Optional: Clean up dangling images
docker image prune -f
```

Make the script executable.

## First start

### Create Database

```
docker compose exec kasseapparat /app/kasseapparat-tool
```

### Create first user

Create /app/kasseapparat/data/user.txt with the following structure (please, edit accordingly)

```
username,email@example.com,true
```

Call

```
docker compose exec kasseapparat /app/kasseapparat-tool --import-users /data/user.txt
```

to create a user called "username" with the email "email@example.com" as an admin. You should receive an email to change your password.

## Startup

```
docker compose up -d
```

## Update

To update the docker image call the update.sh. A backup is performed and stored in the backup directory.
