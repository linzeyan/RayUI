package model

// DNSItem holds DNS configuration.
type DNSItem struct {
	RemoteDNS      string `json:"remoteDns"`
	DirectDNS      string `json:"directDns"`
	BootstrapDNS   string `json:"bootstrapDns"`
	UseSystemHosts bool   `json:"useSystemHosts"`
	FakeIP         bool   `json:"fakeIP"`
	Hosts          string `json:"hosts"`
	DomainStrategy string `json:"domainStrategy"`
}

// DefaultDNSItem returns a DNSItem with sensible defaults.
func DefaultDNSItem() DNSItem {
	return DNSItem{
		RemoteDNS:      "https://dns.google/dns-query",
		DirectDNS:      "https://dns.alidns.com/dns-query",
		BootstrapDNS:   "1.1.1.1",
		UseSystemHosts: false,
		FakeIP:         false,
		DomainStrategy: "prefer_ipv4",
	}
}
