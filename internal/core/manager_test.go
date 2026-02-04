package core

import (
	"testing"

	"github.com/RayUI/RayUI/internal/model"
)

func TestSelectCore(t *testing.T) {
	tests := []struct {
		name    string
		profile model.ProfileItem
		want    model.ECoreType
	}{
		{
			name:    "default → xray",
			profile: model.ProfileItem{ConfigType: model.ConfigVMess},
			want:    model.CoreXray,
		},
		{
			name:    "hysteria2 → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigHysteria2},
			want:    model.CoreSingbox,
		},
		{
			name:    "tuic → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigTUIC},
			want:    model.CoreSingbox,
		},
		{
			name:    "grpc → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, Network: "grpc"},
			want:    model.CoreSingbox,
		},
		{
			name:    "h2 → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, Network: "h2"},
			want:    model.CoreSingbox,
		},
		{
			name:    "override xray",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, CoreType: model.CoreXray, Network: "grpc"},
			want:    model.CoreXray,
		},
		{
			name:    "override singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVMess, CoreType: model.CoreSingbox},
			want:    model.CoreSingbox,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectCore(tt.profile)
			if got != tt.want {
				t.Errorf("SelectCore = %v, want %v", got, tt.want)
			}
		})
	}
}
