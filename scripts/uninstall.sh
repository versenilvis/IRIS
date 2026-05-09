#!/bin/bash

echo "Uninstalling Iris..."


BIN_LOCATIONS=(
    "$HOME/.local/bin/iris"
    "/usr/local/bin/iris"
)

for loc in "${BIN_LOCATIONS[@]}"; do
    if [ -f "$loc" ]; then
        echo "Removing binary: $loc"
        if [ -w "$(dirname "$loc")" ]; then
            rm -f "$loc"
        else
            sudo rm -f "$loc"
        fi
    fi
done


CONFIG_FILES=(
    "$HOME/.zshrc"
    "$HOME/.bashrc"
    "$HOME/.config/fish/config.fish"
)

for file in "${CONFIG_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "Removing integration from $file..."
        sed -i '/# Iris Autocomplete/d' "$file"
        sed -i '/iris init/d' "$file"
        sed -i '/^$/N;/^\n$/D' "$file"
    fi
done

if [ -f "iris.log" ]; then
    rm -f "iris.log"
fi

echo "✓ Iris has been successfully uninstalled"
echo "Please restart your terminal or source your config file to clear the environment"
