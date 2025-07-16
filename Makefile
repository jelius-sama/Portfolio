CONFIG_FILE := ./config/config.json
CONFIG_DIR := $(dir $(CONFIG_FILE))

VERSION := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("version")) | .version' $(CONFIG_DIR){} )
APP_NAME := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("title")) | .title' $(CONFIG_DIR){} )
DEV_PORT := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("dev_port")) | .dev_port' $(CONFIG_DIR){} )

BIN_DIR := ./bin
ENV := production

.PHONY: portfolio dev archive_prod

portfolio:
	(cd client && bun run build)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X main.Environment=$(ENV) -X main.Port=$(DEV_PORT) -X main.Version=$(VERSION)" -trimpath -buildvcs=false -o $(BIN_DIR)/$(APP_NAME)-$(VERSION) ./cmd

dev:
	@echo "Starting development servers..."
	@( \
		pids=""; \
		cleanup() { \
			echo "Gracefully shutting down all servers..."; \
			for pid in $$pids; do \
				if kill -0 $$pid 2>/dev/null; then \
					kill -TERM $$pid 2>/dev/null || true; \
				fi; \
			done; \
			sleep 0.5; \
			for pid in $$pids; do \
				if kill -0 $$pid 2>/dev/null; then \
					kill -KILL $$pid 2>/dev/null || true; \
				fi; \
			done; \
			wait 2>/dev/null || true; \
			echo "All servers stopped."; \
			exit 0; \
		}; \
		trap cleanup INT TERM EXIT; \
		echo "Starting backend server..."; \
		air & pids="$$pids $$!"; \
		echo "Starting frontend dev server..."; \
		(cd client && bun run dev) & pids="$$pids $$!"; \
		echo "All development servers started! Press Ctrl+C to stop all servers."; \
		wait; \
	)

# Files and folders to include in the archive
INCLUDE_FILES = \
	bin \
	config \
	LICENSE \
	README.md \

archive_prod:
	tar -czf $(APP_NAME)-$(VERSION) $(INCLUDE_FILES)
