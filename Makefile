SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
DEPLOY_TARGET?=10.104.86.12
DEPLOY_USER?=root
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

.PHONY: release
release: clean
	./sbin/build_platform.sh grandctl

.PHONY: clean
clean:
	rm -rf dist

.PHONY: local
local:
	go build .
	go install

.PHONY: deploy
deploy: release
ifeq ($(DEPLOY_USER),ubuntu)
	scp ./dist/grandctl-linux-amd64 $(DEPLOY_USER)@$(DEPLOY_TARGET):/tmp
	ssh -t $(DEPLOY_USER)@$(DEPLOY_TARGET) 'sudo cp /tmp/grandctl-linux-amd64 /usr/local/bin/grandctl'
else
	scp ./dist/grandctl-linux-amd64 $(DEPLOY_USER)@$(DEPLOY_TARGET):/usr/local/bin/grandctl
endif

.DEFAULT_GOAL := release