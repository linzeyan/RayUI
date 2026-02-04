package core

import (
	"io"

	"github.com/RayUI/RayUI/internal/model"
)

// CoreManager controls the lifecycle of a proxy core binary.
type CoreManager interface {
	Start(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, config model.Config) error
	Stop() error
	Restart() error
	IsRunning() bool
	GetStatus() model.CoreStatus
	GenerateConfig(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, config model.Config) ([]byte, error)
	CoreType() model.ECoreType
	Version() (string, error)
	BinaryPath() string
	SetLogWriter(w io.Writer)
}

// SelectCore determines which core to use for a given profile.
func SelectCore(profile model.ProfileItem) model.ECoreType {
	if profile.CoreType != model.CoreAuto {
		return profile.CoreType
	}
	switch {
	case profile.ConfigType == model.ConfigHysteria2:
		return model.CoreSingbox
	case profile.ConfigType == model.ConfigTUIC:
		return model.CoreSingbox
	case profile.Network == "grpc" || profile.Network == "h2":
		return model.CoreSingbox
	default:
		return model.CoreXray
	}
}

// NewCoreManager returns a CoreManager for the given core type.
func NewCoreManager(coreType model.ECoreType, dataDir string) CoreManager {
	switch coreType {
	case model.CoreSingbox:
		return NewSingboxCore(dataDir)
	default:
		return NewXrayCore(dataDir)
	}
}
