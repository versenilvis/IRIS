# iris debugger
[group('debug')]
debug:
    @./iris --debug

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

# test the changelog command (fetch, grouping, cache, rate limiting)
# usage: just debug-changelog v1.99.0
[group('debug')]
debug-changelog version="v1.99.0":
    #!/bin/sh
    tmp=$(mktemp -d)
    body='## Changelog\n### Bug fixes\n* abc1234  fix something\n\n## Update\n```bash\niris update\n```\n'
    printf '[{"tag_name":"%s","published_at":"2026-01-01T00:00:00Z","body":"%s","prerelease":false},{"tag_name":"v0.0.1","published_at":"2025-12-01T00:00:00Z","body":"## Changelog\\n### Documentation\\n* def5678  older release\\n","prerelease":false}]' "{{version}}" "$body" > "$tmp/response.json"
    python3 -m http.server 19997 --directory "$tmp" 2>/dev/null &
    SERVER_PID=$!
    sleep 0.3

    # isolate the changelog cache so this doesn't read or write the real
    # ~/.cache/iris/changelog-cache.json
    export XDG_CACHE_HOME="$tmp/cache"

    echo "--- testing iris changelog (latest only) ---"
    IRIS_CHANGELOG_URL="http://localhost:19997/response.json" ./iris changelog

    echo ""
    echo "--- testing iris changelog -n 2 ---"
    IRIS_CHANGELOG_URL="http://localhost:19997/response.json" ./iris changelog -n 2

    echo ""
    echo "--- testing cache hit (server killed, should still work) ---"
    kill $SERVER_PID 2>/dev/null
    IRIS_CHANGELOG_URL="http://localhost:19997/response.json" ./iris changelog

    echo ""
    echo "--- testing --refresh with no server (should error, not hang) ---"
    IRIS_CHANGELOG_URL="http://localhost:19997/response.json" ./iris changelog --refresh

    rm -rf "$tmp"

# test the in-session update notification banner (requires iris.zsh hook to be active)
# usage: just debug-notify v1.99.0
[group('debug')]
debug-notify version="v1.99.0":
    IRIS_PID="" IRIS_MOCK_LATEST_VERSION="{{version}}" ./iris

# test the background auto-updater's escalation ladder (silent install,
# then repeated failures escalating to confirm, then giving up). fully
# offline - mocks the release API, install script, and release asset.
# usage: just debug-autoupdate
[group('debug')]
debug-autoupdate:
    #!/bin/sh
    PORT=19996
    TMP=$(mktemp -d)
    PASS=0
    FAIL=0
    trap "kill \$SERVER_PID 2>/dev/null; pkill -f \"$TMP/iris\" 2>/dev/null; rm -rf $TMP" EXIT

    ok()   { PASS=$((PASS+1)); printf "  \033[32m✓\033[0m %s\n" "$1"; }
    fail() { FAIL=$((FAIL+1)); printf "  \033[31m✗\033[0m %s\n" "$1"; }

    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    MACH=$(uname -m)
    case "$MACH" in
        x86_64|amd64)  ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "unsupported arch: $MACH"; exit 1 ;;
    esac
    ARCHIVE_NAME="iris_${OS}_${ARCH}.tar.gz"

    echo "building old binary (v0.0.1)..."
    go build -ldflags="-X github.com/versenilvis/iris/root.Version=v0.0.1" -o "$TMP/iris" ./cmd/iris

    echo "building new binary (v1.99.0) + archive..."
    mkdir -p "$TMP/new"
    go build -ldflags="-X github.com/versenilvis/iris/root.Version=v1.99.0" -o "$TMP/new/iris" ./cmd/iris
    tar -czf "$TMP/$ARCHIVE_NAME" -C "$TMP/new" iris

    cp scripts/install.sh "$TMP/install.sh"
    printf '#!/bin/sh\nexit 1\n' > "$TMP/install-fail.sh"

    # --- mock server: release JSON, install script, and the tarball ---
    python3 -c 'import sys, textwrap; exec(textwrap.dedent(r"""
    import os, sys
    from http.server import HTTPServer, BaseHTTPRequestHandler

    port = int(sys.argv[1])
    tmp_dir = sys.argv[2]
    archive_name = sys.argv[3]

    class mock_handler(BaseHTTPRequestHandler):
        def do_GET(self):
            if self.path == "/update-response.json":
                self.send_response(200)
                self.send_header("Content-Type", "application/json")
                self.end_headers()
                body = "## Changelog\\n### Bug fixes\\n* abc1234  a test fix\\n"
                resp = "{\"tag_name\":\"v1.99.0\",\"published_at\":\"2026-01-01T00:00:00Z\",\"body\":\"" + body + "\",\"prerelease\":false}"
                self.wfile.write(resp.encode("utf-8"))
            elif self.path in ("/install.sh", "/install-fail.sh"):
                self.send_response(200)
                self.send_header("Content-Type", "text/plain")
                self.end_headers()
                with open(os.path.join(tmp_dir, self.path.lstrip("/")), "rb") as f:
                    self.wfile.write(f.read())
            elif "/repos/versenilvis/iris/releases/latest" in self.path:
                self.send_response(200)
                self.send_header("Content-Type", "application/json")
                self.end_headers()
                download_url = "http://localhost:" + str(port) + "/" + archive_name
                self.wfile.write(("{\"tag_name\":\"v1.99.0\",\"assets\":[{\"browser_download_url\":\"" + download_url + "\"}]}").encode("utf-8"))
            elif self.path == "/" + archive_name:
                filepath = os.path.join(tmp_dir, archive_name)
                self.send_response(200)
                self.send_header("Content-Type", "application/octet-stream")
                self.end_headers()
                with open(filepath, "rb") as f:
                    self.wfile.write(f.read())
            else:
                self.send_response(404)
                self.end_headers()

        def log_message(self, *args):
            pass

    HTTPServer(("", port), mock_handler).serve_forever()
    """))' "$PORT" "$TMP" "$ARCHIVE_NAME" 2>/dev/null &
    SERVER_PID=$!

    i=0
    until curl -sf "http://localhost:${PORT}/update-response.json" >/dev/null 2>&1; do
        i=$((i+1))
        [ $i -ge 30 ] && echo "mock server failed to start" && exit 1
        sleep 0.1
    done

    mkdir -p "$TMP/cfg/iris"
    cat > "$TMP/cfg/iris/config.toml" << 'CFGEOF'
    [updater]
    check-on-startup = true
    channel = "stable"
    check-interval = "0s"
    auto-update = 1
    CFGEOF

    # fresh process = fresh check cycle, since check-interval is 0s
    run_check() {
        HOME="$TMP/home" \
        XDG_CONFIG_HOME="$TMP/cfg" \
        XDG_DATA_HOME="$TMP/data" \
        BIN_DIR="$TMP/installed" \
        IRIS_UPDATE_URL="http://localhost:${PORT}/update-response.json" \
        IRIS_INSTALL_URL="$1" \
        "$TMP/iris" >"$TMP/out.log" 2>&1 &
        RUN_PID=$!
        sleep 1.5
        kill $RUN_PID 2>/dev/null
        wait $RUN_PID 2>/dev/null
    }

    state_field() {
        # tolerate BurntSushi/toml's indentation of nested-table fields
        grep -E "^[[:space:]]*$1" "$TMP/data/iris/state.toml" 2>/dev/null | head -1
    }

    echo ""
    echo "--- running test cases ---"

    # --- test 1: attempt 1 installs silently and succeeds ---
    printf "test 1: silent auto-install (attempt 1)... "
    run_check "http://localhost:${PORT}/install.sh"
    if [ -x "$TMP/installed/iris" ] && "$TMP/installed/iris" version 2>&1 | grep -q "v1.99.0"; then
        ok "binary auto-installed to v1.99.0"
    else
        fail "binary was not installed; see $TMP/out.log"
    fi
    if state_field "auto-update-attempt" | grep -q " 0$"; then
        ok "attempt counter cleared after success"
    else
        fail "expected attempt counter 0 after success, got: $(state_field auto-update-attempt)"
    fi

    # --- tests 2-5: forced failures escalate, then give up instead of looping ---
    rm -rf "$TMP/data" "$TMP/installed"

    printf "test 2: attempt 1 (forced failure)... "
    run_check "http://localhost:${PORT}/install-fail.sh"
    if state_field "auto-update-attempt" | grep -q " 1$"; then
        ok "attempt counter is 1"
    else
        fail "expected attempt counter 1, got: $(state_field auto-update-attempt)"
    fi

    printf "test 3: attempt 2 (forced failure again, escalates)... "
    run_check "http://localhost:${PORT}/install-fail.sh"
    if state_field "auto-update-attempt" | grep -q " 2$"; then
        ok "attempt counter is 2 (would show a confirm prompt)"
    else
        fail "expected attempt counter 2, got: $(state_field auto-update-attempt)"
    fi

    printf "test 4: attempt 3 (gives up)... "
    run_check "http://localhost:${PORT}/install-fail.sh"
    if state_field "auto-update-attempt" | grep -q " 3$"; then
        ok "attempt counter is 3 (gave up)"
    else
        fail "expected attempt counter 3, got: $(state_field auto-update-attempt)"
    fi

    printf "test 5: attempt 4+ (stays given up, no infinite retry)... "
    run_check "http://localhost:${PORT}/install-fail.sh"
    if state_field "auto-update-attempt" | grep -q " 3$"; then
        ok "attempt counter stayed at 3"
    else
        fail "expected attempt counter to stay at 3, got: $(state_field auto-update-attempt)"
    fi

    echo ""
    printf "results: \033[32m%d passed\033[0m, \033[31m%d failed\033[0m\n" "$PASS" "$FAIL"
    echo ""
    echo "note: the confirm-prompt UI (auto-update=2, or auto-update=1's 2nd"
    echo "attempt) is interactive and needs a real terminal with iris shell"
    echo "hooks active to see rendered - set auto-update in your real config"
    echo "and watch a live session for that part."

    [ $FAIL -eq 0 ] || exit 1

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
    go build -ldflags="-X github.com/versenilvis/iris/root.Version={{version}}" -o "$TMP/iris" ./cmd/iris

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
