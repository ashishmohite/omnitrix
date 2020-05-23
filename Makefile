TIMEOUT  = 30

.PHONY: clean
clean:
	rm -rf build/*

.PHONY: fmt
fmt:
	go list -f {{.Dir}} ./... | xargs gofmt -w -s -d

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: test
test:
	go test ./... -timeout $(TIMEOUT)s

.PHONY: build
build: clean fmt lint test
	go build -o build/omnitrix main.go

