SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
DEPLOY_TARGET?=10.104.86.12
OS=$(shell uname -s)

export PATH := ./bin:$(PATH)

.PHONY: setup
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
ifeq ($(OS), Darwin)
	brew install dep
else
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -vendor-only

.PHONY: test
test:
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:
	golangci-lint run --enable-all ./...

ci: lint test

.PHONY: build
build:
	go build .

.PHONY: install
install: build
	go install

.PHONY: deploy
deploy: build install
	scp ${GOPATH}/bin/grandctl root@${DEPLOY_TARGET}:/usr/local/bin

.DEFAULT_GOAL := install