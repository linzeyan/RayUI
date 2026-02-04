package config

import (
	"encoding/json"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// XrayConfigGenerator produces an xray-core JSON config.
type XrayConfigGenerator struct{}

func (g *XrayConfigGenerator) Generate(
	profile model.ProfileItem,
	routing model.RoutingItem,
	dns model.DNSItem,
	cfg model.Config,
) ([]byte, error) {
	root := map[string]any{
		"log":       buildXrayLog(cfg),
		"dns":       buildXrayDNS(dns),
		"inbounds":  buildXrayInbounds(cfg),
		"outbounds": buildXrayOutbounds(profile),
		"routing":   buildXrayRouting(routing),
		"stats":     map[string]any{},
		"api": map[string]any{
			"tag":      "api",
			"services": []string{"StatsService"},
		},
	}
	return json.MarshalIndent(root, "", "  ")
}

func buildXrayLog(cfg model.Config) map[string]any {
	level := cfg.CoreBasic.LogLevel
	if level == "" {
		level = "info"
	}
	return map[string]any{"loglevel": level}
}

func buildXrayDNS(dns model.DNSItem) map[string]any {
	var servers []any
	if dns.RemoteDNS != "" {
		servers = append(servers, map[string]any{
			"address": dns.RemoteDNS,
			"domains": []string{"geosite:geolocation-!cn"},
		})
	}
	if dns.DirectDNS != "" {
		servers = append(servers, map[string]any{
			"address": dns.DirectDNS,
			"domains": []string{"geosite:cn"},
		})
	}
	if dns.BootstrapDNS != "" {
		servers = append(servers, dns.BootstrapDNS)
	}
	return map[string]any{"servers": servers}
}

func buildXrayInbounds(cfg model.Config) []map[string]any {
	var inbounds []map[string]any

	// API inbound for stats.
	inbounds = append(inbounds, map[string]any{
		"tag":      "api-in",
		"protocol": "dokodemo-door",
		"listen":   "127.0.0.1",
		"port":     10813,
		"settings": map[string]any{"address": "127.0.0.1"},
	})

	for _, ib := range cfg.Inbounds {
		listen := ib.ListenAddr
		if listen == "" {
			listen = "127.0.0.1"
		}
		if ib.AllowLAN {
			listen = "0.0.0.0"
		}

		inbound := map[string]any{
			"tag":      ib.Protocol + "-in",
			"protocol": ib.Protocol,
			"listen":   listen,
			"port":     ib.Port,
		}

		settings := map[string]any{}
		if ib.Protocol == "socks" {
			settings["udp"] = ib.UDPEnabled
		}
		inbound["settings"] = settings

		if ib.SniffingEnabled {
			inbound["sniffing"] = map[string]any{
				"enabled":      true,
				"destOverride": []string{"http", "tls"},
			}
		}

		inbounds = append(inbounds, inbound)
	}

	return inbounds
}

func buildXrayOutbounds(profile model.ProfileItem) []map[string]any {
	proxy := buildXrayProxyOutbound(profile)
	direct := map[string]any{"protocol": "freedom", "tag": "direct"}
	block := map[string]any{"protocol": "blackhole", "tag": "block", "settings": map[string]any{"response": map[string]any{"type": "none"}}}
	return []map[string]any{proxy, direct, block}
}

func buildXrayProxyOutbound(p model.ProfileItem) map[string]any {
	ob := map[string]any{
		"tag": "proxy",
	}

	settings := map[string]any{}
	switch p.ConfigType {
	case model.ConfigVMess:
		ob["protocol"] = "vmess"
		settings["vnext"] = []map[string]any{{
			"address": p.Address,
			"port":    p.Port,
			"users": []map[string]any{{
				"id":       p.UUID,
				"alterId":  p.AlterID,
				"security": nonEmpty(p.Security, "auto"),
			}},
		}}

	case model.ConfigVLESS:
		ob["protocol"] = "vless"
		user := map[string]any{
			"id":         p.UUID,
			"encryption": "none",
		}
		if p.Flow != "" {
			user["flow"] = p.Flow
		}
		settings["vnext"] = []map[string]any{{
			"address": p.Address,
			"port":    p.Port,
			"users":   []map[string]any{user},
		}}

	case model.ConfigTrojan:
		ob["protocol"] = "trojan"
		settings["servers"] = []map[string]any{{
			"address":  p.Address,
			"port":     p.Port,
			"password": p.UUID,
		}}

	case model.ConfigShadowsocks:
		ob["protocol"] = "shadowsocks"
		settings["servers"] = []map[string]any{{
			"address":  p.Address,
			"port":     p.Port,
			"method":   p.Security,
			"password": p.UUID,
		}}
	}

	ob["settings"] = settings
	ob["streamSettings"] = buildXrayStreamSettings(p)

	return ob
}

func buildXrayStreamSettings(p model.ProfileItem) map[string]any {
	ss := map[string]any{
		"network": nonEmpty(p.Network, "tcp"),
	}

	// Transport settings.
	switch p.Network {
	case "ws":
		ws := map[string]any{}
		if p.Path != "" {
			ws["path"] = p.Path
		}
		if p.Host != "" {
			ws["headers"] = map[string]any{"Host": p.Host}
		}
		ss["wsSettings"] = ws

	case "h2":
		h2 := map[string]any{}
		if p.Host != "" {
			h2["host"] = strings.Split(p.Host, ",")
		}
		if p.Path != "" {
			h2["path"] = p.Path
		}
		ss["httpSettings"] = h2

	case "grpc":
		grpc := map[string]any{}
		if p.Path != "" {
			grpc["serviceName"] = p.Path
		}
		ss["grpcSettings"] = grpc

	case "kcp":
		kcp := map[string]any{}
		if p.HeaderType != "" {
			kcp["header"] = map[string]any{"type": p.HeaderType}
		}
		ss["kcpSettings"] = kcp

	case "tcp":
		if p.HeaderType == "http" {
			ss["tcpSettings"] = map[string]any{
				"header": map[string]any{
					"type": "http",
					"request": map[string]any{
						"path":    []string{nonEmpty(p.Path, "/")},
						"headers": map[string]any{"Host": strings.Split(nonEmpty(p.Host, ""), ",")},
					},
				},
			}
		}

	case "httpupgrade":
		hu := map[string]any{}
		if p.Host != "" {
			hu["host"] = p.Host
		}
		if p.Path != "" {
			hu["path"] = p.Path
		}
		ss["httpupgradeSettings"] = hu
	}

	// Security settings.
	switch p.StreamSecurity {
	case "tls":
		ss["security"] = "tls"
		tls := map[string]any{}
		if p.SNI != "" {
			tls["serverName"] = p.SNI
		}
		if p.AllowInsecure {
			tls["allowInsecure"] = true
		}
		if p.ALPN != "" {
			tls["alpn"] = strings.Split(p.ALPN, ",")
		}
		if p.Fingerprint != "" {
			tls["fingerprint"] = p.Fingerprint
		}
		ss["tlsSettings"] = tls

	case "reality":
		ss["security"] = "reality"
		reality := map[string]any{}
		if p.SNI != "" {
			reality["serverName"] = p.SNI
		}
		if p.Fingerprint != "" {
			reality["fingerprint"] = p.Fingerprint
		}
		if p.PublicKey != "" {
			reality["publicKey"] = p.PublicKey
		}
		if p.ShortID != "" {
			reality["shortId"] = p.ShortID
		}
		if p.SpiderX != "" {
			reality["spiderX"] = p.SpiderX
		}
		ss["realitySettings"] = reality

	default:
		ss["security"] = "none"
	}

	return ss
}

func buildXrayRouting(routing model.RoutingItem) map[string]any {
	strategy := routing.DomainStrategy
	if strategy == "" {
		strategy = "AsIs"
	}

	var rules []map[string]any

	// API rule.
	rules = append(rules, map[string]any{
		"inboundTag":  []string{"api-in"},
		"outboundTag": "api",
		"type":        "field",
	})

	for _, r := range routing.Rules {
		if !r.Enabled {
			continue
		}
		rule := map[string]any{
			"outboundTag": r.OutboundTag,
			"type":        "field",
		}
		if len(r.Domain) > 0 {
			rule["domain"] = r.Domain
		}
		if len(r.DomainSuffix) > 0 {
			for _, ds := range r.DomainSuffix {
				rule["domain"] = append(toStringSlice(rule["domain"]), "domain:"+ds)
			}
		}
		if len(r.DomainKeyword) > 0 {
			for _, dk := range r.DomainKeyword {
				rule["domain"] = append(toStringSlice(rule["domain"]), "keyword:"+dk)
			}
		}
		if len(r.Geosite) > 0 {
			for _, gs := range r.Geosite {
				rule["domain"] = append(toStringSlice(rule["domain"]), "geosite:"+gs)
			}
		}
		if len(r.GeoIP) > 0 {
			var ips []string
			for _, gi := range r.GeoIP {
				ips = append(ips, "geoip:"+gi)
			}
			rule["ip"] = ips
		}
		if len(r.IPCIDR) > 0 {
			rule["ip"] = append(toStringSlice(rule["ip"]), r.IPCIDR...)
		}
		if r.Port != "" {
			rule["port"] = r.Port
		}
		if len(r.Protocol) > 0 {
			rule["protocol"] = r.Protocol
		}
		if r.Network != "" {
			rule["network"] = r.Network
		}

		rules = append(rules, rule)
	}

	return map[string]any{
		"domainStrategy": strategy,
		"rules":          rules,
	}
}

func nonEmpty(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	if s, ok := v.([]string); ok {
		return s
	}
	return nil
}
