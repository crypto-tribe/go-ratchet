.POSIX:

TEST_CMD = go list -f '{{.Dir}}/...' -m | xargs go test

bin/:
	mkdir -p $@

bin/golangci-lint: | bin/
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh

.PHONY: clean
clean:
	rm -rf bin test-coverage test-coverage.html

.PHONY: lint
lint: bin/golangci-lint
	go list -f '{{.Dir}}/...' -m | xargs ./bin/golangci-lint run --fix

.PHONY: test
test:
	$(TEST_CMD)

.PHONY: test-coverage
test-coverage:
	$(TEST_CMD) -coverprofile test-coverage

.PHONY: test-coverage-html
test-coverage-html: test-coverage
	go tool cover -html=$< -o $<.html
	xdg-open $<.html
