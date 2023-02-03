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

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

# Test for presence of gh.
if ! command -v gh >/dev/null 2>&1; then
    echo "gh not found. Please install it with `make setup-dev`."
    exit 1
fi

if ! version=$(scripts/detect-version.sh); then
    echo "Failed to detect version number."
    exit 1
fi

# Ensure no release exists.
if gh release view "$version" >/dev/null 2>&1; then
    fail_vercheck=1
    echo "Release $version already exists."
fi

# Ensure no tag exists.
if [ "$(git tag -l \"$version\")" == "$version" ]; then
    fail_vercheck=1
    echo "Tag $version already exists."
fi

if [ "${fail_vercheck:-}" = "1" ]; then
    exit 1
fi

# Scrape the changelog for the release notes. Grab the block from
# `## $version` to the next `## `.
if ! grep -q "## $version" CHANGELOG.md; then
    echo "No release notes found in CHANGELOG.md."
    exit 1
fi
if ! notes=$(sed -n "/^## $version/,/^## /p" CHANGELOG.md | sed '$d'); then
    echo "Failed to scrape release notes from CHANGELOG.md."
    exit 1
fi

# Create the release.
if ! gh release create "$version" -t "github.com/neuralnorthwest/mu@$version" -n "$notes"; then
    echo "Failed to create release."
    exit 1
fi
