#!/bin/bash
# Description: Environment initialization for devcontainer (chown, pre-commit, terraform cache)
# Usage: ./init.sh

sudo chown -R "$(id -u)":"$(id -g)" /home/vscode/.aws
sudo chown -R "$(id -u)":"$(id -g)" /home/vscode/.gitconfig
sudo chown -R "$(id -u)":"$(id -g)" /home/vscode/.ssh
chmod 600 /home/vscode/.ssh/id_rsa

# for precommit
pre-commit install

# for terraform
mkdir -p "$HOME/.terraform.d/plugin-cache"
if command -v tfenv >/dev/null 2>&1; then
    tfenv use || { echo "Error: tfenv use failed." >&2; exit 1; }
else
    echo "tfenv not found. Skipping tfenv use."
fi
