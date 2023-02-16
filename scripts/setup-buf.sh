#!/bin/bash

# Copyright 2023 Scott M. Long
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Install gh

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

BUF_VERSION=${1:-1.14.0}
INSTALL_PREFIX=/usr/local

install-buf-linux() {
    curl -sSL \
        "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
        sudo tar -xzf - -C "${INSTALL_PREFIX}" --strip-components 1
}

install-buf-macos() {
    brew install bufbuild/buf/buf@${BUF_VERSION}
}

upgrade-buf-linux() {
    sudo rm -f \
        "${INSTALL_PREFIX}/bin/buf" \
        "${INSTALL_PREFIX}/bin/protoc-gen-buf-breaking" \
        "${INSTALL_PREFIX}/bin/protoc-gen-buf-lint" \
        "${INSTALL_PREFIX}/etc/bash_completion.d/buf" \
        "${INSTALL_PREFIX}/etc/fish/vendor_completions.d/buf.fish" \
        "${INSTALL_PREFIX}/etc/zsh/site-functions/_buf"
    install-buf-linux
}

upgrade-buf-macos() {
    brew upgrade bufbuild/buf/buf@${BUF_VERSION}
}

# Install buf.
if [ "$(uname)" = "Darwin" ]; then
    if command -v brew >/dev/null 2>&1; then
        if brew list bufbuild/buf/buf@${BUF_VERSION} >/dev/null 2>&1; then
            upgrade-buf-macos
        else
            install-buf-macos
        fi
    else
        echo "Please install Homebrew and re-run this script."
        exit 1
    fi
elif command -v apt-get >/dev/null 2>&1; then
    if command -v buf >/dev/null 2>&1; then
        if [ "$(buf --version)" = "${BUF_VERSION}" ]; then
            exit 0
        fi
        upgrade-buf-linux
    else
        install-buf-linux
    fi
else
    echo "Unsupported OS: $(uname)"
    exit 1
fi

echo "buf installed at version $(buf --version)"
