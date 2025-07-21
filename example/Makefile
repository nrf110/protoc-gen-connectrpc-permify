.PHONY: clean
clean:
	rm -rf internal/gen

.PHONY: update
update:
	buf dep update
	go mod tidy

.PHONY: build
build: gen update
	mkdir -p ./bin
	go build -o ./bin/protoc-gen-connectrpc-permify main.go

.PHONY: gen
gen: clean update
	buf generate
	go mod tidy