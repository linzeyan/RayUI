package model

import "testing"

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
