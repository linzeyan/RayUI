package config

import (
	"encoding/json"
	"testing"

	"github.com/RayUI/RayUI/internal/model"
)

func TestSingboxGenerateVLESSReality(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigVLESS,
		Address:        "1.2.3.4",
		Port:           443,
		UUID:           "test-uuid",
		Network:        "tcp",
		StreamSecurity: "reality",
		SNI:            "www.google.com",
		Fingerprint:    "chrome",
		PublicKey:      "pubkey123",
		ShortID:        "ab12",
	}
	routing := model.DefaultRoutingItems()[0] // Global
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	// Check outbounds.
	outbounds, ok := result["outbounds"].([]any)
	if !ok || len(outbounds) < 1 {
		t.Fatal("missing outbounds")
	}
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "vless" {
		t.Errorf("type = %v", proxy["type"])
	}
	if proxy["uuid"] != "test-uuid" {
		t.Errorf("uuid = %v", proxy["uuid"])
	}

	// Check TLS/Reality.
	tls, ok := proxy["tls"].(map[string]any)
	if !ok {
		t.Fatal("missing tls")
	}
	reality, ok := tls["reality"].(map[string]any)
	if !ok {
		t.Fatal("missing reality")
	}
	if reality["public_key"] != "pubkey123" {
		t.Errorf("public_key = %v", reality["public_key"])
	}
}

func TestSingboxGenerateVMessWS(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigVMess,
		Address:        "example.com",
		Port:           443,
		UUID:           "vmess-uuid",
		Security:       "auto",
		Network:        "ws",
		Host:           "example.com",
		Path:           "/ws",
		StreamSecurity: "tls",
		SNI:            "example.com",
	}
	routing := model.DefaultRoutingItems()[1] // BypassLAN
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "vmess" {
		t.Errorf("type = %v", proxy["type"])
	}

	transport, ok := proxy["transport"].(map[string]any)
	if !ok {
		t.Fatal("missing transport")
	}
	if transport["type"] != "ws" {
		t.Errorf("transport type = %v", transport["type"])
	}
}

func TestXrayGenerateVLESSReality(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigVLESS,
		Address:        "1.2.3.4",
		Port:           443,
		UUID:           "test-uuid",
		Network:        "tcp",
		StreamSecurity: "reality",
		SNI:            "www.google.com",
		Fingerprint:    "chrome",
		PublicKey:      "pubkey123",
		ShortID:        "ab12",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["protocol"] != "vless" {
		t.Errorf("protocol = %v", proxy["protocol"])
	}

	ss := proxy["streamSettings"].(map[string]any)
	if ss["security"] != "reality" {
		t.Errorf("security = %v", ss["security"])
	}
	reality := ss["realitySettings"].(map[string]any)
	if reality["publicKey"] != "pubkey123" {
		t.Errorf("publicKey = %v", reality["publicKey"])
	}
}

func TestSingboxGenerateHysteria2(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigHysteria2,
		Address:        "hy2.example.com",
		Port:           443,
		UUID:           "myauth",
		HeaderType:     "salamander",
		Path:           "obfs-secret",
		StreamSecurity: "tls",
		SNI:            "hy2.example.com",
		ALPN:           "h3",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "hysteria2" {
		t.Errorf("type = %v", proxy["type"])
	}
	if proxy["password"] != "myauth" {
		t.Errorf("password = %v", proxy["password"])
	}
	obfs, ok := proxy["obfs"].(map[string]any)
	if !ok {
		t.Fatal("missing obfs")
	}
	if obfs["type"] != "salamander" {
		t.Errorf("obfs type = %v", obfs["type"])
	}
	if obfs["password"] != "obfs-secret" {
		t.Errorf("obfs password = %v", obfs["password"])
	}
	// Should have TLS.
	if _, ok := proxy["tls"]; !ok {
		t.Error("missing tls for hysteria2")
	}
}

func TestSingboxGenerateTUIC(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigTUIC,
		Address:        "tuic.example.com",
		Port:           443,
		UUID:           "tuic-uuid",
		Security:       "tuic-pass",
		HeaderType:     "bbr",
		Path:           "native",
		StreamSecurity: "tls",
		SNI:            "tuic.example.com",
		ALPN:           "h3",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "tuic" {
		t.Errorf("type = %v", proxy["type"])
	}
	if proxy["uuid"] != "tuic-uuid" {
		t.Errorf("uuid = %v", proxy["uuid"])
	}
	if proxy["password"] != "tuic-pass" {
		t.Errorf("password = %v", proxy["password"])
	}
	if proxy["congestion_control"] != "bbr" {
		t.Errorf("congestion_control = %v", proxy["congestion_control"])
	}
	if proxy["udp_relay_mode"] != "native" {
		t.Errorf("udp_relay_mode = %v", proxy["udp_relay_mode"])
	}
}

func TestSingboxGenerateWireGuard(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigWireGuard,
		Address:    "wg.example.com",
		Port:       51820,
		UUID:       "privkey123",
		PublicKey:  "pubkey456",
		Host:       "10.0.0.2/32",
		ShortID:    "1,2,3",
		Extra:      "1420",
		Network:    "wireguard",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "wireguard" {
		t.Errorf("type = %v", proxy["type"])
	}
	if proxy["private_key"] != "privkey123" {
		t.Errorf("private_key = %v", proxy["private_key"])
	}
	if proxy["peer_public_key"] != "pubkey456" {
		t.Errorf("peer_public_key = %v", proxy["peer_public_key"])
	}
	// Should NOT have tls or transport (WireGuard returns early).
	if _, ok := proxy["tls"]; ok {
		t.Error("wireguard should not have tls")
	}
	if _, ok := proxy["transport"]; ok {
		t.Error("wireguard should not have transport")
	}
	// Check reserved.
	reserved, ok := proxy["reserved"].([]any)
	if !ok {
		t.Fatal("missing reserved")
	}
	if len(reserved) != 3 {
		t.Errorf("reserved length = %d", len(reserved))
	}
	// Check MTU.
	mtu, ok := proxy["mtu"].(float64)
	if !ok || int(mtu) != 1420 {
		t.Errorf("mtu = %v", proxy["mtu"])
	}
}

func TestXrayRoutingRules(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigVMess,
		Address:    "1.2.3.4",
		Port:       443,
		UUID:       "test",
		Security:   "auto",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[3] // BypassLAN+CN
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	routingSection := result["routing"].(map[string]any)
	rules := routingSection["rules"].([]any)
	// Should have: api rule + 5 routing rules = 6
	if len(rules) < 5 {
		t.Errorf("expected at least 5 rules, got %d", len(rules))
	}
}

func TestXrayGenerateTrojan(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigTrojan,
		Address:        "trojan.example.com",
		Port:           443,
		UUID:           "trojan-password",
		Network:        "tcp",
		StreamSecurity: "tls",
		SNI:            "trojan.example.com",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["protocol"] != "trojan" {
		t.Errorf("protocol = %v, want trojan", proxy["protocol"])
	}

	settings := proxy["settings"].(map[string]any)
	servers := settings["servers"].([]any)
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
	server := servers[0].(map[string]any)
	if server["password"] != "trojan-password" {
		t.Errorf("password = %v", server["password"])
	}
	if server["address"] != "trojan.example.com" {
		t.Errorf("address = %v", server["address"])
	}

	ss := proxy["streamSettings"].(map[string]any)
	if ss["security"] != "tls" {
		t.Errorf("security = %v", ss["security"])
	}
}

func TestXrayGenerateShadowsocks(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigShadowsocks,
		Address:    "ss.example.com",
		Port:       8388,
		UUID:       "ss-password",
		Security:   "aes-256-gcm",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["protocol"] != "shadowsocks" {
		t.Errorf("protocol = %v, want shadowsocks", proxy["protocol"])
	}

	settings := proxy["settings"].(map[string]any)
	servers := settings["servers"].([]any)
	server := servers[0].(map[string]any)
	if server["method"] != "aes-256-gcm" {
		t.Errorf("method = %v", server["method"])
	}
	if server["password"] != "ss-password" {
		t.Errorf("password = %v", server["password"])
	}
}

func TestSingboxGenerateTrojan(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType:     model.ConfigTrojan,
		Address:        "trojan.example.com",
		Port:           443,
		UUID:           "trojan-pass",
		Network:        "ws",
		Host:           "trojan.example.com",
		Path:           "/ws",
		StreamSecurity: "tls",
		SNI:            "trojan.example.com",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "trojan" {
		t.Errorf("type = %v, want trojan", proxy["type"])
	}
	if proxy["password"] != "trojan-pass" {
		t.Errorf("password = %v", proxy["password"])
	}

	transport := proxy["transport"].(map[string]any)
	if transport["type"] != "ws" {
		t.Errorf("transport type = %v", transport["type"])
	}
}

func TestSingboxGenerateShadowsocks(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigShadowsocks,
		Address:    "ss.example.com",
		Port:       8388,
		UUID:       "ss-pass",
		Security:   "chacha20-ietf-poly1305",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	proxy := outbounds[0].(map[string]any)
	if proxy["type"] != "shadowsocks" {
		t.Errorf("type = %v, want shadowsocks", proxy["type"])
	}
	if proxy["method"] != "chacha20-ietf-poly1305" {
		t.Errorf("method = %v", proxy["method"])
	}
	if proxy["password"] != "ss-pass" {
		t.Errorf("password = %v", proxy["password"])
	}
}

func TestXrayDNSConfig(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigVMess,
		Address:    "1.2.3.4",
		Port:       443,
		UUID:       "test",
		Security:   "auto",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DNSItem{
		RemoteDNS:    "https://dns.google/dns-query",
		DirectDNS:    "223.5.5.5",
		BootstrapDNS: "8.8.8.8",
	}
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	dnsSection := result["dns"].(map[string]any)
	servers := dnsSection["servers"].([]any)
	if len(servers) < 3 {
		t.Errorf("expected at least 3 DNS servers, got %d", len(servers))
	}
}

func TestSingboxDNSFakeIP(t *testing.T) {
	g := &SingboxConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigVLESS,
		Address:    "1.2.3.4",
		Port:       443,
		UUID:       "test",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	dns.FakeIP = true
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	dnsSection := result["dns"].(map[string]any)
	fakeip, ok := dnsSection["fakeip"].(map[string]any)
	if !ok {
		t.Fatal("missing fakeip section")
	}
	if fakeip["enabled"] != true {
		t.Errorf("fakeip.enabled = %v", fakeip["enabled"])
	}
}

func TestXrayOutboundsStructure(t *testing.T) {
	g := &XrayConfigGenerator{}
	profile := model.ProfileItem{
		ConfigType: model.ConfigVMess,
		Address:    "1.2.3.4",
		Port:       443,
		UUID:       "test",
		Security:   "auto",
		Network:    "tcp",
	}
	routing := model.DefaultRoutingItems()[0]
	dns := model.DefaultDNSItem()
	cfg := model.DefaultConfig()

	data, err := g.Generate(profile, routing, dns, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	outbounds := result["outbounds"].([]any)
	if len(outbounds) != 3 {
		t.Fatalf("expected 3 outbounds (proxy, direct, block), got %d", len(outbounds))
	}

	// Check tags.
	tags := make(map[string]bool)
	for _, ob := range outbounds {
		m := ob.(map[string]any)
		tags[m["tag"].(string)] = true
	}
	for _, want := range []string{"proxy", "direct", "block"} {
		if !tags[want] {
			t.Errorf("missing outbound tag %q", want)
		}
	}
}
