.PHONY: build run test lint clean

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
