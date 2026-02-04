package parser

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// ParseURI parses a single proxy share URI into a ProfileItem.
func ParseURI(uri string) (*model.ProfileItem, error) {
	uri = strings.TrimSpace(uri)
	switch {
	case strings.HasPrefix(uri, "vmess://"):
		return parseVMess(uri)
	case strings.HasPrefix(uri, "vless://"):
		return parseVLESS(uri)
	case strings.HasPrefix(uri, "trojan://"):
		return parseTrojan(uri)
	case strings.HasPrefix(uri, "ss://"):
		return parseShadowsocks(uri)
	case strings.HasPrefix(uri, "hysteria2://"), strings.HasPrefix(uri, "hy2://"):
		return parseHysteria2(uri)
	case strings.HasPrefix(uri, "tuic://"):
		return parseTUIC(uri)
	case strings.HasPrefix(uri, "wireguard://"), strings.HasPrefix(uri, "wg://"):
		return parseWireGuard(uri)
	default:
		return nil, &ErrUnsupportedScheme{URI: uri}
	}
}

// ParseBatch parses text containing multiple proxy URIs (one per line).
// If the text appears to be Base64-encoded, it is decoded first.
// If the text is a sing-box JSON config, outbounds are extracted.
// If the text is SIP008 JSON, servers are extracted.
func ParseBatch(text string) ([]model.ProfileItem, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}

	format := DetectFormat(text)
	switch format {
	case "json-singbox":
		return parseSingboxOutbounds(text)
	case "json-sip008":
		return parseSIP008(text)
	case "yaml-clash":
		return parseClashYAML(text)
	case "base64":
		decoded, err := tryBase64Decode(text)
		if err != nil {
			return nil, err
		}
		return parseLines(decoded)
	default:
		return parseLines(text)
	}
}

// DetectFormat returns the detected subscription format.
func DetectFormat(content string) string {
	trimmed := strings.TrimSpace(content)

	// Try JSON detection first.
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var obj map[string]json.RawMessage
		if json.Unmarshal([]byte(trimmed), &obj) == nil {
			if _, ok := obj["outbounds"]; ok {
				return "json-singbox"
			}
			if _, ok := obj["servers"]; ok {
				return "json-sip008"
			}
		}
	}

	// Check Clash YAML (contains "proxies:" key).
	if looksLikeClashYAML(trimmed) {
		return "yaml-clash"
	}

	// Check if it looks like plain URI lines.
	if looksLikeURILines(trimmed) {
		return "uri-lines"
	}

	// Try Base64 decode.
	if _, err := tryBase64Decode(trimmed); err == nil {
		return "base64"
	}

	return "unknown"
}

// ToShareURI converts a ProfileItem back to a share URI.
func ToShareURI(item model.ProfileItem) (string, error) {
	switch item.ConfigType {
	case model.ConfigVMess:
		return toVMessURI(item), nil
	case model.ConfigVLESS:
		return toVLESSURI(item), nil
	case model.ConfigTrojan:
		return toTrojanURI(item), nil
	case model.ConfigShadowsocks:
		return toShadowsocksURI(item), nil
	case model.ConfigHysteria2:
		return toHysteria2URI(item), nil
	case model.ConfigTUIC:
		return toTUICURI(item), nil
	case model.ConfigWireGuard:
		return toWireGuardURI(item), nil
	default:
		return "", &ErrUnsupportedScheme{URI: item.ConfigType.String()}
	}
}

func parseLines(text string) ([]model.ProfileItem, error) {
	var items []model.ProfileItem
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		item, err := ParseURI(line)
		if err != nil {
			continue // skip unparseable lines
		}
		items = append(items, *item)
	}
	return items, nil
}

func tryBase64Decode(s string) (string, error) {
	// Remove whitespace/newlines for base64 decoding.
	cleaned := strings.Join(strings.Fields(s), "")
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(cleaned)
	}
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(cleaned)
	}
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(cleaned)
	}
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func looksLikeClashYAML(s string) bool {
	return strings.Contains(s, "proxies:") &&
		(strings.Contains(s, "- name:") || strings.Contains(s, "- {name:"))
}

func looksLikeURILines(s string) bool {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return false
	}
	first := strings.TrimSpace(lines[0])
	return strings.HasPrefix(first, "vmess://") ||
		strings.HasPrefix(first, "vless://") ||
		strings.HasPrefix(first, "trojan://") ||
		strings.HasPrefix(first, "ss://") ||
		strings.HasPrefix(first, "hysteria2://") ||
		strings.HasPrefix(first, "hy2://") ||
		strings.HasPrefix(first, "tuic://") ||
		strings.HasPrefix(first, "wireguard://") ||
		strings.HasPrefix(first, "wg://")
}

// ErrUnsupportedScheme is returned for unknown URI schemes.
type ErrUnsupportedScheme struct {
	URI string
}

func (e *ErrUnsupportedScheme) Error() string {
	return "unsupported scheme: " + e.URI
}
