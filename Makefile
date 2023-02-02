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

.PHONY: check
check: check-license lint-go test-go

.PHONY: setup-git-hooks
setup-git-hooks:
	@cp scripts/git-hook/pre-commit .git/hooks/pre-commit

.PHONY: check-license
check-license:
	@./scripts/check-license.sh > /dev/null 2>&1
	@echo "License check passed"

.PHONY: lint-go
lint-go:
	@golangci-lint run > /dev/null 2>&1
	@echo "Go lint passed"

.PHONY: test-go
test-go:
	@go test -v -parallel 4 ./... > /dev/null 2>&1
	@echo "Go tests passed"
