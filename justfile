# build binary file
[group('dev')]
build:
    @go build -o iris main.go    
optimized-build:
    @GOAMD64=v3 go build -pgo=auto -ldflags="-s -w" -trimpath -o iris main.go

# run iris
[group('dev')]
run:
    @./iris

# re-build and reload iris
[group('dev')]
[linux, macos]
reload:
    @go build -o iris main.go
    @if [ -n "${IRIS_PID:-}" ]; then kill -USR1 $IRIS_PID; fi
    @if [ -z "${IRIS_FD:-}" ]; then ./iris; fi

# update pkg
[group('dev')]
pkg:
    @go mod tidy

# debugger
[group('dev')]
debug:
    @./iris --debug

# copy to local bin
[group('dev')]
copy:
    @rm ~/.local/bin/iris
    @cp ./iris ~/.local/bin/iris
# run all tests
[group('dev')]
test:
    @go test ./... -v

# run project health and scoring analyzer
alias ana := analyze
[group('dev')]
analyze:
    @go run scripts/test_analyzer.go

# run linter 
[group('dev')]
lint:
    @golangci-lint run ./...

# test the update command (version check + comparison), no full iris session needed
# usage: just debug-update v1.99.0
[group('debug')]
debug-update version="v1.99.0":
    #!/bin/sh
    tmp=$(mktemp -d)
    printf '{"tag_name":"%s"}' "{{version}}" > "$tmp/response.json"
    python3 -m http.server 19999 --directory "$tmp" 2>/dev/null &
    SERVER_PID=$!
    sleep 0.3
    echo "--- testing iris update command ---"
    IRIS_UPDATE_URL="http://localhost:19999/response.json" ./iris update
    kill $SERVER_PID 2>/dev/null
    rm -rf "$tmp"

# test the in-session update notification banner (requires iris.zsh hook to be active)
# usage: just debug-notify v1.99.0
[group('debug')]
debug-notify version="v1.99.0":
    IRIS_PID="" IRIS_MOCK_LATEST_VERSION="{{version}}" ./iris

# build a versioned release binary
# usage: just build-release v1.2.0
[group('debug')]
build-release version:
    @GOAMD64=v3 go build -pgo=auto -ldflags="-s -w -X github.com/versenilvis/iris/root.Version={{version}}" -trimpath -o iris main.go
