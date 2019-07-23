build:
	CGO_ENABLED=0 go build -o go-dd

GO_TEST_PATHS := $(shell command go list ./... | grep -v "vendor")
test:
	go test $(GO_TEST_PATHS) -v

.DEFAULT_GOAL := build