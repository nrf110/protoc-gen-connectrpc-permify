clean:
	rm -rf ./bin
sync:
	go mod tidy
build: clean sync
	mkdir -p ./bin
	go build -o ./bin/protoc-gen-connectrpc-permit main.go