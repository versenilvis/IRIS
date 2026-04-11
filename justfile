# build binary file
[group('dev')]
build:
    @GOAMD64=v3 go build -pgo=auto -ldflags="-s -w" -trimpath -o iris main.go

# run iris
[group('dev')]
run:
    @./iris
