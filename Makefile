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

# Targets:
#
# check: runs all pre-commit checks
# setup-dev: sets up the development environment
# setup-git-hooks: sets up the git hooks
# setup-venv: sets up all virtual environments
# setup-gh: sets up the github cli
# check-license: checks the license headers
# lint-go: runs the go linter
# test-go: runs the go tests
# release: releases the project (called from CI)

.PHONY: check
check: generate check-license lint test

.PHONY: generate
generate: generate-go generate-proto

.PHONY: generate-proto
generate-proto:
	@find . -type d -name proto | while read -r dir; do \
		echo "Generating $$dir"; \
		cd $$dir; \
		buf generate; \
	done
	@echo "Buf generate passed"

.PHONY: generate-go
generate-go:
	@go generate ./...

.PHONY: setup-dev
setup-dev: setup-git-hooks setup-venv setup-gh setup-buf

.PHONY: setup-git-hooks
setup-git-hooks:
	@cp scripts/git-hook/* .git/hooks/

.PHONY: setup-venv
setup-venv:
	@./scripts/setup-venv.sh

.PHONY: setup-gh
setup-gh:
	@./scripts/setup-gh.sh

.PHONY: setup-buf
setup-buf:
	@./scripts/setup-buf.sh

.PHONY: check-license
check-license:
	@./scripts/check-license.sh > /dev/null
	@echo "License check passed"

.PHONY: lint
lint: lint-proto lint-go

.PHONY: lint-proto
lint-proto:
	@find . -type d -name proto -exec buf lint {} \;
	@echo "Buf lint passed"

.PHONY: lint-go
lint-go:
	@golangci-lint run > /dev/null 2>&1
	@echo "Go lint passed"

.PHONY: test
test: test-go

.PHONY: test-go
test-go:
	@go test -v -parallel 4 ./... > /dev/null 2>&1
	@echo "Go tests passed"

.PHONY: release
release:
	@./scripts/release.sh

.PHONY: roll-version
roll-version:
	@./scripts/roll-version.sh
