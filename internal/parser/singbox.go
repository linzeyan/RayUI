package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/util"
)

// singboxConfig represents relevant parts of a sing-box config.
type singboxConfig struct {
	Outbounds []singboxOutbound `json:"outbounds"`
}

type singboxOutbound struct {
	Type       string `json:"type"`
	Tag        string `json:"tag"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`

	// VMess / VLESS
	UUID     string `json:"uuid,omitempty"`
	AlterID  int    `json:"alter_id,omitempty"`
	Security string `json:"security,omitempty"`
	Flow     string `json:"flow,omitempty"`

	// Shadowsocks
	Method   string `json:"method,omitempty"`
	Password string `json:"password,omitempty"`

	// Transport
	Transport *singboxTransport `json:"transport,omitempty"`

	// TLS
	TLS *singboxTLS `json:"tls,omitempty"`
}

type singboxTransport struct {
	Type        string            `json:"type"`
	Host        string            `json:"host,omitempty"`
	Path        string            `json:"path,omitempty"`
	ServiceName string            `json:"service_name,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

type singboxTLS struct {
	Enabled    bool     `json:"enabled"`
	ServerName string   `json:"server_name,omitempty"`
	Insecure   bool     `json:"insecure,omitempty"`
	ALPN       []string `json:"alpn,omitempty"`

	// Reality
	Reality *singboxReality `json:"reality,omitempty"`

	// UTLS
	UTLS *singboxUTLS `json:"utls,omitempty"`
}

type singboxReality struct {
	Enabled   bool   `json:"enabled"`
	PublicKey string `json:"public_key"`
	ShortID   string `json:"short_id"`
}

type singboxUTLS struct {
	Fingerprint string `json:"fingerprint"`
}

func parseSingboxOutbounds(content string) ([]model.ProfileItem, error) {
	var cfg singboxConfig
	if err := json.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("singbox config parse: %w", err)
	}

	var items []model.ProfileItem
	for _, ob := range cfg.Outbounds {
		p, ok := singboxOutboundToProfile(ob)
		if !ok {
			continue
		}
		items = append(items, p)
	}
	return items, nil
}

func singboxOutboundToProfile(ob singboxOutbound) (model.ProfileItem, bool) {
	p := model.ProfileItem{
		ID:             util.GenerateUUID(),
		Address:        ob.Server,
		Port:           ob.ServerPort,
		Remarks:        ob.Tag,
		Network:        "tcp",
		StreamSecurity: "none",
	}

	switch ob.Type {
	case "vmess":
		p.ConfigType = model.ConfigVMess
		p.UUID = ob.UUID
		p.AlterID = ob.AlterID
		p.Security = ob.Security
		if p.Security == "" {
			p.Security = "auto"
		}
	case "vless":
		p.ConfigType = model.ConfigVLESS
		p.UUID = ob.UUID
		p.Security = "none"
		p.Flow = ob.Flow
	case "trojan":
		p.ConfigType = model.ConfigTrojan
		p.UUID = ob.Password
	case "shadowsocks":
		p.ConfigType = model.ConfigShadowsocks
		p.Security = ob.Method
		p.UUID = ob.Password
	default:
		return p, false
	}

	// Transport
	if ob.Transport != nil {
		switch ob.Transport.Type {
		case "ws":
			p.Network = "ws"
			p.Host = ob.Transport.Host
			p.Path = ob.Transport.Path
		case "http":
			p.Network = "h2"
			p.Host = ob.Transport.Host
			p.Path = ob.Transport.Path
		case "grpc":
			p.Network = "grpc"
			p.Path = ob.Transport.ServiceName
		case "httpupgrade":
			p.Network = "httpupgrade"
			p.Host = ob.Transport.Host
			p.Path = ob.Transport.Path
		}
	}

	// TLS
	if ob.TLS != nil && ob.TLS.Enabled {
		p.StreamSecurity = "tls"
		p.SNI = ob.TLS.ServerName
		p.AllowInsecure = ob.TLS.Insecure
		if len(ob.TLS.ALPN) > 0 {
			p.ALPN = strings.Join(ob.TLS.ALPN, ",")
		}
		if ob.TLS.UTLS != nil {
			p.Fingerprint = ob.TLS.UTLS.Fingerprint
		}
		if ob.TLS.Reality != nil && ob.TLS.Reality.Enabled {
			p.StreamSecurity = "reality"
			p.PublicKey = ob.TLS.Reality.PublicKey
			p.ShortID = ob.TLS.Reality.ShortID
		}
	}

	return p, true
}
