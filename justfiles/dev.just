# run iris
[group('dev')]
run:
    @./iris

# initialize default config file
[group('dev')]
config-init:
    @./iris config init

# re-build and reload iris
[group('dev')]
[linux, macos]
reload:
    @go build -o iris ./cmd/iris
    @if [ -f ~/.local/bin/iris ]; then rm -f ~/.local/bin/iris && cp ./iris ~/.local/bin/iris; fi
    @if [ -n "${IRIS_PID:-}" ]; then kill -USR1 $IRIS_PID 2>/dev/null || true; fi
    @if [ -z "${IRIS_FD:-}" ]; then ./iris; fi

# copy to local bin
[group('dev')]
copy:
    @rm ~/.local/bin/iris
    @cp ./iris ~/.local/bin/iris
