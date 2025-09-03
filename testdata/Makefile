.PHONY: clean
clean:
	rm -f gen.log
	rm -rf output

.PHONY: update
update:
	buf dep update
	go mod tidy

.PHONY: gen
gen: clean update
	buf generate
	go mod tidy

.PHONY: golden
golden:
	rm -rf golden
	cp -R output golden