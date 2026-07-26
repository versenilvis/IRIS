# build binary file
[group('build')]
build:
    @go build -o iris ./cmd/iris

# build with optimized binary file
[group('build')]
optimized-build:
    @GOAMD64=v4 go build -ldflags="-s -w" -trimpath -o iris ./cmd/iris

# build a versioned release binary
# usage: just build-release v1.2.0
[group('build')]
build-release version:
    @GOAMD64=v3 go build -pgo=auto -ldflags="-s -w -X github.com/versenilvis/iris/root.Version={{version}}" -trimpath -o iris ./cmd/iris
