BINARY=server
MAIN=cmd/server/main.go

build:
	go build -o bin/$(BINARY) $(MAIN)

run: build
	./bin/$(BINARY)
run:
	set DATABASE_URL=postgres://loguser:secret@localhost:5432/logdb?sslmode=disable&& go run cmd/server/main.go
test:
	go test ./... -v

lint:
	go vet ./...

proto:
	protoc --go_out=paths=source_relative:. \
	       --go-grpc_out=paths=source_relative:. \
	       pkg/proto/logs.proto