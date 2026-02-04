package model

import "github.com/RayUI/RayUI/internal/util"

// SubItem represents a subscription source.
type SubItem struct {
	ID                 string `json:"id"`
	Remarks            string `json:"remarks"`
	URL                string `json:"url"`
	Enabled            bool   `json:"enabled"`
	Sort               int    `json:"sort"`
	Filter             string `json:"filter,omitempty"`
	AutoUpdateInterval int    `json:"autoUpdateInterval"`
	UpdateTime         int64  `json:"updateTime"`
	UserAgent          string `json:"userAgent,omitempty"`
}

// NewSubItem returns a SubItem with a generated UUID and defaults.
func NewSubItem() SubItem {
	return SubItem{
		ID:      util.GenerateUUID(),
		Enabled: true,
	}
}
