package service

import (
	"testing"
	"time"

	"github.com/RayUI/RayUI/internal/model"
)

func TestStatsServiceSetCallback(t *testing.T) {
	s := &StatsService{}

	called := false
	s.SetCallback(func(stats model.TrafficStats) {
		called = true
	})

	s.mu.Lock()
	cb := s.callback
	s.mu.Unlock()

	if cb == nil {
		t.Error("callback should not be nil after SetCallback")
	}

	cb(model.TrafficStats{})
	if !called {
		t.Error("callback was not invoked")
	}
}

func TestStatsServiceGetCurrentTraffic(t *testing.T) {
	s := &StatsService{}

	// Default should be zero
	stats := s.GetCurrentTraffic()
	if stats.Up != 0 || stats.Down != 0 {
		t.Errorf("GetCurrentTraffic() = {Up:%d, Down:%d}, want {0, 0}", stats.Up, stats.Down)
	}

	// Set current and verify
	s.mu.Lock()
	s.current = model.TrafficStats{Up: 100, Down: 200}
	s.mu.Unlock()

	stats = s.GetCurrentTraffic()
	if stats.Up != 100 {
		t.Errorf("Up = %d, want 100", stats.Up)
	}
	if stats.Down != 200 {
		t.Errorf("Down = %d, want 200", stats.Down)
	}
}

func TestStatsServiceStartStopCollecting(t *testing.T) {
	s := &StatsService{}

	// StopCollecting with nil cancel should not panic
	s.StopCollecting()

	// Start and stop should not panic
	s.StartCollecting("profile-1", model.CoreXray)

	// Allow goroutine to start
	time.Sleep(10 * time.Millisecond)

	s.StopCollecting()

	// Verify cancel is cleared
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()
	if cancel != nil {
		t.Error("cancel should be nil after StopCollecting")
	}
}

func TestStatsServiceStartCollectingRestartsLoop(t *testing.T) {
	s := &StatsService{}

	// Start twice - should stop previous and start new
	s.StartCollecting("profile-1", model.CoreXray)
	time.Sleep(10 * time.Millisecond)
	s.StartCollecting("profile-2", model.CoreSingbox)
	time.Sleep(10 * time.Millisecond)

	// Verify cancel is not nil (active collection)
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()
	if cancel == nil {
		t.Error("cancel should not be nil during active collection")
	}

	s.StopCollecting()
}

func TestStatsServiceCallbackInvocation(t *testing.T) {
	s := &StatsService{}

	var receivedStats model.TrafficStats
	s.SetCallback(func(stats model.TrafficStats) {
		receivedStats = stats
	})

	// Simulate what collectLoop does: set current and invoke callback
	s.mu.Lock()
	s.current = model.TrafficStats{Up: 50, Down: 75}
	cb := s.callback
	s.mu.Unlock()

	if cb != nil {
		cb(model.TrafficStats{Up: 50, Down: 75})
	}

	if receivedStats.Up != 50 {
		t.Errorf("receivedStats.Up = %d, want 50", receivedStats.Up)
	}
	if receivedStats.Down != 75 {
		t.Errorf("receivedStats.Down = %d, want 75", receivedStats.Down)
	}
}

func TestFetchTraffic(t *testing.T) {
	s := &StatsService{}

	// fetchTraffic is currently a stub that returns 0,0
	up, down := s.fetchTraffic(model.CoreXray)
	if up != 0 || down != 0 {
		t.Errorf("fetchTraffic() = (%d, %d), want (0, 0)", up, down)
	}
}
