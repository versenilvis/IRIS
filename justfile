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
