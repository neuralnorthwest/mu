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

# Test for presence of gh.
if command -v gh >/dev/null 2>&1; then
    exit 0
fi

# Install gh.
if [ "$(uname)" = "Darwin" ]; then
    brew install gh
elif command -v apt-get >/dev/null 2>&1; then
    sudo apt-get install gh
else
    echo "Unsupported OS: $(uname)"
    exit 1
fi
