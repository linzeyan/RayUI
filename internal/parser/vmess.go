package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/util"
)

// vmessJSON represents the JSON body inside a vmess:// URI.
type vmessJSON struct {
	V    string `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"` // string or int
	ID   string `json:"id"`
	AID  any    `json:"aid"` // string or int
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
	FP   string `json:"fp"`
}

func parseVMess(uri string) (*model.ProfileItem, error) {
	raw := strings.TrimPrefix(uri, "vmess://")
	// Try base64 decode.
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(raw)
	}
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(raw)
	}
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(raw)
	}
	if err != nil {
		return nil, fmt.Errorf("vmess base64 decode: %w", err)
	}

	var v vmessJSON
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, fmt.Errorf("vmess json: %w", err)
	}

	p := model.NewProfileItem()
	p.ConfigType = model.ConfigVMess
	p.Remarks = v.PS
	p.Address = v.Add
	p.Port = anyToInt(v.Port)
	p.UUID = v.ID
	p.AlterID = anyToInt(v.AID)
	p.Security = v.Scy
	if p.Security == "" {
		p.Security = "auto"
	}
	p.Network = v.Net
	if p.Network == "" {
		p.Network = "tcp"
	}
	p.HeaderType = v.Type
	p.Host = v.Host
	p.Path = v.Path
	p.StreamSecurity = v.TLS
	if p.StreamSecurity == "" {
		p.StreamSecurity = "none"
	}
	p.SNI = v.SNI
	p.ALPN = v.ALPN
	p.Fingerprint = v.FP
	p.ShareURI = uri

	return &p, nil
}

func toVMessURI(item model.ProfileItem) string {
	v := vmessJSON{
		V:    "2",
		PS:   item.Remarks,
		Add:  item.Address,
		Port: strconv.Itoa(item.Port),
		ID:   item.UUID,
		AID:  strconv.Itoa(item.AlterID),
		Scy:  item.Security,
		Net:  item.Network,
		Type: item.HeaderType,
		Host: item.Host,
		Path: item.Path,
		TLS:  item.StreamSecurity,
		SNI:  item.SNI,
		ALPN: item.ALPN,
		FP:   item.Fingerprint,
	}
	data, _ := json.Marshal(v)
	return "vmess://" + base64.StdEncoding.EncodeToString(data)
}

func anyToInt(v any) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		n, _ := strconv.Atoi(val)
		return n
	case int:
		return val
	default:
		return 0
	}
}

// GenerateUUID is a utility re-export for tests.
var generateUUID = util.GenerateUUID
