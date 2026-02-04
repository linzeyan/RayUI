package autostart

// AutoStart manages launch-at-login behaviour.
type AutoStart interface {
	Enable(execPath string) error
	Disable() error
	IsEnabled() (bool, error)
}

// NewAutoStart returns a platform-specific AutoStart implementation.
func NewAutoStart() AutoStart {
	return newPlatformAutoStart()
}
