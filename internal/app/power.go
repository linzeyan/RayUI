package app

// PowerEventCallback is called when power state changes.
type PowerEventCallback func(isSleeping bool)

// PowerMonitor monitors system sleep/wake events.
type PowerMonitor struct {
	callback PowerEventCallback
	stopCh   chan struct{}
}

// NewPowerMonitor creates a new power monitor.
func NewPowerMonitor(callback PowerEventCallback) *PowerMonitor {
	return &PowerMonitor{
		callback: callback,
		stopCh:   make(chan struct{}),
	}
}

// Stop stops the power monitor.
func (m *PowerMonitor) Stop() {
	select {
	case <-m.stopCh:
		// Already stopped
	default:
		close(m.stopCh)
	}
}
