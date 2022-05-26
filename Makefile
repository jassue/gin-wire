GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
APP_MAIN_DIR=cmd/app

.PHONY: run
# run
run:
	cd $(APP_MAIN_DIR) && go run main.go wire_gen.go app.go

.PHONY: build
# 自动根据平台编译二进制文件
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# 生成应用所需的文件
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: wire
# wire
wire:
	cd $(APP_MAIN_DIR) && wire

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help