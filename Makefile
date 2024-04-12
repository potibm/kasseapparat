# Makefile

.PHONY: run

run:
	go run ./backend/cmd/main.go 3001: &
	cd frontend && yarn start &

run-be:
	 go run ./backend/cmd/main.go 3001: 


run-fe:
	 cd frontend && yarn start

build:
	go build -o ./build/diekassa ./backend/cmd/main.go
	cd frontend && BUILD_PATH=../build/public yarn build
