package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// parseHysteria2 parses a hysteria2:// URI.
// Format: hysteria2://auth@host:port?sni=xxx&insecure=1&obfs=salamander&obfs-password=xxx&alpn=h3#name
func parseHysteria2(uri string) (*model.ProfileItem, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("hysteria2 url parse: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigHysteria2
	p.CoreType = model.CoreSingbox // Hysteria2 only supported by sing-box
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Remarks = u.Fragment

	// Auth password is in the userinfo.
	p.UUID = u.User.Username()

	q := u.Query()
	p.SNI = q.Get("sni")
	p.ALPN = q.Get("alpn")
	p.Fingerprint = q.Get("fp")

	if q.Get("insecure") == "1" {
		p.AllowInsecure = true
	}

	// Obfuscation.
	p.HeaderType = q.Get("obfs")
	p.Path = q.Get("obfs-password")

	p.StreamSecurity = "tls"
	p.ShareURI = uri
	return &p, nil
}

func toHysteria2URI(item model.ProfileItem) string {
	q := url.Values{}
	if item.SNI != "" {
		q.Set("sni", item.SNI)
	}
	if item.ALPN != "" {
		q.Set("alpn", item.ALPN)
	}
	if item.Fingerprint != "" {
		q.Set("fp", item.Fingerprint)
	}
	if item.AllowInsecure {
		q.Set("insecure", "1")
	}
	if item.HeaderType != "" {
		q.Set("obfs", item.HeaderType)
	}
	if item.Path != "" {
		q.Set("obfs-password", item.Path)
	}

	var sb strings.Builder
	sb.WriteString("hysteria2://")
	sb.WriteString(url.PathEscape(item.UUID))
	sb.WriteString("@")
	sb.WriteString(item.Address)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(item.Port))
	if encoded := q.Encode(); encoded != "" {
		sb.WriteString("?")
		sb.WriteString(encoded)
	}
	if item.Remarks != "" {
		sb.WriteString("#")
		sb.WriteString(url.PathEscape(item.Remarks))
	}
	return sb.String()
}
