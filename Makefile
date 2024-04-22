# Makefile

.PHONY: run

run:
	cd backend && go run ./cmd/main.go 3001 &
	cd frontend && yarn start &

stop:
	lsof -t -i:3000 | xargs kill -9
	lsof -t -i:3001 | xargs kill -9

run-be:
	cd backend && go run ./cmd/main.go 3001

run-tool:
	cd backend && go run ./tools/main.go --seed --purge

run-fe:
	cd frontend && yarn start

linter:
	cd backend && golangci-lint run
	cd frontend && yarn run eslint src/

linter-fix:
	cd backend && golangci-lint run --fix
	cd frontend && yarn run eslint src/ --fix

test:
	cd frontend && yarn test --coverage --watchAll=false

build:
	cd backend && go build -o ../dist/diekassa ./cmd/main.go
	cd backend && go build -o ../dist/diekassa-tool ./tools/main.go
	cd frontend && BUILD_PATH=../dist/public yarn build
	mkdir -p dist/data
	cd dist && ./diekassa-tool --seed --purge
