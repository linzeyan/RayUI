package model

import "github.com/RayUI/RayUI/internal/util"

// RoutingItem is an ordered set of routing rules.
type RoutingItem struct {
	ID             string     `json:"id"`
	Remarks        string     `json:"remarks"`
	DomainStrategy string     `json:"domainStrategy"`
	Rules          []RuleItem `json:"rules"`
	Enabled        bool       `json:"enabled"`
	Locked         bool       `json:"locked"`
	Sort           int        `json:"sort"`
}

// RuleItem is a single routing rule inside a RoutingItem.
type RuleItem struct {
	ID          string `json:"id"`
	OutboundTag string `json:"outboundTag"`
	Enabled     bool   `json:"enabled"`
	Remarks     string `json:"remarks,omitempty"`

	// Match conditions
	Domain        []string `json:"domain,omitempty"`
	DomainSuffix  []string `json:"domainSuffix,omitempty"`
	DomainKeyword []string `json:"domainKeyword,omitempty"`
	DomainRegex   []string `json:"domainRegex,omitempty"`
	Geosite       []string `json:"geosite,omitempty"`
	IP            []string `json:"ip,omitempty"`
	IPCIDR        []string `json:"ipCidr,omitempty"`
	GeoIP         []string `json:"geoip,omitempty"`
	Port          string   `json:"port,omitempty"`
	Protocol      []string `json:"protocol,omitempty"`
	ProcessName   []string `json:"processName,omitempty"`
	Network       string   `json:"network,omitempty"`
	Inbound       []string `json:"inbound,omitempty"`
	RuleSet       []string `json:"ruleSet,omitempty"`
}

func newRule(outbound, remarks string) RuleItem {
	return RuleItem{
		ID:          util.GenerateUUID(),
		OutboundTag: outbound,
		Enabled:     true,
		Remarks:     remarks,
	}
}

// DefaultRoutingItems returns the 4 built-in routing presets.
func DefaultRoutingItems() []RoutingItem {
	return []RoutingItem{
		{
			ID:             util.GenerateUUID(),
			Remarks:        "Global",
			DomainStrategy: "AsIs",
			Locked:         true,
			Enabled:        false,
			Sort:           0,
			Rules: []RuleItem{
				func() RuleItem {
					r := newRule("proxy", "All traffic → proxy")
					r.Network = "tcp,udp"
					return r
				}(),
			},
		},
		{
			ID:             util.GenerateUUID(),
			Remarks:        "BypassLAN",
			DomainStrategy: "AsIs",
			Locked:         true,
			Enabled:        false,
			Sort:           1,
			Rules: []RuleItem{
				func() RuleItem {
					r := newRule("direct", "Private IPs → direct")
					r.GeoIP = []string{"private"}
					return r
				}(),
				func() RuleItem {
					r := newRule("proxy", "Remaining → proxy")
					r.Network = "tcp,udp"
					return r
				}(),
			},
		},
		{
			ID:             util.GenerateUUID(),
			Remarks:        "BypassCN",
			DomainStrategy: "IPIfNonMatch",
			Locked:         true,
			Enabled:        false,
			Sort:           2,
			Rules: []RuleItem{
				func() RuleItem {
					r := newRule("direct", "CN sites → direct")
					r.Geosite = []string{"cn"}
					return r
				}(),
				func() RuleItem {
					r := newRule("direct", "CN IPs → direct")
					r.GeoIP = []string{"cn"}
					return r
				}(),
				func() RuleItem {
					r := newRule("proxy", "Remaining → proxy")
					r.Network = "tcp,udp"
					return r
				}(),
			},
		},
		{
			ID:             util.GenerateUUID(),
			Remarks:        "BypassLAN+CN",
			DomainStrategy: "IPIfNonMatch",
			Locked:         true,
			Enabled:        true,
			Sort:           3,
			Rules: []RuleItem{
				func() RuleItem {
					r := newRule("block", "Ads → block")
					r.Geosite = []string{"category-ads-all"}
					return r
				}(),
				func() RuleItem {
					r := newRule("direct", "Private IPs → direct")
					r.GeoIP = []string{"private"}
					return r
				}(),
				func() RuleItem {
					r := newRule("direct", "CN sites → direct")
					r.Geosite = []string{"cn"}
					return r
				}(),
				func() RuleItem {
					r := newRule("direct", "CN IPs → direct")
					r.GeoIP = []string{"cn"}
					return r
				}(),
				func() RuleItem {
					r := newRule("proxy", "Remaining → proxy")
					r.Network = "tcp,udp"
					return r
				}(),
			},
		},
	}
}
