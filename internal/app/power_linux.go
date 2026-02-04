//go:build linux

package app

// Start starts monitoring power events on Linux.
// TODO: Implement using D-Bus org.freedesktop.login1
func (m *PowerMonitor) Start() {
	// Linux power monitoring not yet implemented
	// Would use D-Bus to listen to org.freedesktop.login1.Manager PrepareForSleep signal
}

// StopPlatform stops platform-specific monitoring.
func (m *PowerMonitor) StopPlatform() {
	// Linux power monitoring not yet implemented
}
