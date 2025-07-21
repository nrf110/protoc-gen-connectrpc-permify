.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: update
update:
	go mod tidy

.PHONY: test
test: clean update
	go test -v ./...

build: clean update
	mkdir -p ./bin
	go build -o ./bin/protoc-gen-connectrpc-permify main.go