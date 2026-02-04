package model

// EConfigType represents a proxy protocol type.
type EConfigType int

const (
	ConfigVMess       EConfigType = 1
	ConfigShadowsocks EConfigType = 3
	ConfigSOCKS       EConfigType = 4
	ConfigVLESS       EConfigType = 5
	ConfigTrojan      EConfigType = 6
	ConfigHysteria2   EConfigType = 7  // P1
	ConfigTUIC        EConfigType = 8  // P1
	ConfigWireGuard   EConfigType = 9  // P1
	ConfigHTTP        EConfigType = 10
)

func (e EConfigType) String() string {
	switch e {
	case ConfigVMess:
		return "vmess"
	case ConfigShadowsocks:
		return "shadowsocks"
	case ConfigSOCKS:
		return "socks"
	case ConfigVLESS:
		return "vless"
	case ConfigTrojan:
		return "trojan"
	case ConfigHysteria2:
		return "hysteria2"
	case ConfigTUIC:
		return "tuic"
	case ConfigWireGuard:
		return "wireguard"
	case ConfigHTTP:
		return "http"
	default:
		return "unknown"
	}
}

// ECoreType represents a proxy core engine.
type ECoreType int

const (
	CoreAuto    ECoreType = 0
	CoreXray    ECoreType = 1
	CoreSingbox ECoreType = 2
)

func (e ECoreType) String() string {
	switch e {
	case CoreAuto:
		return "auto"
	case CoreXray:
		return "xray"
	case CoreSingbox:
		return "sing-box"
	default:
		return "unknown"
	}
}

// EProxyMode represents a proxy interception mode.
type EProxyMode int

const (
	ProxyModeManual EProxyMode = 0
	ProxyModeSystem EProxyMode = 1
	ProxyModeTUN    EProxyMode = 2
	ProxyModePAC    EProxyMode = 3 // P2
)

func (e EProxyMode) String() string {
	switch e {
	case ProxyModeManual:
		return "manual"
	case ProxyModeSystem:
		return "system"
	case ProxyModeTUN:
		return "tun"
	case ProxyModePAC:
		return "pac"
	default:
		return "unknown"
	}
}
