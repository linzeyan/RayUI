# RayUI

A modern, cross-platform desktop proxy client built with Go and Wails v2. RayUI provides a clean interface for managing proxy connections through [Xray-core](https://github.com/XTLS/Xray-core) and [sing-box](https://github.com/SagerNet/sing-box), with automatic core selection based on protocol.

## Features

- **Dual-core support** — Xray-core and sing-box with automatic selection
- **Multi-protocol** — VMess, VLESS, Trojan, Shadowsocks, Hysteria2, TUIC, WireGuard
- **Subscription management** — Base64, SIP008, sing-box JSON, Clash YAML formats
- **Routing rules** — Built-in presets (Global, Bypass LAN, Bypass CN) with GeoIP/GeoSite
- **DNS configuration** — Remote, direct, and bootstrap DNS with FakeIP support
- **Proxy modes** — Manual, System Proxy, TUN (via sing-box)
- **Traffic monitoring** — Real-time upload/download statistics
- **Speed testing** — TCP ping and download speed tests
- **System tray** — Quick toggle, show/hide, quit
- **Internationalization** — English, 正體中文
- **Themes** — Light, dark, and system-follow

## Supported Platforms

| Platform | Architecture |
|----------|-------------|
| macOS    | amd64, arm64 (Universal Binary) |
| Windows  | amd64 |
| Linux    | amd64 |

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) 20+
- [pnpm](https://pnpm.io/) 9+
- [Wails CLI](https://wails.io/) v2

### Platform-specific

**macOS**
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian)**
```bash
sudo apt-get install -y build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev libglib2.0-dev
```

**Windows**
- [Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) with C++ workload

## Getting Started

### 1. Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Clone and build

```bash
git clone https://github.com/RayUI/RayUI.git
cd RayUI
make build
```

The built application is in `build/bin/`.

### 3. Run in development mode

```bash
wails dev
```

This starts the Go backend and frontend dev server with hot reload.

### 4. Install proxy cores

Launch RayUI, go to **Settings > Core**, and click **Check Update** for Xray-core and/or sing-box. The cores are downloaded automatically to `~/.RayUI/cores/`.

### 5. Add a profile

- Click **Add Profile** to manually enter server details, or
- Go to **Subscriptions** and add a subscription URL to import nodes in bulk

### 6. Connect

Select a profile and click the connect button. Choose your proxy mode in **Settings > Proxy**:

- **Manual** — Configure your applications to use SOCKS5 (default 10808) or HTTP (default 10809)
- **System Proxy** — Automatically sets OS-level proxy settings
- **TUN** — Transparent proxy for all traffic (requires sing-box)

## Build Commands

```bash
make build              # Build for current platform
make build-darwin       # macOS (amd64 + arm64)
make build-windows      # Windows (amd64)
make build-linux        # Linux (amd64)

make package-darwin     # Build + create macOS universal .zip
make package-windows    # Build + create Windows .zip
make package-linux      # Build + create Linux .tar.gz
make package            # Package all platforms

make test               # Run all tests (Go + frontend)
make test-backend       # Go tests only
make test-frontend      # Frontend tests only
```

Override the version tag:
```bash
make package-darwin VERSION=v1.2.3
```

## Project Structure

```
RayUI/
├── main.go                  # Wails entry point
├── wails.json               # Wails configuration
├── Makefile                 # Build automation
├── frontend/                # React + TypeScript frontend
│   ├── src/
│   │   ├── pages/           # Profiles, Subscriptions, Routing, DNS, Settings, Logs
│   │   ├── components/      # shadcn/ui based components
│   │   ├── stores/          # Zustand state management
│   │   ├── hooks/           # Custom React hooks
│   │   ├── i18n/            # Translations (en, zh-TW)
│   │   └── lib/             # Utilities
│   └── wailsjs/             # Auto-generated Wails bindings
├── internal/
│   ├── app/                 # Wails app bindings
│   ├── core/                # Core process manager (xray, sing-box)
│   ├── config/              # Core config generators
│   ├── model/               # Data models
│   ├── store/               # JSON file persistence
│   ├── service/             # Business logic (subscriptions, updater, stats, geodata)
│   ├── parser/              # Protocol URI parsers
│   ├── sysproxy/            # OS-level proxy configuration
│   ├── autostart/           # Auto-start at login
│   ├── tray/                # System tray
│   └── security/            # Encryption utilities
├── scripts/                 # Build helpers (icon generation, zip packaging)
└── build/                   # Build resources and output
```

## Data Directory

All user data is stored in `~/.RayUI/`:

```
~/.RayUI/
├── config.json              # Global settings
├── profiles.json            # Proxy profiles
├── subscriptions.json       # Subscription sources
├── routing.json             # Routing rules
├── dns.json                 # DNS configuration
├── cores/                   # Core binaries (xray, sing-box)
├── data/                    # GeoIP/GeoSite databases
└── logs/                    # Application logs
```

## Tech Stack

**Backend** — Go, [Wails v2](https://github.com/wailsapp/wails)

**Frontend** — React 19, TypeScript, [Tailwind CSS v4](https://tailwindcss.com/), [shadcn/ui](https://ui.shadcn.com/) (Radix UI), Zustand, [Rsbuild](https://rsbuild.dev/)

**Testing** — `go test`, [Vitest](https://vitest.dev/), Testing Library

## Acknowledgements

- [Wails](https://github.com/wailsapp/wails) — Desktop application framework for Go
- [XTLS/Xray-core](https://github.com/XTLS/Xray-core) — Proxy core with VLESS, VMess, Trojan support
- [SagerNet/sing-box](https://github.com/SagerNet/sing-box) — Universal proxy platform
- [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) — GeoIP/GeoSite data for Xray
- [SagerNet/sing-geoip](https://github.com/SagerNet/sing-geoip) / [sing-geosite](https://github.com/SagerNet/sing-geosite) — GeoIP/GeoSite data for sing-box
- [shadcn/ui](https://ui.shadcn.com/) — UI component library
- [Radix UI](https://www.radix-ui.com/) — Accessible component primitives
- [Tailwind CSS](https://tailwindcss.com/) — Utility-first CSS framework

## License

MIT
