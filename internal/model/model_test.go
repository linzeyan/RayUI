package model

import (
	"encoding/json"
	"testing"
)

func TestEnumStrings(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ConfigVMess", ConfigVMess.String(), "vmess"},
		{"ConfigVLESS", ConfigVLESS.String(), "vless"},
		{"ConfigTrojan", ConfigTrojan.String(), "trojan"},
		{"ConfigShadowsocks", ConfigShadowsocks.String(), "shadowsocks"},
		{"ConfigHTTP", ConfigHTTP.String(), "http"},
		{"CoreAuto", CoreAuto.String(), "auto"},
		{"CoreXray", CoreXray.String(), "xray"},
		{"CoreSingbox", CoreSingbox.String(), "sing-box"},
		{"ProxyModeManual", ProxyModeManual.String(), "manual"},
		{"ProxyModeSystem", ProxyModeSystem.String(), "system"},
		{"ProxyModeTUN", ProxyModeTUN.String(), "tun"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.UI.Language != "en" {
		t.Errorf("default language = %q, want %q", cfg.UI.Language, "en")
	}
	if cfg.UI.Theme != "system" {
		t.Errorf("default theme = %q, want %q", cfg.UI.Theme, "system")
	}
	if len(cfg.Inbounds) != 2 {
		t.Errorf("expected 2 inbounds, got %d", len(cfg.Inbounds))
	}
	if cfg.Inbounds[0].Port != 10808 {
		t.Errorf("socks port = %d, want 10808", cfg.Inbounds[0].Port)
	}
	if cfg.Inbounds[1].Port != 10809 {
		t.Errorf("http port = %d, want 10809", cfg.Inbounds[1].Port)
	}
	if cfg.TUN.Stack != "gvisor" {
		t.Errorf("tun stack = %q, want %q", cfg.TUN.Stack, "gvisor")
	}
}

func TestDefaultRoutingItems(t *testing.T) {
	items := DefaultRoutingItems()
	if len(items) != 4 {
		t.Fatalf("expected 4 routing items, got %d", len(items))
	}

	names := []string{"Global", "BypassLAN", "BypassCN", "BypassLAN+CN"}
	for i, want := range names {
		if items[i].Remarks != want {
			t.Errorf("item[%d].Remarks = %q, want %q", i, items[i].Remarks, want)
		}
		if !items[i].Locked {
			t.Errorf("item[%d] should be locked", i)
		}
	}
}

func TestDefaultDNSItem(t *testing.T) {
	dns := DefaultDNSItem()
	if dns.RemoteDNS == "" {
		t.Error("RemoteDNS should not be empty")
	}
	if dns.DomainStrategy != "prefer_ipv4" {
		t.Errorf("DomainStrategy = %q, want %q", dns.DomainStrategy, "prefer_ipv4")
	}
}

func TestProfileItemValidate(t *testing.T) {
	p := NewProfileItem()
	if err := p.Validate(); err == nil {
		t.Error("expected validation error for empty address")
	}

	p.Address = "example.com"
	p.Port = 443
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	p.Port = 0
	if err := p.Validate(); err == nil {
		t.Error("expected validation error for port 0")
	}

	p.Port = 70000
	if err := p.Validate(); err == nil {
		t.Error("expected validation error for port > 65535")
	}
}

func TestEnumStringsComplete(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ConfigSOCKS", ConfigSOCKS.String(), "socks"},
		{"ConfigHysteria2", ConfigHysteria2.String(), "hysteria2"},
		{"ConfigTUIC", ConfigTUIC.String(), "tuic"},
		{"ConfigWireGuard", ConfigWireGuard.String(), "wireguard"},
		{"ConfigHTTP", ConfigHTTP.String(), "http"},
		{"ProxyModePAC", ProxyModePAC.String(), "pac"},
		{"UnknownConfigType", EConfigType(99).String(), "unknown"},
		{"UnknownCoreType", ECoreType(99).String(), "unknown"},
		{"UnknownProxyMode", EProxyMode(99).String(), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestNewProfileItem(t *testing.T) {
	p := NewProfileItem()
	if p.ID == "" {
		t.Error("ID should not be empty")
	}
	if p.Network != "tcp" {
		t.Errorf("Network = %q, want tcp", p.Network)
	}
	if p.StreamSecurity != "none" {
		t.Errorf("StreamSecurity = %q, want none", p.StreamSecurity)
	}
	if p.Security != "auto" {
		t.Errorf("Security = %q, want auto", p.Security)
	}
}

func TestNewSubItem(t *testing.T) {
	s := NewSubItem()
	if s.ID == "" {
		t.Error("ID should not be empty")
	}
	if !s.Enabled {
		t.Error("Enabled should be true by default")
	}
}

func TestProfileItemValidateBoundary(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		port    int
		wantErr bool
	}{
		{"port 1 valid", "host", 1, false},
		{"port 65535 valid", "host", 65535, false},
		{"port -1 invalid", "host", -1, true},
		{"port 65536 invalid", "host", 65536, true},
		{"empty address", "", 443, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProfileItem{Address: tt.addr, Port: tt.port}
			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultRoutingItemsRuleIDs(t *testing.T) {
	items := DefaultRoutingItems()
	for i, item := range items {
		if item.ID == "" {
			t.Errorf("item[%d] has empty ID", i)
		}
		for j, rule := range item.Rules {
			if rule.ID == "" {
				t.Errorf("item[%d].Rules[%d] has empty ID", i, j)
			}
		}
	}
}

func TestProfileItemJSONRoundTrip(t *testing.T) {
	mux := true
	p := ProfileItem{
		ID: "test-id", ConfigType: ConfigVLESS, Remarks: "Test",
		Address: "1.2.3.4", Port: 443, UUID: "uuid-test",
		Security: "none", Network: "ws", Host: "example.com",
		Path: "/ws", StreamSecurity: "tls", SNI: "example.com",
		ALPN: "h2", Fingerprint: "chrome", PublicKey: "pk",
		ShortID: "sid", SpiderX: "/", CoreType: CoreSingbox,
		Extra: "extra", MuxEnabled: &mux, AllowInsecure: true,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got ProfileItem
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.ID != p.ID || got.Address != p.Address || got.Port != p.Port {
		t.Errorf("basic fields mismatch: got %+v", got)
	}
	if got.ConfigType != p.ConfigType || got.CoreType != p.CoreType {
		t.Errorf("enum fields mismatch")
	}
	if got.MuxEnabled == nil || *got.MuxEnabled != true {
		t.Errorf("MuxEnabled mismatch")
	}
	if got.PublicKey != p.PublicKey || got.ShortID != p.ShortID {
		t.Errorf("reality fields mismatch")
	}
}
