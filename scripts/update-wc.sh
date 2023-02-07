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

# Count lines of Go code, convert to thousands (k), and update the README.md
# badge.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

WC=$(find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' | xargs wc -l | tail -n 1 | awk '{print $1}')
WCK=$(((WC + 500) / 1000))

# Update README.md, looking for `code-<N>k-blue` and replacing it with
# `code-${WCK}k-blue`.
sed -i "s/code-[0-9]\+k-blue/code-${WCK}k-blue/" README.md
