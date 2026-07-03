#!/bin/sh
set -e

REPO="versenilvis/iris"
BIN_DIR="${BIN_DIR:-/usr/local/bin}"
# allow overriding the GitHub API base URL for local testing
IRIS_API_URL="${IRIS_API_URL:-https://api.github.com}"

main() {
    echo "Installing iris..."

    arch=$(get_arch)
    echo "Detected architecture: ${arch}"

    tmp_dir=$(mktemp -d)
    trap 'rm -rf '"${tmp_dir}" EXIT
    cd "${tmp_dir}"

    download_url=$(get_download_url "${arch}")
    if [ -z "${download_url}" ]; then
        err "Could not find release for architecture: ${arch}"
    fi
    echo "Downloading: ${download_url}"

    if command -v curl >/dev/null 2>&1; then
        curl -sLO "${download_url}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${download_url}"
    else
        err "curl or wget is required"
    fi

    archive=$(basename "${download_url}")
    case "${archive}" in
        *.tar.gz) tar -xzf "${archive}" ;;
        *.zip)
            if command -v unzip >/dev/null 2>&1; then
                unzip -q "${archive}"
            else
                err "unzip is required to extract ${archive}"
            fi
            ;;
        *) err "Unknown archive format: ${archive}" ;;
    esac

    bin=$(find . -name "iris" -type f | head -1)
    [ -z "$bin" ] && err "Binary not found in archive"

    # check if we have write permission to the install directory
    has_write_permission=0
    if [ -d "${BIN_DIR}" ]; then
        if [ -w "${BIN_DIR}" ]; then
            has_write_permission=1
        fi
    else
        parent_dir=$(dirname "${BIN_DIR}")
        if [ -w "${parent_dir}" ]; then
            has_write_permission=1
        fi
    fi

    if [ "${has_write_permission}" -eq 1 ]; then
        mkdir -p "${BIN_DIR}"
        cp "$bin" "${BIN_DIR}/iris"
        chmod +x "${BIN_DIR}/iris"
        if "${BIN_DIR}/iris" version >/dev/null 2>&1; then
            echo "Installation verified."
            echo ""
            "${BIN_DIR}/iris" setup
        else
            echo "Warning: could not verify installed binary at ${BIN_DIR}/iris"
        fi
    else
        # fallback to ~/.local/bin which is user-writable without sudo
        local_bin="${HOME}/.local/bin"
        mkdir -p "${local_bin}"
        chmod +x "$bin"
        cp "$bin" "${local_bin}/iris"
        chmod +x "${local_bin}/iris"
        if "${local_bin}/iris" version >/dev/null 2>&1; then
            echo "Installation verified."
            echo ""
            "${local_bin}/iris" setup
        else
            # both locations failed, sudo install
            echo ""
            echo "Installation requires elevated permissions, enter your password:"
            if sudo cp "$bin" "${BIN_DIR}/iris" && sudo chmod +x "${BIN_DIR}/iris"; then
                echo "Installation verified."
                echo ""
                "${BIN_DIR}/iris" setup
            else
                tmp_iris=$(mktemp "${TMPDIR:-/tmp}/iris.XXXXXX")
                cp "$bin" "${tmp_iris}"
                chmod +x "${tmp_iris}"
                echo ""
                printf "Failed. Run manually: \033[32msudo cp %s %s/iris\033[0m\n" "${tmp_iris}" "${BIN_DIR}"
            fi
        fi
    fi
}

get_arch() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "${os}" in
        linux)  os="linux" ;;
        darwin) os="darwin" ;;
        *) err "Unsupported OS: ${os}" ;;
    esac

    case "${arch}" in
        x86_64 | amd64)   arch="amd64" ;;
        aarch64 | arm64)  arch="arm64" ;;
        *) err "Unsupported architecture: ${arch}" ;;
    esac

    echo "${os}_${arch}"
}

get_download_url() {
    arch="$1"

    if command -v curl >/dev/null 2>&1; then
        http_response=$(curl -sL -w "\n%{http_code}" \
            ${GITHUB_TOKEN:+-H "Authorization: Bearer ${GITHUB_TOKEN}"} \
            "${IRIS_API_URL}/repos/${REPO}/releases/latest")
        http_code=$(echo "${http_response}" | tail -1)
        releases=$(echo "${http_response}" | sed '$d')
    elif command -v wget >/dev/null 2>&1; then
        tmp_headers=$(mktemp)
        releases=$(wget -S -qO- \
            ${GITHUB_TOKEN:+--header "Authorization: Bearer ${GITHUB_TOKEN}"} \
            "${IRIS_API_URL}/repos/${REPO}/releases/latest" 2>"$tmp_headers" || true)
        http_code=$(grep "HTTP/" "$tmp_headers" | tail -1 | sed -e 's/^[[:space:]]*//' | cut -d' ' -f2)
        [ -z "${http_code}" ] && http_code="000"
        rm -f "$tmp_headers"
    else
        err "curl or wget is required"
    fi

    if [ "${http_code}" = "404" ]; then
        err "no releases found for ${REPO}. the project may not have published a release yet"
    fi

    if [ "${http_code}" = "403" ] || echo "${releases}" | grep -q "rate limit"; then
        err "GitHub API rate limited. try again later or set GITHUB_TOKEN env variable"
    fi

    if [ "${http_code}" != "200" ]; then
        msg=$(echo "${releases}" | grep '"message"' | head -1 | cut -d '"' -f 4)
        err "GitHub API error (HTTP ${http_code}): ${msg}"
    fi

    url=$(echo "${releases}" | grep "browser_download_url" | grep "${arch}" | head -1 | cut -d '"' -f 4)
    echo "${url}"
}

err() {
    echo "Error: $1" >&2
    exit 1
}

main "$@"