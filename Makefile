BIN_FILE:=server

.PHONY: clean
clean:
	-rm -f ./dist/$(BIN_FILE)-macos-amd64
	-rm -f ./dist/$(BIN_FILE)-linux-amd64
	-rm -f ./dist/$(BIN_FILE)-win-amd64.exe

.PHONY: mock
mock:
	mockery --all

.PHONY: wire
wire:
	wire ./...

$(BIN_FILE): clean wire
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -o ./dist/$(BIN_FILE)-win-amd64.exe ./cmd/server

	go env -w GOOS=darwin
	go env -w GOARCH=amd64
	go build -o ./dist/$(BIN_FILE)-macos-amd64 ./cmd/server

	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o ./dist/$(BIN_FILE)-linux-amd64 ./cmd/server

.PHONY: prod
prod: clean wire
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o ./dist/$(BIN_FILE)-linux-amd64 -ldflags "-s -w" ./cmd/server
	scp ./dist/$(BIN_FILE)-linux-amd64 ai:~/proj3


.PHONY: run
run: $(BIN_FILE)
	./dist/$(BIN_FILE) -f configs/server.yml

.PHONY: dev
dev: clean wire
	CompileDaemon -build="go build -o ./dist/$(BIN_FILE) ./cmd/server" -command="./dist/$(BIN_FILE) -f configs/server.yml"


#.PHONY: deploy
#deploy:
#	docker-compose -f deployments/docker-compose.yml up --build
