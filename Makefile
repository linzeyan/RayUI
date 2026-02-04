# Makefile for RayUI

BINARY_NAME ?= RayUI
BUILD_DIR ?= build/bin
ARTIFACTS_DIR ?= build/artifacts
ICON_SRC ?= build/appicon.png
WIN_ICON_OUT ?= build/windows/icon.ico

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo v0.1.0)
WAILS ?= $(shell command -v wails 2>/dev/null || echo "$(shell go env GOPATH)/bin/wails")

.PHONY: all clean icons build-frontend build \
	build-darwin build-windows build-linux \
	package package-darwin package-windows package-linux \
	test test-backend test-frontend test-frontend-coverage test-all

all: build

clean:
	@rm -rf $(BUILD_DIR)/*
	@rm -rf $(ARTIFACTS_DIR)/*
	@rm -rf frontend/dist/*

icons:
	@echo "Normalizing $(ICON_SRC) to a real PNG (in-place) if needed"
	@go run ./scripts/normalize_appicon "$(ICON_SRC)" "$(ICON_SRC).tmp"
	@mv "$(ICON_SRC).tmp" "$(ICON_SRC)"
	@echo "Generating Windows icon from $(ICON_SRC) -> $(WIN_ICON_OUT)"
	@go run ./scripts/gen_windows_ico "$(ICON_SRC)" "$(WIN_ICON_OUT)"

build-frontend:
	@echo "Building frontend..."
	@cd frontend && pnpm install && pnpm run build

build: clean build-frontend icons
	@echo "Building application..."
	$(WAILS) build -trimpath -ldflags "-X main.AppVersion=$(VERSION)"

build-darwin: clean build-frontend icons
	$(WAILS) build -trimpath -platform darwin/amd64,darwin/arm64 -ldflags "-X main.AppVersion=$(VERSION)"

build-windows: clean build-frontend icons
	$(WAILS) build -trimpath -platform windows/amd64 -ldflags "-X main.AppVersion=$(VERSION)"

build-linux: clean build-frontend
	$(WAILS) build -trimpath -platform linux/amd64 -ldflags "-X main.AppVersion=$(VERSION)"

package-darwin: build-darwin
	@mkdir -p $(ARTIFACTS_DIR)
	@# Create universal binary from amd64 and arm64 builds
	@cp -R "$(BUILD_DIR)/$(BINARY_NAME)-arm64.app" "$(BUILD_DIR)/$(BINARY_NAME).app"
	@lipo -create \
		"$(BUILD_DIR)/$(BINARY_NAME)-amd64.app/Contents/MacOS/$(BINARY_NAME)" \
		"$(BUILD_DIR)/$(BINARY_NAME)-arm64.app/Contents/MacOS/$(BINARY_NAME)" \
		-output "$(BUILD_DIR)/$(BINARY_NAME).app/Contents/MacOS/$(BINARY_NAME)"
	@go run ./scripts/zip_artifact "$(ARTIFACTS_DIR)/RayUI-macOS-$(VERSION).zip" "$(BUILD_DIR)/$(BINARY_NAME).app"

package-windows: build-windows
	@mkdir -p $(ARTIFACTS_DIR)
	@go run ./scripts/zip_artifact "$(ARTIFACTS_DIR)/RayUI-windows-$(VERSION).zip" "$(BUILD_DIR)/$(BINARY_NAME).exe"

package-linux: build-linux
	@mkdir -p $(ARTIFACTS_DIR)
	@cd $(BUILD_DIR) && tar -czf ../../$(ARTIFACTS_DIR)/RayUI-linux-$(VERSION).tar.gz $(BINARY_NAME)

package: package-darwin package-windows package-linux
	@echo "All packages built in $(ARTIFACTS_DIR)"

# ─── Testing ───────────────────────────────────────────────────

test-backend:
	@echo "Running Go tests..."
	@go test ./... -v

test-frontend:
	@echo "Running frontend tests..."
	@cd frontend && pnpm test:run

test-frontend-coverage:
	@echo "Running frontend tests with coverage..."
	@cd frontend && pnpm test:coverage

test: test-backend test-frontend

test-all: test
