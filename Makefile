# Makefile for RayUI

BINARY_NAME ?= RayUI
BUILD_DIR ?= build/bin
ARTIFACTS_DIR ?= build/artifacts
ICON_SRC ?= build/appicon.png
WIN_ICON_OUT ?= build/windows/icon.ico

# VERSION:
# - Manual packaging: override via `make package-darwin VERSION=v1.2.3`
# - Default: best-effort from git
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo v0.1.0)

WAILS ?= $(shell command -v wails 2>/dev/null || echo "$(shell go env GOPATH)/bin/wails")

.PHONY: all clean icons icons-windows build-frontend build \
	build-darwin build-windows build-linux \
	package package-darwin package-windows package-linux \
	test test-backend test-frontend test-frontend-coverage test-all

all: build

clean:
	@rm -rf build/bin/*
	@rm -rf build/artifacts/*
	@rm -rf frontend/dist/*

icons: icons-windows

icons-windows:
	@echo "Normalizing $(ICON_SRC) to a real PNG (in-place) if needed"
	@go run ./scripts/normalize_appicon "$(ICON_SRC)" "$(ICON_SRC).tmp"
	@mv "$(ICON_SRC).tmp" "$(ICON_SRC)"
	@echo "Generating Windows icon from $(ICON_SRC) -> $(WIN_ICON_OUT)"
	@go run ./scripts/gen_windows_ico "$(ICON_SRC)" "$(WIN_ICON_OUT)"

build-frontend:
	@echo "Building frontend..."
	@cd frontend && \
		pnpm install && \
		pnpm run build

# build-darwin: clean build-frontend
# 	$(WAILS) build -trimpath -platform darwin/amd64,darwin/arm64 -ldflags "-X main.version=$(VERSION)"

# build-windows: clean build-frontend
# 	$(WAILS) build -trimpath -platform windows/amd64 -ldflags "-X main.version=$(VERSION)"

# package-windows: build-windows
# 	mkdir -p build/artifacts
# 	cd $(BUILD_DIR) && zip -j ../artifacts/RayUI-windows-$(VERSION).zip $(BINARY_NAME).exe

# build-linux: clean build-frontend
# 	$(WAILS) build -trimpath -platform linux/amd64 -ldflags "-X main.version=$(VERSION)"

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
	@go run ./scripts/zip_artifact "$(ARTIFACTS_DIR)/RayUI-macOS-$(VERSION).zip" "$(BUILD_DIR)/$(BINARY_NAME).app"

package-windows: build-windows
	@mkdir -p $(ARTIFACTS_DIR)
	@go run ./scripts/zip_artifact "$(ARTIFACTS_DIR)/RayUI-windows-$(VERSION).zip" "$(BUILD_DIR)/$(BINARY_NAME).exe"

package-linux: build-linux
	mkdir -p $(ARTIFACTS_DIR)
	cd $(BUILD_DIR) && tar -czf ../artifacts/RayUI-linux-$(VERSION).tar.gz $(BINARY_NAME)

package: package-darwin package-windows package-linux
	@echo "All packages built in build/artifacts"

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
