.PHONY: build run test lint clean ui-install ui-dev ui-build ui-test ui-lint

# Go
build:
	go build -o bin/stream-intel ./cmd/stream-intel

run:
	go run ./cmd/stream-intel

test:
	go test ./... -v

lint:
	go vet ./...
	gofmt -l .

clean:
	rm -rf bin/

# UI
ui-install:
	cd ui && pnpm install

ui-dev:
	cd ui && pnpm dev

ui-build:
	cd ui && pnpm build

ui-test:
	cd ui && pnpm test

ui-lint:
	cd ui && pnpm lint
