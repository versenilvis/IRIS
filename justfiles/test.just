# run all tests
[group('test')]
test:
    @go test ./... -v

# run project health and scoring analyzer
alias ana := analyze
[group('test')]
analyze:
    @go run scripts/test_analyzer.go

# run linter
[group('test')]
lint:
    @golangci-lint run ./...
