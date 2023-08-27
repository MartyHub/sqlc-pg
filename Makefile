default: watch

all: tidy vet lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/sqlc-pg/internal.Version=version dev'" -race

examples: build
	rm -f internal/testdata/sqlc/examples/*/sqlc/*.go
	find internal/testdata/sqlc/examples/ -name 'sqlc.yaml' | xargs -I {} -P 4 sqlc generate --file {}

lint:
	golangci-lint run

test:
	go test -race -timeout 10s ./...

tidy:
	go mod tidy

vet:
	go vet ./...

watch:
	modd --file=.modd.conf

.PHONY: all build examples lint test tidy vet watch
