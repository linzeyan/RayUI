package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

func parseVLESS(uri string) (*model.ProfileItem, error) {
	// vless://uuid@host:port?params#fragment
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("vless url parse: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigVLESS
	p.UUID = u.User.Username()
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Remarks = u.Fragment

	q := u.Query()
	p.Security = q.Get("encryption")
	if p.Security == "" {
		p.Security = "none"
	}
	p.StreamSecurity = q.Get("security")
	if p.StreamSecurity == "" {
		p.StreamSecurity = "none"
	}
	p.SNI = q.Get("sni")
	p.Network = q.Get("type")
	if p.Network == "" {
		p.Network = "tcp"
	}
	p.Host = q.Get("host")
	p.Path = q.Get("path")
	p.Fingerprint = q.Get("fp")
	p.ALPN = q.Get("alpn")
	p.Flow = q.Get("flow")
	p.HeaderType = q.Get("headerType")

	// Reality fields
	p.PublicKey = q.Get("pbk")
	p.ShortID = q.Get("sid")
	p.SpiderX = q.Get("spx")

	p.ShareURI = uri
	return &p, nil
}

func toVLESSURI(item model.ProfileItem) string {
	q := url.Values{}
	if item.Security != "" && item.Security != "none" {
		q.Set("encryption", item.Security)
	}
	if item.StreamSecurity != "" && item.StreamSecurity != "none" {
		q.Set("security", item.StreamSecurity)
	}
	if item.SNI != "" {
		q.Set("sni", item.SNI)
	}
	if item.Network != "" && item.Network != "tcp" {
		q.Set("type", item.Network)
	}
	if item.Host != "" {
		q.Set("host", item.Host)
	}
	if item.Path != "" {
		q.Set("path", item.Path)
	}
	if item.Fingerprint != "" {
		q.Set("fp", item.Fingerprint)
	}
	if item.ALPN != "" {
		q.Set("alpn", item.ALPN)
	}
	if item.Flow != "" {
		q.Set("flow", item.Flow)
	}
	if item.HeaderType != "" {
		q.Set("headerType", item.HeaderType)
	}
	if item.PublicKey != "" {
		q.Set("pbk", item.PublicKey)
	}
	if item.ShortID != "" {
		q.Set("sid", item.ShortID)
	}
	if item.SpiderX != "" {
		q.Set("spx", item.SpiderX)
	}

	var sb strings.Builder
	sb.WriteString("vless://")
	sb.WriteString(item.UUID)
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
