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

ensure-tree-is-clean() {
    if ! git diff --quiet; then
        echo "Tree is dirty."
        exit 1
    fi
}

# Check out and update the main branch.
git checkout main
git pull

if ! main_version=$(scripts/detect-version.sh); then
    echo "Failed to detect version number."
    exit 1
fi

# Check out and update the develop branch.
git checkout develop
git pull origin main

# Ensure the develop branch is up to date with main.
if ! git merge-base --is-ancestor main develop; then
    echo "The develop branch is not up to date with main."
    exit 1
fi

# Get current version.
if ! version=$(scripts/detect-version.sh); then
    echo "Failed to detect version number."
    exit 1
fi

# If the version is not the same as the main branch, exit.
if [ "$version" != "$main_version" ]; then
    echo "The version is not the same as the main branch."
    echo "This probably means that the version has already been rolled."
    exit 1
fi

# Increment patch version.
new_version=$(echo "$version" | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')

# Replace version number in version.go.
sed -i "s/$version/$new_version/g" version.go

# Add a block to the top of CHANGELOG.md.
(
    cat <<EOF
## $new_version

TODO

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

EOF
) | cat - CHANGELOG.md > temp && mv temp CHANGELOG.md

# Commit changes.
git add CHANGELOG.md version.go
git commit -m "Bump version to $new_version"
git push
