.PHONY: clean
clean:
	rm -rf output

.PHONY: update
update:
	buf dep update
	go mod tidy

.PHONY: gen
gen: clean update
	buf generate
	go mod tidy