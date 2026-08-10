BUILD_DIR := ./build

# Config Options:
HOST := https://jelius.dev
ASSET_CDN_HOST := https://cdn.jelius.dev
VERSION := 3.0.4
PORT := :6969

DEV_BIN := $(BUILD_DIR)/portfolio-dev
BIN := $(BUILD_DIR)/portfolio

.PHONY: dev build

dev:
	@mkdir -p $(BUILD_DIR)
	tailwindcss -i ./template/input.css -o ./assets/css/output-$(VERSION).css
	bun run --cwd ./legacy/markdown/ build
	templ generate
	CGO_ENABLED=0 GOOS=linux go build -ldflags "\
		    -s -w \
		    -X main.Environment=development  \
		    -X main.Host=$(HOST) \
		    -X main.AssetCDNHost=http://shogun.local$(PORT) \
		    -X main.Version=$(VERSION) \
		    -X main.Port=$(PORT)" \
		    -trimpath -buildvcs=false -o $(DEV_BIN) ./cmd

build:
	@mkdir -p $(BUILD_DIR)
	tailwindcss -i ./template/input.css -o ./assets/css/output-$(VERSION).css
	bun run --cwd ./legacy/markdown/ build
	templ generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "\
		    -s -w \
		    -X main.Environment=production \
		    -X main.Host=$(HOST) \
		    -X main.AssetCDNHost=$(ASSET_CDN_HOST) \
		    -X main.Version=$(VERSION) \
			-X main.Port=:6500" \
		    -trimpath -buildvcs=false -o $(BIN) ./cmd
