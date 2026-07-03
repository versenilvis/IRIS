#!/bin/sh

echo "Uninstalling Iris..."

if command -v iris >/dev/null 2>&1; then
    if iris uninstall --help >/dev/null 2>&1; then
        iris uninstall
        exit 0
    fi
fi

for loc in "${HOME}/.local/bin/iris" "/usr/local/bin/iris"; do
    if [ -f "${loc}" ]; then
        echo "Removing binary: ${loc}"
        if [ -w "$(dirname "${loc}")" ]; then
            /bin/rm -f "${loc}" 2>/dev/null || rm -f "${loc}"
        else
            sudo /bin/rm -f "${loc}" 2>/dev/null || sudo rm -f "${loc}"
        fi
    fi
done

for file in "${HOME}/.zshrc" "${HOME}/.bashrc" "${HOME}/.config/fish/config.fish"; do
    if [ -f "${file}" ]; then
        echo "Removing integration from ${file}..."
        tmp_file=$(mktemp)
        grep -v -i -E "(# iris autocomplete|# iris autostart|iris init)" "${file}" > "${tmp_file}" 2>/dev/null
        status=$?
        if [ "${status}" -eq 0 ] || [ "${status}" -eq 1 ]; then
            mv "${tmp_file}" "${file}"
        else
            /bin/rm -f "${tmp_file}" 2>/dev/null || rm -f "${tmp_file}"
        fi
    fi
done

/bin/rm -rf "${HOME}/.config/iris" 2>/dev/null || rm -rf "${HOME}/.config/iris"
/bin/rm -rf "${HOME}/.local/share/iris" 2>/dev/null || rm -rf "${HOME}/.local/share/iris"
/bin/rm -rf "${HOME}/.cache/iris" 2>/dev/null || rm -rf "${HOME}/.cache/iris"
/bin/rm -f "iris.log" 2>/dev/null || rm -f "iris.log"

echo "✓ Iris has been successfully uninstalled"
if [ -n "${IRIS_PID}" ]; then
    echo ""
    echo "⚠️  You are currently inside an active Iris session."
    echo "Iris runs as the parent process of this terminal - do NOT run 'pkill iris'"
    echo "as it will immediately close this terminal window."
    echo ""
    echo "To fully exit, simply close this terminal window and open a new one."
    echo "Iris will not start again since the shell config has been cleaned up."
else
    echo "Please close and reopen your terminal to complete the uninstall."
fi
