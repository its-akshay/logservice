BINARY=server
MAIN=cmd/server/main.go

# ===== ENV VARIABLES =====
DB_URL=postgres://loguser:secret@localhost:5432/logdb?sslmode=disable
JWT_SECRET=super-secret-access
JWT_REFRESH_SECRET=super-secret-refresh

# ===== BUILD =====
build:
	go build -o bin/$(BINARY) $(MAIN)

# ===== RUN (with envs) =====
run:
	set "DATABASE_URL=$(DB_URL)" && \
	set "JWT_SECRET=$(JWT_SECRET)" && \
	set "JWT_REFRESH_SECRET=$(JWT_REFRESH_SECRET)" && \
	go run $(MAIN)

# ===== RUN BINARY =====
run-bin: build
	set DATABASE_URL=$(DB_URL) && \
	set JWT_SECRET=$(JWT_SECRET) && \
	set JWT_REFRESH_SECRET=$(JWT_REFRESH_SECRET) && \
	./bin/$(BINARY)

# ===== TEST =====
test:
	go test ./... -v

# ===== LINT =====
lint:
	go vet ./...

# ===== PROTO =====
proto:
	protoc --go_out=paths=source_relative:. \
	       --go-grpc_out=paths=source_relative:. \
	       pkg/proto/logs.proto