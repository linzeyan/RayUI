package sysproxy

// SysProxy controls the operating system's proxy settings.
type SysProxy interface {
	Set(httpAddr string, httpPort int, socksAddr string, socksPort int) error
	Clear() error
	GetCurrent() (*ProxyState, error)
}

// ProxyState holds the current system proxy configuration.
type ProxyState struct {
	Enabled   bool   `json:"enabled"`
	HTTPHost  string `json:"httpHost"`
	HTTPPort  int    `json:"httpPort"`
	SOCKSHost string `json:"socksHost"`
	SOCKSPort int    `json:"socksPort"`
}

// NewSysProxy returns a platform-specific SysProxy implementation.
func NewSysProxy() SysProxy {
	return newPlatformSysProxy()
}
