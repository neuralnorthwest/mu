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

# Check for files missing license headers and not in .licenseignore.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

NO_LIC_FILES=$(grep --exclude-dir=.git --exclude-dir=venv --exclude-dir=gen --exclude-dir=mock --exclude-dir=docs -HLr 'Licensed under the Apache License' . | sort)
IGNORE_FILES=$(cat .licenseignore | sort)
if [ "$NO_LIC_FILES" != "$IGNORE_FILES" ]; then
    echo "The following files are missing license headers and are not in .licenseignore:" >&2
    echo "$NO_LIC_FILES" | grep -v -F -f .licenseignore >&2
    exit 1
fi
