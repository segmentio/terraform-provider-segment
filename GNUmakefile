# Runs all the necessary targets to ensure code is in good shape
all: deps generate fmt lint test

# Downloads build dependencies
deps:
	@go mod download
.PHONY: deps

# Generates documentation
generate:
	@go generate ./...
.PHONY: generate

# Reformats the code (using gofumpt)
fmt:
	@go run mvdan.cc/gofumpt@latest -l -w .
.PHONY: fmt

GOLANGCI_LINT_VERSION='1.53.3'

# Lints the source code and tests
lint:
ifdef CI
	@bash -c '(\
			(which golangci-lint && golangci-lint --version) | grep -q $(GOLANGCI_LINT_VERSION) ||\
			curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh |\
			sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION) \
		) >\
		/dev/null && PATH="$$PATH:$(shell go env GOPATH)/bin" \
		golangci-lint run'
else
	@golangci-lint run
endif
.PHONY: lint

coverage-threshold = ( \
	bash -c '\
	COVERAGE_THRESHOLD=89;\
	ACTUAL_COVERAGE=`go tool cover -func=coverage.out | grep total | grep -Eo "[0-9]{2,3}"`; \
	if [ $$ACTUAL_COVERAGE -ge $$COVERAGE_THRESHOLD ]; then \
	  echo "✅ Code coverage OK: $$ACTUAL_COVERAGE% (threshold $$COVERAGE_THRESHOLD%)"; \
	else \
	  echo -e "⛔️ Code coverage below threshold: $$ACTUAL_COVERAGE%, should be at least $$COVERAGE_THRESHOLD%\
	  \nPlease inspect the coverage output (with \`make coverage\`) and add more tests!" && exit -1; \
	fi' \
)

# Runs acceptance tests
test:
	@TF_ACC=1 gotestsum -f short-verbose -- -v $(TESTARGS) -timeout 120m -coverpkg=./... -coverprofile coverage.out -covermode=atomic -cover ./...
	@$(call coverage-threshold)
.PHONY: test

# Opens the HTML code coverage report in the default browser.
coverage:
	@go tool cover -html=coverage.out
.PHONY: coverage

