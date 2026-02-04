package parser

import (
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
	"gopkg.in/yaml.v3"
)

// clashConfig represents the minimal Clash YAML structure for proxy parsing.
type clashConfig struct {
	Proxies []clashProxy `yaml:"proxies"`
}

// clashProxy represents a single proxy entry in Clash format.
type clashProxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password,omitempty"`
	UUID     string `yaml:"uuid,omitempty"`
	AlterID  int    `yaml:"alterId,omitempty"`
	Cipher   string `yaml:"cipher,omitempty"`
	// TLS
	TLS            bool   `yaml:"tls,omitempty"`
	SkipCertVerify bool   `yaml:"skip-cert-verify,omitempty"`
	ServerName     string `yaml:"servername,omitempty"`
	SNI            string `yaml:"sni,omitempty"`
	ALPN           []string `yaml:"alpn,omitempty"`
	Fingerprint    string `yaml:"client-fingerprint,omitempty"`
	// Transport
	Network string `yaml:"network,omitempty"`
	// WS options
	WSOpts *clashWSOptions `yaml:"ws-opts,omitempty"`
	// gRPC options
	GRPCOpts *clashGRPCOptions `yaml:"grpc-opts,omitempty"`
	// HTTP/2 options
	H2Opts *clashH2Options `yaml:"h2-opts,omitempty"`
	// VMess specific
	Security string `yaml:"security,omitempty"`
	// VLESS specific
	Flow string `yaml:"flow,omitempty"`
	// Reality
	RealityOpts *clashRealityOptions `yaml:"reality-opts,omitempty"`
	// Hysteria2
	Auth     string `yaml:"auth,omitempty"`
	Obfs     string `yaml:"obfs,omitempty"`
	ObfsPass string `yaml:"obfs-password,omitempty"`
	// TUIC
	CongestionController string `yaml:"congestion-controller,omitempty"`
	UDPRelayMode         string `yaml:"udp-relay-mode,omitempty"`
	// WireGuard
	PrivateKey string   `yaml:"private-key,omitempty"`
	PublicKey  string   `yaml:"public-key,omitempty"`
	IP         string   `yaml:"ip,omitempty"`
	IPv6       string   `yaml:"ipv6,omitempty"`
	Reserved   []int    `yaml:"reserved,omitempty"`
	MTU        int      `yaml:"mtu,omitempty"`
	Peers      []interface{} `yaml:"peers,omitempty"` // ignored for now
}

type clashWSOptions struct {
	Path    string            `yaml:"path,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type clashGRPCOptions struct {
	ServiceName string `yaml:"grpc-service-name,omitempty"`
}

type clashH2Options struct {
	Host []string `yaml:"host,omitempty"`
	Path string   `yaml:"path,omitempty"`
}

type clashRealityOptions struct {
	PublicKey string `yaml:"public-key,omitempty"`
	ShortID   string `yaml:"short-id,omitempty"`
}

// parseClashYAML parses Clash YAML proxies into ProfileItems.
func parseClashYAML(data string) ([]model.ProfileItem, error) {
	var cfg clashConfig
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}

	var items []model.ProfileItem
	for _, proxy := range cfg.Proxies {
		item := clashProxyToProfile(proxy)
		if item != nil {
			items = append(items, *item)
		}
	}
	return items, nil
}

func clashProxyToProfile(p clashProxy) *model.ProfileItem {
	item := model.NewProfileItem()
	item.Remarks = p.Name
	item.Address = p.Server
	item.Port = p.Port

	switch strings.ToLower(p.Type) {
	case "vmess":
		item.ConfigType = model.ConfigVMess
		item.UUID = p.UUID
		item.AlterID = p.AlterID
		item.Security = p.Cipher
		if item.Security == "" {
			item.Security = "auto"
		}
	case "vless":
		item.ConfigType = model.ConfigVLESS
		item.UUID = p.UUID
		item.Flow = p.Flow
		item.Security = "none"
	case "trojan":
		item.ConfigType = model.ConfigTrojan
		item.UUID = p.Password
	case "ss", "shadowsocks":
		item.ConfigType = model.ConfigShadowsocks
		item.Security = p.Cipher
		item.UUID = p.Password
	case "hysteria2", "hy2":
		item.ConfigType = model.ConfigHysteria2
		item.CoreType = model.CoreSingbox
		item.UUID = p.Password
		if item.UUID == "" {
			item.UUID = p.Auth
		}
		item.HeaderType = p.Obfs
		item.Path = p.ObfsPass
	case "tuic":
		item.ConfigType = model.ConfigTUIC
		item.CoreType = model.CoreSingbox
		item.UUID = p.UUID
		item.Security = p.Password
		item.HeaderType = p.CongestionController
		item.Path = p.UDPRelayMode
	case "wireguard", "wg":
		item.ConfigType = model.ConfigWireGuard
		item.CoreType = model.CoreSingbox
		item.UUID = p.PrivateKey
		item.PublicKey = p.PublicKey
		item.Network = "wireguard"
		if p.IP != "" {
			item.Host = p.IP
			if p.IPv6 != "" {
				item.Host += "," + p.IPv6
			}
		}
		if len(p.Reserved) > 0 {
			parts := make([]string, len(p.Reserved))
			for i, v := range p.Reserved {
				parts[i] = strconv.Itoa(v)
			}
			item.ShortID = strings.Join(parts, ",")
		}
		if p.MTU > 0 {
			item.Extra = strconv.Itoa(p.MTU)
		}
		item.StreamSecurity = "none"
		return &item
	default:
		return nil // unsupported type
	}

	// Common TLS/transport fields.
	if p.TLS {
		item.StreamSecurity = "tls"
	}
	if p.SkipCertVerify {
		item.AllowInsecure = true
	}
	sni := p.SNI
	if sni == "" {
		sni = p.ServerName
	}
	item.SNI = sni
	if len(p.ALPN) > 0 {
		item.ALPN = strings.Join(p.ALPN, ",")
	}
	item.Fingerprint = p.Fingerprint

	// Reality.
	if p.RealityOpts != nil {
		item.StreamSecurity = "reality"
		item.PublicKey = p.RealityOpts.PublicKey
		item.ShortID = p.RealityOpts.ShortID
	}

	// Transport.
	net := p.Network
	if net == "" {
		net = "tcp"
	}
	item.Network = net

	switch net {
	case "ws":
		if p.WSOpts != nil {
			item.Path = p.WSOpts.Path
			if host, ok := p.WSOpts.Headers["Host"]; ok {
				item.Host = host
			}
		}
	case "grpc":
		if p.GRPCOpts != nil {
			item.Path = p.GRPCOpts.ServiceName
		}
	case "h2":
		if p.H2Opts != nil {
			item.Path = p.H2Opts.Path
			if len(p.H2Opts.Host) > 0 {
				item.Host = p.H2Opts.Host[0]
			}
		}
	}

	return &item
}
