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