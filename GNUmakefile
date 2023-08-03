all: generate test

generate:
	go generate ./...
.PHONY: generate

# Run acceptance tests
test:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
.PHONY: test

