#!/bin/sh
action="$1"
if [ "$action" = "remove" ] || [ "$action" = "0" ]; then
    echo ""
    echo "========================================================================="
    echo " WARNING: To completely remove IRIS configurations and shell hooks,"
    echo " please run 'iris uninstall' BEFORE removing this package, or manually"
    echo " clean up ~/.config/iris and your shell RC files (e.g. ~/.zshrc)."
    echo "========================================================================="
    echo ""
fi
