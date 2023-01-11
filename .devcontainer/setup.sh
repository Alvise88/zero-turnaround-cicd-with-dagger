#!/usr/bin/env bash

echo "Installing Starship"

# curl -sS https://starship.rs/install.sh | sh - -y
sh -c "$(curl -fsSL https://starship.rs/install.sh)" - -y
grep -qF "/usr/local/bin/starship" ~/.bashrc || echo "$(starship init bash)" >> ~/.bashrc

go install github.com/magefile/mage@v1.14.0