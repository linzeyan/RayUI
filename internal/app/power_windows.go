//go:build windows

package app

// Start starts monitoring power events on Windows.
// TODO: Implement using RegisterPowerSettingNotification
func (m *PowerMonitor) Start() {
	// Windows power monitoring not yet implemented
	// Would use RegisterPowerSettingNotification with GUID_SYSTEM_AWAYMODE
}

// StopPlatform stops platform-specific monitoring.
func (m *PowerMonitor) StopPlatform() {
	// Windows power monitoring not yet implemented
}
