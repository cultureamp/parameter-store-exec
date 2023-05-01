VERSION = $(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w

parameter-store-exec: main.go paramstore/store.go
	go build -ldflags="$(FLAGS)"

.PHONY: test
test:
	go test ./...

.PHONY: mod
mod:
	go mod download

# ensures that `go mod tidy` has been run after any dependency changes
.PHONY: ensure-deps
ensure-deps: mod
	@go mod tidy
	@git diff --exit-code
