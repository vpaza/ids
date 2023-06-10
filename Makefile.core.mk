FINDFILES=find . \( -path ./.git -o -path ./out -o -path ./.github -o -path ./vendor -o -path ./frontend/node_modules \) -prune -o -type f
XARGS=xargs -0 -r
RELEASE_LDFLAGS='-extldflags -static -s -w'
BINARIES=./cmd/api

lint-copyright:
	@${FINDFILES} \( -name '*.go' -o -name '*.sh' \) \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 |\
		${XARGS} scripts/lint_copyright_license.sh

lint-go:
	@${FINDFILES} -name '*.go' \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 | ${XARGS} scripts/lint_go.sh

lint-markdown:
	@${FINDFILES} -name '*.md' -print0 | ${XARGS} mdl --ignore-front-matter --style .mdl.rb

frontend-lint:
	@cd frontend && yarn lint

frontend-lint-fix:
	@cd frontend && yarn lint-fix

frontend-format:
	@cd frontend && yarn format

lint: lint-copyright lint-go lint-markdown

fix-copyright:
	@${FINDFILES} \( -name '*.go' -o -name '*.sh' \) \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 |\
		${XARGS} scripts/fix_copyright_license.sh

.PHONY: default
default: init build

.PHONY: init
init:
	@mkdir -p out

.PHONY: frontend-build
frontend-build:
	@cd frontend && yarn && yarn build

.PHONY: backend-build
backend-build:
	@LDFLAGS=${RELEASE_LDFLAGS} scripts/gobuild.sh out/ ${BINARIES}

.PHONY: assemble-release
assemble-release:
	@cp -r static out/

.PHONY: build
build: frontend-build backend-build assemble-release build-image

.PHONY: mod-vendor
mod-vendor:
	@go mod vendor

.PHONY: dev
dev:
	@go run ./cmd/api/main.go server

.PHONY: clean
clean:
	@rm -rf out
	@rm -rf frontend/dist
	@rm -rf static

.PHONY: dist-clean
dist-clean: clean
	@rm -rf vendor
	@rm -rf frontend/node_modules
