package model

// Config holds the global application configuration.
type Config struct {
	ActiveProfileID string `json:"activeProfileId"`
	ActiveRoutingID string `json:"activeRoutingId"`
	ActiveDNSPreset string `json:"activeDnsPreset"`

	CoreBasic CoreBasicConfig `json:"coreBasic"`
	Inbounds  []InboundConfig `json:"inbounds"`
	ProxyMode EProxyMode      `json:"proxyMode"`
	TUN       TUNConfig       `json:"tun"`

	SystemProxy SystemProxyConfig `json:"systemProxy"`
	UI          UIConfig          `json:"ui"`
	SpeedTest   SpeedTestConfig   `json:"speedTest"`
}

// CoreBasicConfig holds core-level settings.
type CoreBasicConfig struct {
	LogEnabled     bool   `json:"logEnabled"`
	LogLevel       string `json:"logLevel"`
	MuxEnabled     bool   `json:"muxEnabled"`
	AllowInsecure  bool   `json:"allowInsecure"`
	Fingerprint    string `json:"fingerprint"`
	EnableFragment bool   `json:"enableFragment"`
}

// InboundConfig describes one inbound listener.
type InboundConfig struct {
	Protocol        string `json:"protocol"`
	ListenAddr      string `json:"listenAddr"`
	Port            int    `json:"port"`
	UDPEnabled      bool   `json:"udpEnabled"`
	SniffingEnabled bool   `json:"sniffingEnabled"`
	AllowLAN        bool   `json:"allowLAN"`
}

// TUNConfig holds TUN-mode settings.
type TUNConfig struct {
	Enabled     bool   `json:"enabled"`
	AutoRoute   bool   `json:"autoRoute"`
	StrictRoute bool   `json:"strictRoute"`
	Stack       string `json:"stack"`
	MTU         int    `json:"mtu"`
	EnableIPv6  bool   `json:"enableIPv6"`
}

// SystemProxyConfig holds system proxy behaviour settings.
type SystemProxyConfig struct {
	Exceptions    string `json:"exceptions"`
	NotProxyLocal bool   `json:"notProxyLocal"`
}

// UIConfig holds user-interface preferences.
type UIConfig struct {
	Theme           string `json:"theme"`
	Language        string `json:"language"`
	FontFamily      string `json:"fontFamily"`
	FontSize        int    `json:"fontSize"`
	AutoHideOnStart bool   `json:"autoHideOnStart"`
	CloseToTray     bool   `json:"closeToTray"`
	ShowInDock      bool   `json:"showInDock"`
}

// SpeedTestConfig holds speed/latency test settings.
type SpeedTestConfig struct {
	Timeout    int    `json:"timeout"`
	URL        string `json:"url"`
	PingURL    string `json:"pingUrl"`
	Concurrent int    `json:"concurrent"`
}

// DefaultConfig returns a Config with all default values.
func DefaultConfig() Config {
	return Config{
		ActiveDNSPreset: "default",
		ProxyMode:       ProxyModeManual,
		CoreBasic: CoreBasicConfig{
			LogEnabled:  true,
			LogLevel:    "info",
			Fingerprint: "chrome",
		},
		Inbounds: []InboundConfig{
			{
				Protocol:        "socks",
				ListenAddr:      "127.0.0.1",
				Port:            10808,
				UDPEnabled:      true,
				SniffingEnabled: true,
			},
			{
				Protocol:        "http",
				ListenAddr:      "127.0.0.1",
				Port:            10809,
				SniffingEnabled: true,
			},
		},
		TUN: TUNConfig{
			AutoRoute:   true,
			StrictRoute: true,
			Stack:       "gvisor",
			MTU:         9000,
		},
		SystemProxy: SystemProxyConfig{
			NotProxyLocal: true,
		},
		UI: UIConfig{
			Theme:      "system",
			Language:   "en",
			FontSize:   14,
			CloseToTray: true,
			ShowInDock:  true,
		},
		SpeedTest: SpeedTestConfig{
			Timeout:    10,
			URL:        "https://speed.cloudflare.com/__down?bytes=10000000",
			PingURL:    "https://www.gstatic.com/generate_204",
			Concurrent: 4,
		},
	}
}
