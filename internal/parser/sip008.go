package parser

import (
	"encoding/json"
	"fmt"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/util"
)

// sip008Response is the SIP008 JSON format.
type sip008Response struct {
	Servers []sip008Server `json:"servers"`
}

type sip008Server struct {
	ID         string `json:"id"`
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

func parseSIP008(content string) ([]model.ProfileItem, error) {
	var resp sip008Response
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return nil, fmt.Errorf("sip008 parse: %w", err)
	}

	var items []model.ProfileItem
	for _, s := range resp.Servers {
		p := model.ProfileItem{
			ID:         util.GenerateUUID(),
			ConfigType: model.ConfigShadowsocks,
			Remarks:    s.Remarks,
			Address:    s.Server,
			Port:       s.ServerPort,
			UUID:       s.Password,
			Security:   s.Method,
			Network:    "tcp",
			StreamSecurity: "none",
		}
		items = append(items, p)
	}
	return items, nil
}
