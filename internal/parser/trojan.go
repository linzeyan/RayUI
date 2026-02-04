package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

func parseTrojan(uri string) (*model.ProfileItem, error) {
	// trojan://password@host:port?params#fragment
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("trojan url parse: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigTrojan
	p.UUID = u.User.Username() // password stored in UUID field
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Remarks = u.Fragment

	q := u.Query()
	p.StreamSecurity = q.Get("security")
	if p.StreamSecurity == "" {
		p.StreamSecurity = "tls"
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

	p.ShareURI = uri
	return &p, nil
}

func toTrojanURI(item model.ProfileItem) string {
	q := url.Values{}
	if item.StreamSecurity != "" && item.StreamSecurity != "tls" {
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

	var sb strings.Builder
	sb.WriteString("trojan://")
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
