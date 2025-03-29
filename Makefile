FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = dist
BACKEND_BUILD_CMD = go build -o ../$(DIST_DIR)

.PHONY: list run run-fe run-be deps-be deps-fe run-tool linter linter-fix test test-fe test-be build docker-build docker-run

list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

run:
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001 &
	docker run -d -p 8025:8025 -p 1025:1025 mailhog/mailhog
	cd $(FRONTEND_DIR) && yarn start &

stop:
	lsof -t -i:3000 | xargs kill -9
	lsof -t -i:3001 | xargs kill -9
	docker ps -q --filter "ancestor=mailhog/mailhog" | xargs docker kill

run-be:
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001

run-tool:
	cd $(BACKEND_DIR) && go run ./tools/main.go --seed --purge

run-fe:
	cd $(FRONTEND_DIR) && yarn dev

run-mailhog:
	docker run -d -p 8025:8025 -p 1025:1025 --platform "linux/amd64" mailhog/mailhog

deps-be:
	cd $(BACKEND_DIR) && go get -u -t ./...
	cd $(BACKEND_DIR) && go mod tidy

deps-fe:
	cd $(FRONTEND_DIR) && yarn up
	cd $(FRONTEND_DIR) && yarn upgrade-interactive

deps-install:
	cd $(FRONTEND_DIR) && yarn install
	cd $(BACKEND_DIR) && go mod download

linter:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(FRONTEND_DIR) && yarn run prettier .. --check
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && yarn run eslint

linter-fix:
	mkdir -p $(BACKEND_DIR)/cmd/assets
	touch $(BACKEND_DIR)/cmd/assets/index.html
	cd $(FRONTEND_DIR) && yarn run prettier .. --write
	cd $(FRONTEND_DIR) && yarn run eslint --fix
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && golangci-lint run --fix

test:
	$(MAKE) test-fe
	$(MAKE) test-be

test-fe:
	cd $(FRONTEND_DIR) && yarn vitest run --coverage

test-be:
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
	cd $(FRONTEND_DIR) && yarn build --outDir ../$(BACKEND_DIR)/cmd/assets -m production
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat ./cmd/main.go
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/kasseapparat-tool ./tools/main.go
	[ -f $(DIST_DIR)/.env ] || cp $(BACKEND_DIR)/.env.example $(DIST_DIR)/.env
	mkdir -p $(DIST_DIR)/data
	cd $(DIST_DIR) && ./kasseapparat-tool --seed --purge

docker-build:
	echo "$(shell date +%y%m%d%H%M)" > VERSION
	docker build -t kasseapparat:latest .
	rm VERSION

manual:
	rm -f doc/manual.pdf
#   docker buildx create --use
#	docker buildx build --platform linux/amd64 -t md-to-pdf-converter tools/md-to-pdf --load
	docker run --platform linux/amd64 --rm -v $(PWD)/doc:/app md-to-pdf-converter manual.md
	mv doc/manual.pdf $(FRONTEND_DIR)/public/manual.pdf

docker-run:
	docker run -p 3003:8080 -e "CORS_ALLOW_ORIGINS=http://localhost:3003" -v ./backend/data:/app/data kasseapparat:latest