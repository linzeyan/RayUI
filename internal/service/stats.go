package service

import (
	"context"
	"sync"
	"time"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/store"
)

// StatsService collects and manages traffic statistics.
type StatsService struct {
	StatsStore *store.StatsStore

	mu       sync.Mutex
	cancel   context.CancelFunc
	current  model.TrafficStats
	callback func(model.TrafficStats)
}

// SetCallback sets a function called on each stats update (for Wails event push).
func (s *StatsService) SetCallback(fn func(model.TrafficStats)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callback = fn
}

// StartCollecting begins periodic traffic collection.
func (s *StatsService) StartCollecting(profileID string, coreType model.ECoreType) {
	s.StopCollecting()

	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	go s.collectLoop(ctx, profileID, coreType)
}

// StopCollecting stops the collection goroutine.
func (s *StatsService) StopCollecting() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

// GetCurrentTraffic returns the last known traffic stats.
func (s *StatsService) GetCurrentTraffic() model.TrafficStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *StatsService) collectLoop(ctx context.Context, profileID string, coreType model.ECoreType) {
	ticker := time.NewTicker(2 * time.Second)
	saveTicker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	defer saveTicker.Stop()

	var lastUp, lastDown int64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			up, down := s.fetchTraffic(coreType)

			var deltaUp, deltaDown int64
			if lastUp > 0 || lastDown > 0 {
				deltaUp = up - lastUp
				deltaDown = down - lastDown
				if deltaUp < 0 {
					deltaUp = 0
				}
				if deltaDown < 0 {
					deltaDown = 0
				}
			}
			lastUp = up
			lastDown = down

			stats := model.TrafficStats{
				Up:   deltaUp / 2, // per second (2s interval)
				Down: deltaDown / 2,
			}

			s.mu.Lock()
			s.current = stats
			cb := s.callback
			s.mu.Unlock()

			if cb != nil {
				cb(stats)
			}

			if deltaUp > 0 || deltaDown > 0 {
				_ = s.StatsStore.UpdateStats(profileID, deltaUp, deltaDown)
			}

		case <-saveTicker.C:
			// Periodic save is handled by UpdateStats above.
		}
	}
}

// fetchTraffic queries the running core for cumulative traffic bytes.
func (s *StatsService) fetchTraffic(coreType model.ECoreType) (up, down int64) {
	// TODO: Implement actual API queries.
	// sing-box: GET http://127.0.0.1:9090/traffic (SSE)
	// xray: gRPC stats API
	return 0, 0
}
