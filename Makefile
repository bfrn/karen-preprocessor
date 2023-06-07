build:
	@go build -o bin/preprocessor cmd/server/main.go

run: build
	./bin/preprocessor

test:
	go test -v ./...
