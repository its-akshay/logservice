# ---------- Stage 1: builder ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# install git (needed for go modules sometimes)
RUN apk add --no-cache git

# copy go mod first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /server ./cmd/server


# ---------- Stage 2: runtime ----------
FROM gcr.io/distroless/static

# copy only binary
COPY --from=builder /server /server

# expose ports
EXPOSE 8080
EXPOSE 50051

# run
ENTRYPOINT ["/server"]