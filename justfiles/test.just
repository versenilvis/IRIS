# run all tests
[group('test')]
test:
    @go test ./... -v

# run project health and scoring analyzer
alias ana := analyze
[group('test')]
analyze:
    @go run scripts/test_analyzer.go

# golangci-lint type checks with the go/types it was built against, so a newer
# toolchain on PATH makes it panic on the standard library rather than report
# anything. Override when your golangci-lint was built against a different Go.
lint-toolchain := env('IRIS_LINT_TOOLCHAIN', 'go1.26.5')

# run linter
[group('test')]
lint:
    @GOTOOLCHAIN={{lint-toolchain}} golangci-lint run ./...
