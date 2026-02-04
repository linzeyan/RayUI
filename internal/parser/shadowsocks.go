package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

func parseShadowsocks(uri string) (*model.ProfileItem, error) {
	// ss://base64(method:password)@host:port#remarks
	// or ss://base64(method:password@host:port)#remarks
	raw := strings.TrimPrefix(uri, "ss://")

	// Split off the fragment (remarks).
	remarks := ""
	if idx := strings.LastIndex(raw, "#"); idx >= 0 {
		remarks, _ = url.PathUnescape(raw[idx+1:])
		raw = raw[:idx]
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigShadowsocks
	p.Remarks = remarks
	p.ShareURI = uri

	// Format 1: base64(method:password)@host:port
	if atIdx := strings.LastIndex(raw, "@"); atIdx >= 0 {
		userPart := raw[:atIdx]
		hostPort := raw[atIdx+1:]

		decoded, err := decodeSSBase64(userPart)
		if err != nil {
			return nil, fmt.Errorf("ss userinfo decode: %w", err)
		}

		method, password, ok := strings.Cut(decoded, ":")
		if !ok {
			return nil, fmt.Errorf("ss userinfo: expected method:password")
		}
		p.Security = method
		p.UUID = password

		host, port, err := splitHostPort(hostPort)
		if err != nil {
			return nil, err
		}
		p.Address = host
		p.Port = port

		return &p, nil
	}

	// Format 2: base64(method:password@host:port)
	decoded, err := decodeSSBase64(raw)
	if err != nil {
		return nil, fmt.Errorf("ss full decode: %w", err)
	}

	methodAndRest, hostPort, ok := splitSSFull(decoded)
	if !ok {
		return nil, fmt.Errorf("ss full format: invalid")
	}

	method, password, ok := strings.Cut(methodAndRest, ":")
	if !ok {
		return nil, fmt.Errorf("ss full format: expected method:password")
	}
	p.Security = method
	p.UUID = password

	host, port, err := splitHostPort(hostPort)
	if err != nil {
		return nil, err
	}
	p.Address = host
	p.Port = port

	return &p, nil
}

func toShadowsocksURI(item model.ProfileItem) string {
	userInfo := item.Security + ":" + item.UUID
	encoded := base64.URLEncoding.EncodeToString([]byte(userInfo))
	// Remove padding.
	encoded = strings.TrimRight(encoded, "=")

	var sb strings.Builder
	sb.WriteString("ss://")
	sb.WriteString(encoded)
	sb.WriteString("@")
	sb.WriteString(item.Address)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(item.Port))
	if item.Remarks != "" {
		sb.WriteString("#")
		sb.WriteString(url.PathEscape(item.Remarks))
	}
	return sb.String()
}

func decodeSSBase64(s string) (string, error) {
	// Try standard, then URL-safe, with and without padding.
	for _, enc := range []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	} {
		decoded, err := enc.DecodeString(s)
		if err == nil {
			return string(decoded), nil
		}
	}
	return "", fmt.Errorf("unable to base64 decode: %s", s)
}

func splitHostPort(s string) (string, int, error) {
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return "", 0, fmt.Errorf("missing port in %q", s)
	}
	host := s[:idx]
	port, err := strconv.Atoi(s[idx+1:])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port in %q: %w", s, err)
	}
	return host, port, nil
}

// splitSSFull splits "method:password@host:port" at the last @ before the host:port.
func splitSSFull(s string) (string, string, bool) {
	atIdx := strings.LastIndex(s, "@")
	if atIdx < 0 {
		return "", "", false
	}
	return s[:atIdx], s[atIdx+1:], true
}
