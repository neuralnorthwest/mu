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

# Set up all virtual environments.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

# Test for presence of python3 and virtualenv.
if ! command -v python3 >/dev/null 2>&1; then
    echo "python3 not found. Please install python3."
    exit 1
fi

if ! command -v virtualenv >/dev/null 2>&1; then
    echo "virtualenv not found. Please install virtualenv."
    exit 1
fi

# Remove all virtual environments if --fresh is passed.
if [ "${1:-}" = "--fresh" ]; then
    echo "Removing all virtual environments."
    find . -name venv -type d -exec rm -rf {} +
fi

# Create virtual environments.
find . -name requirements.txt -type f -exec dirname {} \; | while read -r dir; do
    if [ ! -d "$dir/venv" ]; then
        echo "Creating virtual environment in $dir/venv."
        virtualenv -p python3 "$dir/venv"
        "$dir/venv/bin/pip" install --upgrade pip
        "$dir/venv/bin/pip" install -r "$dir/requirements.txt"
    fi
done
