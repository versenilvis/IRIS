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
    @if [ -z "${IRIS_FD:-}" ]; then ./iris; fi

# update pkg
[group('dev')]
pkg:
    @go mod tidy