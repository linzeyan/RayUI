package parser

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/RayUI/RayUI/internal/model"
)

// --- VMess tests ---

func TestParseVMess(t *testing.T) {
	// Standard VMess URI.
	vmessJSON := `{"v":"2","ps":"Tokyo","add":"1.2.3.4","port":"443","id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","aid":"0","scy":"auto","net":"ws","type":"none","host":"example.com","path":"/ws","tls":"tls","sni":"example.com","alpn":"h2,http/1.1","fp":"chrome"}`
	uri := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}
	if p.ConfigType != model.ConfigVMess {
		t.Errorf("configType = %v, want VMess", p.ConfigType)
	}
	if p.Remarks != "Tokyo" {
		t.Errorf("remarks = %q, want %q", p.Remarks, "Tokyo")
	}
	if p.Address != "1.2.3.4" {
		t.Errorf("address = %q", p.Address)
	}
	if p.Port != 443 {
		t.Errorf("port = %d, want 443", p.Port)
	}
	if p.Network != "ws" {
		t.Errorf("network = %q, want ws", p.Network)
	}
	if p.StreamSecurity != "tls" {
		t.Errorf("streamSecurity = %q, want tls", p.StreamSecurity)
	}
	if p.Fingerprint != "chrome" {
		t.Errorf("fingerprint = %q, want chrome", p.Fingerprint)
	}
}

func TestVMessRoundTrip(t *testing.T) {
	vmessJSON := `{"v":"2","ps":"HK","add":"10.0.0.1","port":"8443","id":"11111111-2222-3333-4444-555555555555","aid":"0","scy":"auto","net":"tcp","type":"none","host":"","path":"","tls":"","sni":"","alpn":"","fp":""}`
	uri := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}

	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}

	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}

	if p.Address != p2.Address || p.Port != p2.Port || p.UUID != p2.UUID || p.Network != p2.Network {
		t.Errorf("round-trip mismatch: %+v vs %+v", p, p2)
	}
}

// --- VLESS tests ---

func TestParseVLESS(t *testing.T) {
	uri := "vless://aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee@example.com:443?encryption=none&security=tls&sni=example.com&type=ws&host=example.com&path=%2Fws&fp=chrome&alpn=h2#Tokyo"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}
	if p.ConfigType != model.ConfigVLESS {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.UUID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("uuid = %q", p.UUID)
	}
	if p.Address != "example.com" {
		t.Errorf("address = %q", p.Address)
	}
	if p.Port != 443 {
		t.Errorf("port = %d", p.Port)
	}
	if p.Network != "ws" {
		t.Errorf("network = %q", p.Network)
	}
	if p.Remarks != "Tokyo" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestParseVLESSReality(t *testing.T) {
	uri := "vless://uuid@1.2.3.4:443?encryption=none&security=reality&sni=www.google.com&fp=chrome&pbk=pubkey123&sid=ab12&spx=%2F&type=tcp#RealityNode"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.StreamSecurity != "reality" {
		t.Errorf("streamSecurity = %q, want reality", p.StreamSecurity)
	}
	if p.PublicKey != "pubkey123" {
		t.Errorf("publicKey = %q", p.PublicKey)
	}
	if p.ShortID != "ab12" {
		t.Errorf("shortId = %q", p.ShortID)
	}
}

func TestVLESSRoundTrip(t *testing.T) {
	uri := "vless://test-uuid@example.com:443?security=tls&sni=example.com&type=grpc&fp=chrome#GRPC"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != p2.UUID || p.Address != p2.Address || p.Port != p2.Port {
		t.Errorf("round-trip mismatch")
	}
}

// --- Trojan tests ---

func TestParseTrojan(t *testing.T) {
	uri := "trojan://mypassword@example.com:443?security=tls&sni=example.com&type=ws&host=example.com&path=%2Fws&fp=chrome#TrojanWS"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigTrojan {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.UUID != "mypassword" {
		t.Errorf("password = %q", p.UUID)
	}
	if p.Network != "ws" {
		t.Errorf("network = %q", p.Network)
	}
	if p.Remarks != "TrojanWS" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestTrojanRoundTrip(t *testing.T) {
	uri := "trojan://pass123@server.com:443?sni=server.com#Node1"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != p2.UUID || p.Address != p2.Address || p.Port != p2.Port {
		t.Errorf("round-trip mismatch")
	}
}

// --- Shadowsocks tests ---

func TestParseShadowsocksFormat1(t *testing.T) {
	// Format 1: ss://base64(method:password)@host:port#remarks
	userInfo := base64.URLEncoding.EncodeToString([]byte("aes-256-gcm:mypass"))
	userInfo = strings.TrimRight(userInfo, "=")
	uri := "ss://" + userInfo + "@1.2.3.4:8388#TestSS"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigShadowsocks {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.Security != "aes-256-gcm" {
		t.Errorf("method = %q", p.Security)
	}
	if p.UUID != "mypass" {
		t.Errorf("password = %q", p.UUID)
	}
	if p.Address != "1.2.3.4" || p.Port != 8388 {
		t.Errorf("host:port = %s:%d", p.Address, p.Port)
	}
	if p.Remarks != "TestSS" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestParseShadowsocksFormat2(t *testing.T) {
	// Format 2: ss://base64(method:password@host:port)#remarks
	full := base64.URLEncoding.EncodeToString([]byte("chacha20-ietf-poly1305:secret@5.6.7.8:9999"))
	full = strings.TrimRight(full, "=")
	uri := "ss://" + full + "#Format2"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.Security != "chacha20-ietf-poly1305" {
		t.Errorf("method = %q", p.Security)
	}
	if p.Address != "5.6.7.8" || p.Port != 9999 {
		t.Errorf("host:port = %s:%d", p.Address, p.Port)
	}
}

func TestShadowsocksRoundTrip(t *testing.T) {
	userInfo := base64.URLEncoding.EncodeToString([]byte("aes-128-gcm:pass123"))
	userInfo = strings.TrimRight(userInfo, "=")
	uri := "ss://" + userInfo + "@server.com:1234#SS-Node"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.Security != p2.Security || p.UUID != p2.UUID || p.Address != p2.Address || p.Port != p2.Port {
		t.Errorf("round-trip mismatch: security=%q/%q addr=%s/%s", p.Security, p2.Security, p.Address, p2.Address)
	}
}

// --- Batch / Format detection tests ---

func TestParseBatchBase64(t *testing.T) {
	lines := []string{
		"vless://uuid@1.2.3.4:443?security=tls&sni=a.com#Node1",
		"trojan://pass@5.6.7.8:443?sni=b.com#Node2",
	}
	raw := strings.Join(lines, "\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	items, err := ParseBatch(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2, got %d", len(items))
	}
	if items[0].ConfigType != model.ConfigVLESS {
		t.Errorf("item[0] type = %v", items[0].ConfigType)
	}
	if items[1].ConfigType != model.ConfigTrojan {
		t.Errorf("item[1] type = %v", items[1].ConfigType)
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"uri-lines", "vless://uuid@host:443#test\ntrojan://pass@host:443#test2", "uri-lines"},
		{"singbox", `{"outbounds":[{"type":"vless"}]}`, "json-singbox"},
		{"sip008", `{"servers":[{"server":"1.2.3.4"}]}`, "json-sip008"},
		{"base64", base64.StdEncoding.EncodeToString([]byte("vless://uuid@host:443#test")), "base64"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFormat(tt.input)
			if got != tt.want {
				t.Errorf("DetectFormat = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseBatchSingbox(t *testing.T) {
	content := `{
		"outbounds": [
			{"type": "vmess", "tag": "vmess-node", "server": "1.2.3.4", "server_port": 443, "uuid": "test-uuid", "security": "auto"},
			{"type": "selector", "tag": "select"},
			{"type": "direct", "tag": "direct"}
		]
	}`
	items, err := ParseBatch(content)
	if err != nil {
		t.Fatal(err)
	}
	// selector and direct are skipped.
	if len(items) != 1 {
		t.Fatalf("expected 1, got %d", len(items))
	}
	if items[0].ConfigType != model.ConfigVMess {
		t.Errorf("type = %v", items[0].ConfigType)
	}
}

func TestParseURIUnsupported(t *testing.T) {
	_, err := ParseURI("http://example.com")
	if err == nil {
		t.Error("expected error for unsupported scheme")
	}
}

// --- Hysteria2 tests ---

func TestParseHysteria2(t *testing.T) {
	uri := "hysteria2://myauth@example.com:443?sni=example.com&obfs=salamander&obfs-password=secret&alpn=h3&insecure=1#Hy2Node"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigHysteria2 {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.UUID != "myauth" {
		t.Errorf("auth = %q, want myauth", p.UUID)
	}
	if p.Address != "example.com" || p.Port != 443 {
		t.Errorf("host:port = %s:%d", p.Address, p.Port)
	}
	if p.HeaderType != "salamander" {
		t.Errorf("obfs = %q", p.HeaderType)
	}
	if p.Path != "secret" {
		t.Errorf("obfs-password = %q", p.Path)
	}
	if !p.AllowInsecure {
		t.Error("allowInsecure should be true")
	}
	if p.Remarks != "Hy2Node" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestParseHy2Alias(t *testing.T) {
	uri := "hy2://pass@host:8443#AliasTest"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigHysteria2 {
		t.Errorf("configType = %v", p.ConfigType)
	}
}

func TestHysteria2RoundTrip(t *testing.T) {
	uri := "hysteria2://auth123@server.com:443?sni=server.com&obfs=salamander&obfs-password=pwd&alpn=h3#HyNode"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != p2.UUID || p.Address != p2.Address || p.Port != p2.Port || p.HeaderType != p2.HeaderType {
		t.Errorf("round-trip mismatch")
	}
}

// --- TUIC tests ---

func TestParseTUIC(t *testing.T) {
	uri := "tuic://aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee:mypassword@example.com:443?congestion_control=bbr&udp_relay_mode=native&alpn=h3&sni=example.com#TuicNode"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigTUIC {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.UUID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("uuid = %q", p.UUID)
	}
	if p.Security != "mypassword" {
		t.Errorf("password = %q", p.Security)
	}
	if p.HeaderType != "bbr" {
		t.Errorf("congestion_control = %q", p.HeaderType)
	}
	if p.Path != "native" {
		t.Errorf("udp_relay_mode = %q", p.Path)
	}
	if p.Remarks != "TuicNode" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestTUICRoundTrip(t *testing.T) {
	uri := "tuic://uuid123:pass456@server.com:443?congestion_control=bbr&udp_relay_mode=native&alpn=h3&sni=server.com#TUIC"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != p2.UUID || p.Security != p2.Security || p.Address != p2.Address || p.Port != p2.Port {
		t.Errorf("round-trip mismatch")
	}
}

// --- WireGuard tests ---

func TestParseWireGuard(t *testing.T) {
	uri := "wireguard://privatekey123@example.com:51820?publickey=pubkey456&address=10.0.0.2/32&reserved=1,2,3&mtu=1420#WGNode"

	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigWireGuard {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.UUID != "privatekey123" {
		t.Errorf("privateKey = %q", p.UUID)
	}
	if p.PublicKey != "pubkey456" {
		t.Errorf("publicKey = %q", p.PublicKey)
	}
	if p.Host != "10.0.0.2/32" {
		t.Errorf("address = %q", p.Host)
	}
	if p.ShortID != "1,2,3" {
		t.Errorf("reserved = %q", p.ShortID)
	}
	if p.Extra != "1420" {
		t.Errorf("mtu = %q", p.Extra)
	}
	if p.Remarks != "WGNode" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestParseWGAlias(t *testing.T) {
	uri := "wg://key@host:51820?publickey=pk#Alias"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.ConfigType != model.ConfigWireGuard {
		t.Errorf("configType = %v", p.ConfigType)
	}
}

func TestWireGuardRoundTrip(t *testing.T) {
	uri := "wireguard://privkey@server.com:51820?publickey=pubkey&address=10.0.0.2/32&reserved=0,0,0&mtu=1280#WG"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := ToShareURI(*p)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ParseURI(generated)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != p2.UUID || p.PublicKey != p2.PublicKey || p.Address != p2.Address || p.Port != p2.Port {
		t.Errorf("round-trip mismatch")
	}
}

// --- Clash YAML tests ---

func TestParseClashYAML(t *testing.T) {
	content := `
proxies:
  - name: "VMess-WS"
    type: vmess
    server: 1.2.3.4
    port: 443
    uuid: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
    alterId: 0
    cipher: auto
    tls: true
    servername: example.com
    network: ws
    ws-opts:
      path: /ws
      headers:
        Host: example.com
  - name: "VLESS-Reality"
    type: vless
    server: 5.6.7.8
    port: 443
    uuid: test-uuid
    flow: xtls-rprx-vision
    network: tcp
    tls: true
    client-fingerprint: chrome
    reality-opts:
      public-key: pubkey123
      short-id: ab12
  - name: "Trojan"
    type: trojan
    server: 9.10.11.12
    port: 443
    password: trojanpass
    sni: trojan.example.com
  - name: "SS"
    type: ss
    server: 13.14.15.16
    port: 8388
    cipher: aes-256-gcm
    password: sspass
`
	items, err := ParseBatch(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 4 {
		t.Fatalf("expected 4, got %d", len(items))
	}

	// VMess
	if items[0].ConfigType != model.ConfigVMess {
		t.Errorf("item[0] type = %v", items[0].ConfigType)
	}
	if items[0].Network != "ws" {
		t.Errorf("item[0] network = %q", items[0].Network)
	}
	if items[0].Path != "/ws" {
		t.Errorf("item[0] path = %q", items[0].Path)
	}
	if items[0].Host != "example.com" {
		t.Errorf("item[0] host = %q", items[0].Host)
	}

	// VLESS Reality
	if items[1].ConfigType != model.ConfigVLESS {
		t.Errorf("item[1] type = %v", items[1].ConfigType)
	}
	if items[1].StreamSecurity != "reality" {
		t.Errorf("item[1] streamSecurity = %q", items[1].StreamSecurity)
	}
	if items[1].PublicKey != "pubkey123" {
		t.Errorf("item[1] publicKey = %q", items[1].PublicKey)
	}
	if items[1].Flow != "xtls-rprx-vision" {
		t.Errorf("item[1] flow = %q", items[1].Flow)
	}

	// Trojan
	if items[2].ConfigType != model.ConfigTrojan {
		t.Errorf("item[2] type = %v", items[2].ConfigType)
	}
	if items[2].UUID != "trojanpass" {
		t.Errorf("item[2] password = %q", items[2].UUID)
	}

	// Shadowsocks
	if items[3].ConfigType != model.ConfigShadowsocks {
		t.Errorf("item[3] type = %v", items[3].ConfigType)
	}
	if items[3].Security != "aes-256-gcm" {
		t.Errorf("item[3] cipher = %q", items[3].Security)
	}
}

func TestDetectFormatClashYAML(t *testing.T) {
	content := "proxies:\n  - name: test\n    type: vmess\n"
	got := DetectFormat(content)
	if got != "yaml-clash" {
		t.Errorf("DetectFormat = %q, want yaml-clash", got)
	}
}

func TestParseSIP008(t *testing.T) {
	content := `{
		"servers": [
			{"id":"s1","remarks":"SIP008 Node","server":"10.0.0.1","server_port":8388,"password":"pass","method":"aes-256-gcm"}
		]
	}`
	items, err := ParseBatch(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1, got %d", len(items))
	}
	if items[0].Security != "aes-256-gcm" {
		t.Errorf("method = %q", items[0].Security)
	}
}

// --- Edge case tests ---

func TestParseURIEmpty(t *testing.T) {
	_, err := ParseURI("")
	if err == nil {
		t.Error("expected error for empty URI")
	}
}

func TestParseURIMalformedScheme(t *testing.T) {
	// These URIs should error.
	errorURIs := []string{
		"://missing-scheme",
		"vmess://not-valid-base64!!!",
	}
	for _, uri := range errorURIs {
		_, err := ParseURI(uri)
		if err == nil {
			t.Errorf("expected error for %q", uri)
		}
	}

	// Minimal URIs with known schemes are parsed leniently (no error).
	lenientURIs := []string{
		"vless://",
		"trojan://",
	}
	for _, uri := range lenientURIs {
		_, err := ParseURI(uri)
		if err != nil {
			t.Errorf("unexpected error for lenient URI %q: %v", uri, err)
		}
	}
}

func TestParseVLESSIPv6(t *testing.T) {
	uri := "vless://uuid@[::1]:443?security=tls&sni=example.com&type=tcp#IPv6Node"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatalf("ParseURI IPv6: %v", err)
	}
	if p.ConfigType != model.ConfigVLESS {
		t.Errorf("configType = %v", p.ConfigType)
	}
	if p.Port != 443 {
		t.Errorf("port = %d, want 443", p.Port)
	}
	if p.Remarks != "IPv6Node" {
		t.Errorf("remarks = %q", p.Remarks)
	}
}

func TestParseTrojanUnicodeRemarks(t *testing.T) {
	uri := "trojan://pass@1.2.3.4:443?sni=a.com#日本東京節點"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.Remarks != "日本東京節點" {
		t.Errorf("remarks = %q, want 日本東京節點", p.Remarks)
	}
}

func TestParseTrojanSpecialCharsPassword(t *testing.T) {
	uri := "trojan://p%40ss%3Aw0rd%21@1.2.3.4:443?sni=a.com#Node"
	p, err := ParseURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if p.UUID != "p@ss:w0rd!" {
		t.Errorf("password = %q, want p@ss:w0rd!", p.UUID)
	}
}

func TestParseBatchEmpty(t *testing.T) {
	items, err := ParseBatch("")
	if err != nil && len(items) != 0 {
		t.Fatalf("empty batch should return 0 items or error, got %d items err=%v", len(items), err)
	}
}

func TestParseBatchMixedValidInvalid(t *testing.T) {
	lines := "vless://uuid@1.2.3.4:443?security=tls&sni=a.com#Good\nhttp://invalid.com\ntrojan://pass@5.6.7.8:443?sni=b.com#Good2"
	items, err := ParseBatch(lines)
	if err != nil {
		t.Fatal(err)
	}
	// Should parse the valid ones and skip invalid.
	if len(items) < 2 {
		t.Errorf("expected at least 2 valid items, got %d", len(items))
	}
}

func TestParseBatchPlainTextLines(t *testing.T) {
	lines := "vless://uuid@1.2.3.4:443?security=tls&sni=a.com#Node1\ntrojan://pass@5.6.7.8:443?sni=b.com#Node2"
	items, err := ParseBatch(lines)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2, got %d", len(items))
	}
}

func TestDetectFormatPlainText(t *testing.T) {
	got := DetectFormat("just some random text without any structure")
	// Should not crash; result should be "unknown" or similar.
	if got == "" {
		t.Error("DetectFormat should return a non-empty string")
	}
}

func TestParseVLESSAllTransports(t *testing.T) {
	transports := []struct {
		name    string
		typeVal string
	}{
		{"tcp", "tcp"},
		{"ws", "ws"},
		{"grpc", "grpc"},
		{"h2", "h2"},
		{"httpupgrade", "httpupgrade"},
	}
	for _, tt := range transports {
		t.Run(tt.name, func(t *testing.T) {
			uri := "vless://uuid@1.2.3.4:443?security=tls&sni=a.com&type=" + tt.typeVal + "#" + tt.name
			p, err := ParseURI(uri)
			if err != nil {
				t.Fatalf("ParseURI %s: %v", tt.name, err)
			}
			if p.Network != tt.typeVal {
				t.Errorf("network = %q, want %q", p.Network, tt.typeVal)
			}
		})
	}
}
