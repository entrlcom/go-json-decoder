#!make

include .env
export $(shell sed 's/\#.*//g; /^$$/d; s/=.*//g' .env)

LOCAL_BIN := ${CURDIR}/bin
PATH := ${PATH}:${LOCAL_BIN}

GOLANGCI_LINT_VERSION := v2.1.2

export REPOSITORY_VERSION ?= `git describe --tags || echo "v0.0.1"`

.go-vuln-install:
	GOBIN=${LOCAL_BIN} go install golang.org/x/vuln/cmd/govulncheck@latest

.golangci-lint-install:
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

.PHONY: go-fmt
go-fmt: .golangci-lint-install
	${LOCAL_BIN}/golangci-lint fmt -c .golangci.yml

.PHONY: go-lint
go-lint: .golangci-lint-install
	${LOCAL_BIN}/golangci-lint run -c .golangci.yml ./...

.PHONY: go-test
go-test:
	go test -count=1 -covermode=atomic -coverprofile=cover.out.tmp -race ./internal/...
	cat cover.out.tmp | grep -v ".minimock.go" > cover.out
	@coverage=$$(go tool cover -func cover.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); echo "$$coverage"

.PHONY: go-vuln
go-vuln: .go-vuln-install
	${LOCAL_BIN}/govulncheck ./...
