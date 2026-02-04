package model

import (
	"errors"

	"github.com/RayUI/RayUI/internal/util"
)

// ProfileItem is the unified node/server configuration model covering all protocols.
type ProfileItem struct {
	// Identity
	ID         string      `json:"id"`
	ConfigType EConfigType `json:"configType"`
	Remarks    string      `json:"remarks"`
	SubID      string      `json:"subId"`
	ShareURI   string      `json:"shareUri"`
	Sort       int         `json:"sort"`

	// Connection
	Address string `json:"address"`
	Port    int    `json:"port"`
	Ports   string `json:"ports,omitempty"`

	// Authentication
	UUID     string `json:"uuid"`
	AlterID  int    `json:"alterId,omitempty"`
	Security string `json:"security"`
	Flow     string `json:"flow,omitempty"`

	// Transport
	Network    string `json:"network"`
	HeaderType string `json:"headerType,omitempty"`
	Host       string `json:"host,omitempty"`
	Path       string `json:"path,omitempty"`

	// TLS / Security
	StreamSecurity string `json:"streamSecurity"`
	AllowInsecure  bool   `json:"allowInsecure"`
	SNI            string `json:"sni,omitempty"`
	ALPN           string `json:"alpn,omitempty"`
	Fingerprint    string `json:"fingerprint,omitempty"`

	// Reality
	PublicKey string `json:"publicKey,omitempty"`
	ShortID   string `json:"shortId,omitempty"`
	SpiderX   string `json:"spiderX,omitempty"`

	// Advanced
	CoreType   ECoreType `json:"coreType,omitempty"`
	Extra      string    `json:"extra,omitempty"`
	MuxEnabled *bool     `json:"muxEnabled,omitempty"`
}

// NewProfileItem returns a ProfileItem with a generated UUID and sensible defaults.
func NewProfileItem() ProfileItem {
	return ProfileItem{
		ID:             util.GenerateUUID(),
		Network:        "tcp",
		StreamSecurity: "none",
		Security:       "auto",
	}
}

// Validate performs basic validation.
func (p ProfileItem) Validate() error {
	if p.Address == "" {
		return errors.New("address is required")
	}
	if p.Port <= 0 || p.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
}
