BIN_FILE:=server.exe

.PHONY: clean
clean:
	-rm -f ./dist/$(BIN_FILE)

.PHONY: mock
mock:
	mockery --all

.PHONY: wire
wire:
	wire ./...

$(BIN_FILE): clean wire
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -o ./dist/$(BIN_FILE)  ./cmd/server

.PHONY: run
run: $(BIN_FILE)
	./dist/$(BIN_FILE) -f configs/server.yml

.PHONY: dev
dev: clean wire
	CompileDaemon -build="go build -o ./dist/$(BIN_FILE) ./cmd/server" -command="./dist/$(BIN_FILE) -f configs/server.yml"


#.PHONY: deploy
#deploy:
#	docker-compose -f deployments/docker-compose.yml up --build
