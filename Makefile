FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = dist
VERSION = 0.0.$(shell date +%y%m%d%H%M)
BACKEND_BUILD_CMD = go build -ldflags "-X main.version=$(VERSION)" -o ../$(DIST_DIR)
NODE_MAJOR := 25
GO_VERSION := 1.26

.PHONY: list run run-fe run-be deps-be deps-fe deps-actions run-tool linter linter-fix test test-fe test-be build docker-build docker-run manual e2e-setup e2e-run e2e-report 

list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

prepare-buildx:
	@if ! docker buildx inspect kasseapparat-builder >/dev/null 2>&1; then \
		echo "🔧 Creating buildx builder 'kasseapparat-builder'..."; \
		docker buildx create --name kasseapparat-builder --use --bootstrap; \
	else \
		echo "✅ Using existing buildx builder 'kasseapparat-builder'"; \
		docker buildx use kasseapparat-builder; \
		docker buildx inspect --bootstrap > /dev/null; \
	fi

run: check-go check-node infra-up run-be run-fe
	@echo "🏃 Starte Application..."

stop: infra-down
	@echo "🎯 Stop local processes..."
	lsof -t -i:3000 | xargs kill -9
	lsof -t -i:3001 | xargs kill -9
	rm -rf $(BACKEND_DIR)/logs/app.json

run-be: check-go
	cd $(BACKEND_DIR) && go run ./cmd/main.go --port=3001 --log-level=debug --otel-endpoint=localhost:4317

run-tool: check-go
	cd $(BACKEND_DIR) && go run ./tools/main.go --seed --purge

run-fe: check-node
	cd $(FRONTEND_DIR) && corepack yarn dev

deps-be: check-go
	cd $(BACKEND_DIR) && go get -u -t ./...
	cd $(BACKEND_DIR) && go mod tidy
	cd $(BACKEND_DIR) && gomajor list
	
deps-fe: check-node
	cd $(FRONTEND_DIR) && corepack yarn up -R
	cd $(FRONTEND_DIR) && corepack yarn upgrade-interactive

deps-install:
	cd $(FRONTEND_DIR) && corepack yarn install
	cd $(BACKEND_DIR) && go mod download

deps-actions:
	npx actions-up	

linter:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && corepack yarn run tsc --noEmit
	cd $(FRONTEND_DIR) && corepack yarn run eslint
	cd $(FRONTEND_DIR) && corepack yarn run prettier .. --check
	cd $(BACKEND_DIR) && dotenv-linter check . --ignore-checks QuoteCharacter,ValueWithoutQuotes
	cd $(FRONTEND_DIR) && dotenv-linter check . 

linter-fix:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(FRONTEND_DIR) && corepack yarn run prettier .. --write
	cd $(FRONTEND_DIR) && corepack yarn run tsc --noEmit
	cd $(FRONTEND_DIR) && corepack yarn run eslint --fix
	cd $(FRONTEND_DIR) && dotenv-linter fix . --no-backup
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && golangci-lint run --fix 
	cd $(BACKEND_DIR) && dotenv-linter fix . --no-backup --ignore-checks QuoteCharacter,ValueWithoutQuotes

test:
	$(MAKE) test-fe
	$(MAKE) test-be

test-fe: check-node
	cd $(FRONTEND_DIR) && NODE_OPTIONS="--no-webstorage" corepack yarn vitest run --coverage

test-be: check-go
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(BACKEND_DIR) && go test -cover -coverprofile=coverage.out -coverpkg=./... -v ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

sec-be:
	cd $(BACKEND_DIR) && gosec ./...

build:
	rm -rf $(BACKEND_DIR)/cmd/assets
	mkdir -p $(BACKEND_DIR)/cmd/assets
	cd $(FRONTEND_DIR) && corepack yarn build --outDir ../$(BACKEND_DIR)/cmd/assets -m production
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat ./cmd/main.go
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat-tool ./tools/main.go
	[ -f $(DIST_DIR)/.env ] || cp $(BACKEND_DIR)/.env.example $(DIST_DIR)/.env
	mkdir -p $(DIST_DIR)/data
	cd $(DIST_DIR) && ./kasseapparat-tool --seed --purge

docker-build: prepare-buildx
	@VERSION=$(VERSION); \
	BUILD_DATE=$$(date -Iseconds); \
	docker buildx build \
	    --builder kasseapparat-builder \
		--build-arg VERSION=$$VERSION \
		--build-arg BUILD_DATE=$$BUILD_DATE \
		--tag kasseapparat:latest \
		--tag kasseapparat:$$VERSION \
		--load \
		. ; \
	
docker-run:
	@if ! docker image inspect kasseapparat:latest >/dev/null 2>&1; then \
		echo "❌ Docker image 'kasseapparat:latest' not found. Please run 'make docker-build' first."; \
		exit 1; \
	fi
	docker run -p 3003:8080 \
		-e "CORS_ALLOW_ORIGINS=http://localhost:3003" \
		-v ./backend/data:/app/data \
		kasseapparat:latest
		
manual: prepare-buildx
	@echo "🧹 Removing old manual.pdf if it exists..."
	@rm -f doc/manual.pdf

	@echo "🐳 Building Docker image (if needed)..."
	@docker buildx build \
		--builder kasseapparat-builder \
		--platform linux/amd64 \
		-t md-to-pdf-converter tools/md-to-pdf \
		--load

	@echo "📄 Generating manual.pdf from markdown..."
	@docker run --platform linux/amd64 --rm \
		-v "$(PWD)/doc:/app" \
		--shm-size=1g \
		md-to-pdf-converter manual.md

	@echo "📁 Moving generated PDF to frontend..."
	@mkdir -p "$(FRONTEND_DIR)/public"
	@mv doc/manual.pdf "$(FRONTEND_DIR)/public/manual.pdf"

check-node:
	@node -v | grep -q "^v$(NODE_MAJOR)\." || \
	(echo "❌ Node $(NODE_MAJOR).x required. Current: $$(node -v)"; exit 1)

check-go:
	@CURRENT=$$(go version | awk '{print $$3}' | sed -E 's/go([0-9]+\.[0-9]+).*/\1/') ; \
	if [ "$$CURRENT" != "$(GO_VERSION)" ]; then \
	  echo "❌ Go $(GO_VERSION).x required. Current: $$(go version)"; \
	  exit 1; \
	fi

e2e-setup:
	@echo "Create test database..."
	cd $(BACKEND_DIR) && go run ./tools/main.go --seed-with-test --purge --db-file "e2e-clean"
	@echo "Copying test database to active location..."
	cd $(BACKEND_DIR) && cp data/e2e-clean.db data/e2e-work.db

e2e-run: e2e-setup
	cd $(FRONTEND_DIR) && corepack yarn playwright test

e2e-report:
	cd $(FRONTEND_DIR) && corepack yarn playwright show-report

infra-up:
	@echo "🚀 Start up infrastructure (Mailhog, OpenObserve, Collector)..."
	docker compose up -d
	@echo "📧 Mailhog is available at http://localhost:8025"
	@echo "📊 OpenObserve is available at http://localhost:5080"

infra-down:
	@echo "🛑 Stop infrastructure..."
	docker compose down