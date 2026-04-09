# build binary file
[group('dev')]
build:
    @go build -o iris main.go

# run iris
[group('dev')]
run:
    @./iris