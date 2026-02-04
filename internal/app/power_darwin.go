//go:build darwin

package app

import (
	"time"
)

// Start starts monitoring power events on macOS.
// Uses time-gap detection to identify sleep/wake cycles.
func (m *PowerMonitor) Start() {
	go m.detectSleepWake()
}

// StopPlatform stops platform-specific monitoring.
func (m *PowerMonitor) StopPlatform() {
	// Stop is handled by the stopCh channel
}

// detectSleepWake uses time gap detection to identify sleep/wake events.
// If the time between ticks is significantly larger than expected,
// the system was likely asleep.
func (m *PowerMonitor) detectSleepWake() {
	const tickInterval = 2 * time.Second
	const sleepThreshold = 5 * time.Second // If gap > 5s, assume sleep occurred

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	lastTick := time.Now()
	wasSleeping := false

	for {
		select {
		case <-m.stopCh:
			return
		case now := <-ticker.C:
			gap := now.Sub(lastTick)
			lastTick = now

			// If the gap is much larger than the tick interval, system was asleep
			if gap > sleepThreshold && !wasSleeping {
				// System just woke up
				wasSleeping = false
				if m.callback != nil {
					m.callback(false) // wake event
				}
			} else if gap <= sleepThreshold && wasSleeping {
				wasSleeping = false
			}

			// Note: We can't directly detect "going to sleep" with this method
			// The callback for sleep would need native APIs (CGO)
			// For now, we only detect wake events reliably
		}
	}
}
