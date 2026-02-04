package config

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// SingboxConfigGenerator produces a sing-box JSON config.
type SingboxConfigGenerator struct{}

func (g *SingboxConfigGenerator) Generate(
	profile model.ProfileItem,
	routing model.RoutingItem,
	dns model.DNSItem,
	cfg model.Config,
) ([]byte, error) {
	root := map[string]any{
		"log":          buildSingboxLog(cfg),
		"dns":          buildSingboxDNS(dns),
		"inbounds":     buildSingboxInbounds(cfg),
		"outbounds":    buildSingboxOutbounds(profile),
		"route":        buildSingboxRoute(routing),
		"experimental": buildSingboxExperimental(),
	}
	return json.MarshalIndent(root, "", "  ")
}

func buildSingboxLog(cfg model.Config) map[string]any {
	level := cfg.CoreBasic.LogLevel
	if level == "" {
		level = "info"
	}
	return map[string]any{
		"level":     level,
		"timestamp": true,
	}
}

func buildSingboxDNS(dns model.DNSItem) map[string]any {
	servers := []map[string]any{
		{"tag": "remote-dns", "address": dns.RemoteDNS, "detour": "proxy"},
		{"tag": "direct-dns", "address": dns.DirectDNS, "detour": "direct"},
	}
	if dns.BootstrapDNS != "" {
		servers = append(servers, map[string]any{
			"tag": "bootstrap-dns", "address": dns.BootstrapDNS, "detour": "direct",
		})
	}

	rules := []map[string]any{
		{"outbound": "any", "server": "direct-dns"},
	}

	result := map[string]any{
		"servers": servers,
		"rules":   rules,
		"final":   "remote-dns",
	}
	if dns.DomainStrategy != "" {
		result["strategy"] = dns.DomainStrategy
	}
	if dns.FakeIP {
		result["fakeip"] = map[string]any{
			"enabled":    true,
			"inet4_range": "198.18.0.0/15",
			"inet6_range": "fc00::/18",
		}
	}
	return result
}

func buildSingboxInbounds(cfg model.Config) []map[string]any {
	var inbounds []map[string]any

	// TUN inbound (if enabled).
	if cfg.TUN.Enabled {
		tun := map[string]any{
			"type":                     "tun",
			"tag":                      "tun-in",
			"inet4_address":            "172.19.0.1/30",
			"auto_route":              cfg.TUN.AutoRoute,
			"strict_route":            cfg.TUN.StrictRoute,
			"stack":                   cfg.TUN.Stack,
			"mtu":                     cfg.TUN.MTU,
			"sniff":                   true,
			"sniff_override_destination": false,
		}
		if cfg.TUN.EnableIPv6 {
			tun["inet6_address"] = "fdfe:dcba:9876::1/126"
		}
		inbounds = append(inbounds, tun)
	}

	// Mixed/SOCKS/HTTP inbounds from config.
	for _, ib := range cfg.Inbounds {
		listen := ib.ListenAddr
		if listen == "" {
			listen = "127.0.0.1"
		}
		if ib.AllowLAN {
			listen = "0.0.0.0"
		}
		inbound := map[string]any{
			"type":        "mixed",
			"tag":         ib.Protocol + "-in",
			"listen":      listen,
			"listen_port": ib.Port,
			"sniff":       ib.SniffingEnabled,
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds
}

func buildSingboxOutbounds(profile model.ProfileItem) []map[string]any {
	proxy := buildSingboxProxyOutbound(profile)
	direct := map[string]any{"type": "direct", "tag": "direct"}
	block := map[string]any{"type": "block", "tag": "block"}
	dns := map[string]any{"type": "dns", "tag": "dns-out"}
	return []map[string]any{proxy, direct, block, dns}
}

func buildSingboxProxyOutbound(p model.ProfileItem) map[string]any {
	ob := map[string]any{
		"tag":         "proxy",
		"server":      p.Address,
		"server_port": p.Port,
	}

	switch p.ConfigType {
	case model.ConfigVMess:
		ob["type"] = "vmess"
		ob["uuid"] = p.UUID
		ob["alter_id"] = p.AlterID
		sec := p.Security
		if sec == "" {
			sec = "auto"
		}
		ob["security"] = sec

	case model.ConfigVLESS:
		ob["type"] = "vless"
		ob["uuid"] = p.UUID
		if p.Flow != "" {
			ob["flow"] = p.Flow
		}

	case model.ConfigTrojan:
		ob["type"] = "trojan"
		ob["password"] = p.UUID

	case model.ConfigShadowsocks:
		ob["type"] = "shadowsocks"
		ob["method"] = p.Security
		ob["password"] = p.UUID

	case model.ConfigHysteria2:
		ob["type"] = "hysteria2"
		ob["password"] = p.UUID
		if p.HeaderType != "" {
			ob["obfs"] = map[string]any{
				"type":     p.HeaderType,
				"password": p.Path,
			}
		}

	case model.ConfigTUIC:
		ob["type"] = "tuic"
		ob["uuid"] = p.UUID
		ob["password"] = p.Security
		if p.HeaderType != "" {
			ob["congestion_control"] = p.HeaderType
		}
		udpRelay := p.Path
		if udpRelay == "" {
			udpRelay = "native"
		}
		ob["udp_relay_mode"] = udpRelay

	case model.ConfigWireGuard:
		ob["type"] = "wireguard"
		ob["private_key"] = p.UUID
		if p.PublicKey != "" {
			ob["peer_public_key"] = p.PublicKey
		}
		if p.Host != "" {
			addrs := strings.Split(p.Host, ",")
			var localAddrs []string
			for _, a := range addrs {
				a = strings.TrimSpace(a)
				if a != "" {
					localAddrs = append(localAddrs, a)
				}
			}
			if len(localAddrs) > 0 {
				ob["local_address"] = localAddrs
			}
		}
		if p.ShortID != "" {
			parts := strings.Split(p.ShortID, ",")
			var reserved []int
			for _, s := range parts {
				if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
					reserved = append(reserved, n)
				}
			}
			if len(reserved) > 0 {
				ob["reserved"] = reserved
			}
		}
		if p.Extra != "" {
			if mtu, err := strconv.Atoi(p.Extra); err == nil {
				ob["mtu"] = mtu
			}
		}
		return ob // WireGuard doesn't use standard transport/TLS
	}

	// Transport
	if p.Network != "" && p.Network != "tcp" {
		transport := map[string]any{"type": singboxTransportType(p.Network)}
		switch p.Network {
		case "ws":
			if p.Path != "" {
				transport["path"] = p.Path
			}
			if p.Host != "" {
				transport["headers"] = map[string]any{"Host": p.Host}
			}
		case "h2":
			transport["type"] = "http"
			if p.Host != "" {
				transport["host"] = []string{p.Host}
			}
			if p.Path != "" {
				transport["path"] = p.Path
			}
		case "grpc":
			if p.Path != "" {
				transport["service_name"] = p.Path
			}
		case "httpupgrade":
			if p.Host != "" {
				transport["host"] = p.Host
			}
			if p.Path != "" {
				transport["path"] = p.Path
			}
		}
		ob["transport"] = transport
	}

	// TLS
	if p.StreamSecurity == "tls" || p.StreamSecurity == "reality" {
		tls := map[string]any{"enabled": true}
		if p.SNI != "" {
			tls["server_name"] = p.SNI
		}
		if p.AllowInsecure {
			tls["insecure"] = true
		}
		if p.ALPN != "" {
			tls["alpn"] = strings.Split(p.ALPN, ",")
		}
		if p.Fingerprint != "" {
			tls["utls"] = map[string]any{"fingerprint": p.Fingerprint}
		}
		if p.StreamSecurity == "reality" {
			tls["reality"] = map[string]any{
				"enabled":    true,
				"public_key": p.PublicKey,
				"short_id":   p.ShortID,
			}
		}
		ob["tls"] = tls
	}

	return ob
}

func buildSingboxRoute(routing model.RoutingItem) map[string]any {
	var rules []map[string]any
	for _, r := range routing.Rules {
		if !r.Enabled {
			continue
		}
		rule := map[string]any{"outbound": r.OutboundTag}
		if len(r.Domain) > 0 {
			rule["domain"] = r.Domain
		}
		if len(r.DomainSuffix) > 0 {
			rule["domain_suffix"] = r.DomainSuffix
		}
		if len(r.DomainKeyword) > 0 {
			rule["domain_keyword"] = r.DomainKeyword
		}
		if len(r.DomainRegex) > 0 {
			rule["domain_regex"] = r.DomainRegex
		}
		if len(r.Geosite) > 0 {
			rule["geosite"] = r.Geosite
		}
		if len(r.IPCIDR) > 0 {
			rule["ip_cidr"] = r.IPCIDR
		}
		if len(r.GeoIP) > 0 {
			rule["geoip"] = r.GeoIP
		}
		if r.Port != "" {
			rule["port"] = r.Port
		}
		if len(r.Protocol) > 0 {
			rule["protocol"] = r.Protocol
		}
		if len(r.ProcessName) > 0 {
			rule["process_name"] = r.ProcessName
		}
		if r.Network != "" {
			rule["network"] = r.Network
		}
		if len(r.RuleSet) > 0 {
			rule["rule_set"] = r.RuleSet
		}
		rules = append(rules, rule)
	}

	result := map[string]any{
		"final":              "proxy",
		"auto_detect_interface": true,
	}
	if len(rules) > 0 {
		result["rules"] = rules
	}
	return result
}

func buildSingboxExperimental() map[string]any {
	return map[string]any{
		"clash_api": map[string]any{
			"external_controller": "127.0.0.1:9090",
			"secret":              "",
		},
		"cache_file": map[string]any{
			"enabled": true,
		},
	}
}

func singboxTransportType(network string) string {
	switch network {
	case "h2":
		return "http"
	default:
		return network
	}
}
