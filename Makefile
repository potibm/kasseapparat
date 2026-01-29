FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = dist
BACKEND_BUILD_CMD = go build -o ../$(DIST_DIR)
NODE_MAJOR := 24
GO_VERSION := 1.25

.PHONY: list run run-fe run-be deps-be deps-fe run-tool linter linter-fix test test-fe test-be build docker-build docker-run manual

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

run: check-go check-node
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001 &
	docker run -d -p 8025:8025 -p 2025:1025 --platform "linux/amd64" mailhog/mailhog
	cd $(FRONTEND_DIR) && yarn start &

stop:
	lsof -t -i:3000 | xargs kill -9
	lsof -t -i:3001 | xargs kill -9
	docker ps -q --filter "ancestor=mailhog/mailhog" | xargs docker kill

run-be: check-go
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001

run-tool: check-go
	cd $(BACKEND_DIR) && go run ./tools/main.go --seed --purge

run-fe: check-node
	cd $(FRONTEND_DIR) && corepack yarn dev

run-mailhog:
	docker run -d -p 8025:8025 -p 2025:1025 --platform "linux/amd64" mailhog/mailhog

deps-be: check-go
	cd $(BACKEND_DIR) && go get -u -t ./...
	cd $(BACKEND_DIR) && go mod tidy
	cd $(BACKEND_DIR) && gomajor list
	
deps-fe: check-node
	cd $(FRONTEND_DIR) && corepack yarn up
	cd $(FRONTEND_DIR) && corepack yarn upgrade-interactive

deps-install:
	cd $(FRONTEND_DIR) && corepack yarn install
	cd $(BACKEND_DIR) && go mod download

linter:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(FRONTEND_DIR) && corepack yarn run prettier .. --check
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && corepack yarn run eslint

linter-fix:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(FRONTEND_DIR) && corepack yarn run prettier .. --write
	cd $(FRONTEND_DIR) && corepack yarn run eslint --fix
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && golangci-lint run --fix

test:
	$(MAKE) test-fe
	$(MAKE) test-be

test-fe: check-node
	cd $(FRONTEND_DIR) && corepack yarn vitest run --coverage

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
	echo "$(shell date +%y%m%d%H%M)" > $(DIST_DIR)/VERSION
	cd $(FRONTEND_DIR) && corepack yarn build --outDir ../$(BACKEND_DIR)/cmd/assets -m production
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat ./cmd/main.go
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat-tool ./tools/main.go
	[ -f $(DIST_DIR)/.env ] || cp $(BACKEND_DIR)/.env.example $(DIST_DIR)/.env
	mkdir -p $(DIST_DIR)/data
	cd $(DIST_DIR) && ./kasseapparat-tool --seed --purge

docker-build: prepare-buildx
	@VERSION=0.0.$$(date +%y%m%d%H%M); \
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