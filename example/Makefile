.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf example/internal/gen

.PHONY: update
update:
	buf dep update
	go mod tidy

.PHONY: gen
gen: clean
	buf generate ./...

.PHONY: build
build: gen update
	mkdir -p ./bin
	go build -o ./bin/protoc-gen-connectrpc-permify main.go