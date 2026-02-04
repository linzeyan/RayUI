package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// parseWireGuard parses a wireguard:// URI.
// Format: wireguard://private-key@host:port?publickey=xxx&address=10.0.0.2/32&mtu=1420&reserved=0,0,0#name
func parseWireGuard(uri string) (*model.ProfileItem, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("wireguard url parse: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigWireGuard
	p.CoreType = model.CoreSingbox // WireGuard best supported by sing-box
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Remarks = u.Fragment

	// Private key in userinfo.
	p.UUID = u.User.Username()

	q := u.Query()
	p.PublicKey = q.Get("publickey")
	p.Host = q.Get("address") // tunnel address
	p.ShortID = q.Get("reserved")
	p.SNI = q.Get("sni")

	if mtu := q.Get("mtu"); mtu != "" {
		p.Extra = mtu
	}

	p.Network = "wireguard"
	p.StreamSecurity = "none"
	p.ShareURI = uri
	return &p, nil
}

func toWireGuardURI(item model.ProfileItem) string {
	q := url.Values{}
	if item.PublicKey != "" {
		q.Set("publickey", item.PublicKey)
	}
	if item.Host != "" {
		q.Set("address", item.Host)
	}
	if item.ShortID != "" {
		q.Set("reserved", item.ShortID)
	}
	if item.SNI != "" {
		q.Set("sni", item.SNI)
	}
	if item.Extra != "" {
		q.Set("mtu", item.Extra)
	}

	var sb strings.Builder
	sb.WriteString("wireguard://")
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
