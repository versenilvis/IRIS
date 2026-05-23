# NOTE: RUN "just -l" TO QUICKLY VIEW ALL COMMANDS

# build binary file
[group('build')]
build:
    @go build -o iris main.go
    
# build with optimized binary file
[group('build')]
optimized-build:
    @GOAMD64=v4 go build -ldflags="-s -w" -trimpath -o iris main.go

# run iris
[group('dev')]
run:
    @./iris

# re-build and reload iris
[group('dev')]
[linux, macos]
reload:
    @go build -o iris main.go
    @if [ -n "${IRIS_PID:-}" ]; then kill -USR1 $IRIS_PID 2>/dev/null || true; fi
    @if [ -z "${IRIS_FD:-}" ]; then ./iris; fi

# update pkg
[group('dev')]
pkg:
    @go mod tidy

# iris debugger
[group('debug')]
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

# test the install script locally
# usage: just debug-install v1.0.0
[group('debug')]
debug-install version="v1.0.0":
    #!/bin/sh
    PORT=19998
    TMP=$(mktemp -d)
    PASS=0
    FAIL=0

    # always cleanup server + tmp, even if a test crashes
    trap "kill \$SERVER_PID 2>/dev/null; rm -rf $TMP" EXIT

    ok()   { PASS=$((PASS+1)); printf "  \033[32m✓\033[0m %s\n" "$1"; }
    fail() { FAIL=$((FAIL+1)); printf "  \033[31m✗\033[0m %s\n" "$1"; }

    # --- detect actual arch (not hardcoded) ---
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    MACH=$(uname -m)
    case "$MACH" in
        x86_64|amd64)  ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "unsupported arch: $MACH"; exit 1 ;;
    esac
    ARCHIVE_NAME="iris_${OS}_${ARCH}.tar.gz"

    echo "building iris ({{version}}, ${OS}_${ARCH})..."
    go build -ldflags="-X github.com/versenilvis/iris/root.Version={{version}}" -o "$TMP/iris" main.go

    echo "packaging archive ($ARCHIVE_NAME)..."
    tar -czf "$TMP/$ARCHIVE_NAME" -C "$TMP" iris

    # --- start mock server ---
    python3 -c 'import sys, textwrap; exec(textwrap.dedent(r"""
    import os, sys
    from http.server import HTTPServer, BaseHTTPRequestHandler

    port = int(sys.argv[1])
    tmp_dir = sys.argv[2]
    version = sys.argv[3]
    os_name = sys.argv[4]
    arch = sys.argv[5]

    class mock_handler(BaseHTTPRequestHandler):
        def do_GET(self):
            if "/repos/versenilvis/iris/releases/latest" in self.path:
                if self.path.startswith("/404"):
                    self.send_response(404)
                    self.send_header("Content-Type", "application/json")
                    self.end_headers()
                    self.wfile.write(b"{\"message\":\"Not Found\",\"documentation_url\":\"https://docs.github.com/rest\"}")
                elif self.path.startswith("/ratelimit"):
                    self.send_response(403)
                    self.send_header("Content-Type", "application/json")
                    self.end_headers()
                    self.wfile.write(b"{\"message\":\"API rate limit exceeded for ...\",\"documentation_url\":\"https://docs.github.com/rest/overview/rate-limits-for-the-rest-api\"}")
                else:
                    self.send_response(200)
                    self.send_header("Content-Type", "application/json")
                    self.end_headers()
                    download_url = "http://localhost:" + str(port) + "/iris_" + os_name + "_" + arch + ".tar.gz"
                    response = "{\n  \"tag_name\": \"" + version + "\",\n  \"assets\": [\n    {\n      \"browser_download_url\": \"" + download_url + "\"\n    }\n  ]\n}\n"
                    self.wfile.write(response.encode("utf-8"))
            elif self.path.endswith(".tar.gz"):
                filename = os.path.basename(self.path)
                filepath = os.path.join(tmp_dir, filename)
                if os.path.exists(filepath):
                    self.send_response(200)
                    self.send_header("Content-Type", "application/octet-stream")
                    self.send_header("Content-Length", str(os.path.getsize(filepath)))
                    self.end_headers()
                    with open(filepath, "rb") as f:
                        self.wfile.write(f.read())
                else:
                    self.send_response(404)
                    self.end_headers()
            else:
                self.send_response(200)
                self.end_headers()

        def log_message(self, *args):
            pass

    HTTPServer(("", port), mock_handler).serve_forever()
    """))' "$PORT" "$TMP" "{{version}}" "$OS" "$ARCH" 2>/dev/null &
    SERVER_PID=$!

    # retry loop instead of blind sleep - wait up to 3s for server to be ready
    i=0
    until curl -sf "http://localhost:${PORT}/" >/dev/null 2>&1; do
        i=$((i+1))
        [ $i -ge 30 ] && echo "mock server failed to start" && exit 1
        sleep 0.1
    done

    echo ""
    echo "--- running test cases ---"

    # --- test 1: happy path ---
    printf "test 1: happy path install... "
    OUT=$(BIN_DIR="$TMP/out1" IRIS_API_URL="http://localhost:${PORT}/happy" sh scripts/install.sh 2>&1) && STATUS=0 || STATUS=$?
    if [ $STATUS -eq 0 ] && echo "$OUT" | grep -q "Installation verified"; then
        ok "binary installed and verified"
    else
        fail "happy path failed\n$OUT"
    fi

    # --- test 2: version output check ---
    printf "test 2: version string matches... "
    BIN="$TMP/out1/iris"
    if [ -x "$BIN" ]; then
        VER=$("$BIN" version 2>&1)
        if echo "$VER" | grep -q "{{version}}"; then
            ok "got '$VER'"
        else
            fail "expected {{version}}, got '$VER'"
        fi
    else
        fail "binary not found at $BIN"
    fi

    # --- test 3: 404 (no release published) ---
    printf "test 3: 404 no release error... "
    OUT=$(BIN_DIR="$TMP/out3" IRIS_API_URL="http://localhost:${PORT}/404" sh scripts/install.sh 2>&1) && STATUS=0 || STATUS=1
    if [ $STATUS -ne 0 ] && echo "$OUT" | grep -qi "no releases found\|not have published"; then
        ok "correct error message for 404"
    else
        fail "expected 404 error, got: $OUT"
    fi

    # --- test 4: 403 rate limit ---
    printf "test 4: rate limit error... "
    OUT=$(BIN_DIR="$TMP/out4" IRIS_API_URL="http://localhost:${PORT}/ratelimit" sh scripts/install.sh 2>&1) && STATUS=0 || STATUS=1
    if [ $STATUS -ne 0 ] && echo "$OUT" | grep -qi "rate limit\|rate limited"; then
        ok "correct error message for rate limit"
    else
        fail "expected rate limit error, got: $OUT"
    fi

    # --- test 5: wget code path (call wget directly, same as install script would) ---
    printf "test 5: wget code path... "
    if ! command -v wget >/dev/null 2>&1; then
        ok "skipped (wget not installed)"
    else
        # directly invoke wget the same way install.sh does and check it parses correctly
        RELEASES=$(wget -qO- "http://localhost:${PORT}/happy/repos/versenilvis/iris/releases/latest" 2>&1) && WS=0 || WS=$?
        URL=$(echo "$RELEASES" | grep "browser_download_url" | grep "${OS}_${ARCH}" | head -1 | cut -d '"' -f 4)
        if [ -n "$URL" ]; then
            ok "wget parses download URL correctly ($URL)"
        else
            fail "wget could not parse URL from response: $RELEASES"
        fi
    fi

    # --- summary ---
    echo ""
    printf "results: \033[32m%d passed\033[0m, \033[31m%d failed\033[0m\n" "$PASS" "$FAIL"
    [ $FAIL -eq 0 ] || exit 1