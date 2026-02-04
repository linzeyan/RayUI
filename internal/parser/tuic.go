package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// parseTUIC parses a tuic:// URI.
// Format: tuic://uuid:password@host:port?congestion_control=bbr&udp_relay_mode=native&alpn=h3&sni=xxx#name
func parseTUIC(uri string) (*model.ProfileItem, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("tuic url parse: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigTUIC
	p.CoreType = model.CoreSingbox // TUIC only supported by sing-box
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Remarks = u.Fragment

	// UUID and password from userinfo.
	p.UUID = u.User.Username()
	p.Security, _ = u.User.Password()

	q := u.Query()
	p.SNI = q.Get("sni")
	p.ALPN = q.Get("alpn")
	p.Fingerprint = q.Get("fp")

	if q.Get("allow_insecure") == "1" || q.Get("insecure") == "1" {
		p.AllowInsecure = true
	}

	// congestion_control → HeaderType, udp_relay_mode → Path
	p.HeaderType = q.Get("congestion_control")
	p.Path = q.Get("udp_relay_mode")

	p.StreamSecurity = "tls"
	p.ShareURI = uri
	return &p, nil
}

func toTUICURI(item model.ProfileItem) string {
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
		q.Set("allow_insecure", "1")
	}
	if item.HeaderType != "" {
		q.Set("congestion_control", item.HeaderType)
	}
	if item.Path != "" {
		q.Set("udp_relay_mode", item.Path)
	}

	var sb strings.Builder
	sb.WriteString("tuic://")
	sb.WriteString(item.UUID)
	sb.WriteString(":")
	sb.WriteString(url.PathEscape(item.Security))
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
