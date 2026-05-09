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
