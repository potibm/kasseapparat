FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = dist
BACKEND_BUILD_CMD = go build -o ../$(DIST_DIR)

.PHONY: run run-fe run-be run-tool linter linter-fix test test-fe test-be build

run:
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001 &
	cd $(FRONTEND_DIR) && yarn start &

stop:
	lsof -t -i:3000 | xargs kill -9
	lsof -t -i:3001 | xargs kill -9

run-be:
	cd $(BACKEND_DIR) && go run ./cmd/main.go 3001

run-tool:
	cd $(BACKEND_DIR) && go run ./tools/main.go --seed --purge

run-fe:
	cd $(FRONTEND_DIR) && yarn start

linter:
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && yarn run eslint src/

linter-fix:
	cd $(BACKEND_DIR) && golangci-lint run --fix
	cd $(FRONTEND_DIR) && yarn run eslint src/ --fix

test:
	$(MAKE) test-fe
	$(MAKE) test-be

test-fe:
	cd $(FRONTEND_DIR) && yarn test --coverage --watchAll=false

test-be:
	cd $(BACKEND_DIR) && go test ./...

build:
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/diekassa ./cmd/main.go
	cd $(BACKEND_DIR) && $(BACKEND_BUILD_CMD)/diekassa-tool ./tools/main.go
	cd $(FRONTEND_DIR) && BUILD_PATH=../$(DIST_DIR)/public yarn build
	mkdir -p $(DIST_DIR)/data
	cd $(DIST_DIR) && ./diekassa-tool --seed --purge