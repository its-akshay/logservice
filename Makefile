BINARY=server
MAIN=cmd/server/main.go

build:
	go build -o bin/$(BINARY) $(MAIN)

run: build
	./bin/$(BINARY)

test:
	go test ./... -v

lint:
	go vet ./...