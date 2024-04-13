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

run-fe:
	 cd frontend && yarn start

build:
	cd backend && go build -o ../dist/diekassa ./cmd/main.go
	cd frontend && BUILD_PATH=../dist/public yarn build
	mkdir -p dist/data
	# @TODO create a script to initialize the sqlite database
