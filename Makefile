BUILD_DIR := ./build

# Config Options:
HOST := jelius.dev
VERSION := 2.0.1
PORT := :6969

DEV_BIN := $(BUILD_DIR)/portfolio-dev
BIN := $(BUILD_DIR)/portfolio

.PHONY: dev build

dev:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags "\
		    -s -w \
		    -X main.Environment=development  \
		    -X main.Host=$(HOST) \
		    -X main.Version=$(VERSION) \
		    -X main.Port=$(PORT)" \
		    -trimpath -buildvcs=false -o $(DEV_BIN) ./cmd

build:
	@mkdir -p $(BUILD_DIR)
	templ generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "\
		    -s -w \
		    -X main.Environment=production \
		    -X main.Host=$(HOST) \
		    -X main.Version=$(VERSION) \
		    -X main.Port=$(PORT)" \
		    -trimpath -buildvcs=false -o $(BIN) ./cmd
