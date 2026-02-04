package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/parser"
	"github.com/RayUI/RayUI/internal/store"
	"github.com/RayUI/RayUI/internal/util"
)

// SubscriptionService handles subscription fetching, parsing and syncing.
type SubscriptionService struct {
	SubStore     *store.SubStore
	ProfileStore *store.ProfileStore
	StatsStore   *store.StatsStore

	// Auto-update
	stopCh     chan struct{}
	stopOnce   sync.Once
	OnAutoSync func(subID string, count int, err error)
}

// Sync fetches and syncs a single subscription, returning the number of imported profiles.
func (s *SubscriptionService) Sync(subID string) (int, error) {
	sub, err := s.SubStore.GetByID(subID)
	if err != nil {
		return 0, err
	}
	if sub == nil {
		return 0, fmt.Errorf("subscription %s not found", subID)
	}

	// Fetch subscription content.
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", sub.URL, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	ua := sub.UserAgent
	if ua == "" {
		ua = "RayUI/1.0"
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch subscription: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response: %w", err)
	}

	// Parse profiles.
	items, err := parser.ParseBatch(string(body))
	if err != nil {
		return 0, fmt.Errorf("parse subscription: %w", err)
	}

	// Apply filter.
	if sub.Filter != "" {
		re, err := regexp.Compile(sub.Filter)
		if err == nil {
			var filtered []model.ProfileItem
			for _, item := range items {
				if re.MatchString(item.Remarks) {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}
	}

	// Load old profiles for merge.
	oldProfiles, _ := s.ProfileStore.GetBySubID(subID)
	oldByURI := make(map[string]model.ProfileItem)
	for _, p := range oldProfiles {
		if p.ShareURI != "" {
			oldByURI[p.ShareURI] = p
		}
	}

	// Prepare new profiles.
	for i := range items {
		items[i].SubID = subID
		// Try to match existing profile by ShareURI to preserve ID.
		if old, ok := oldByURI[items[i].ShareURI]; ok {
			items[i].ID = old.ID
		} else {
			items[i].ID = util.GenerateUUID()
		}
		items[i].Sort = i
	}

	// Replace profiles for this subscription.
	if err := s.ProfileStore.ReplaceBySubID(subID, items); err != nil {
		return 0, fmt.Errorf("save profiles: %w", err)
	}

	// Update subscription timestamp.
	sub.UpdateTime = time.Now().Unix()
	if err := s.SubStore.Update(*sub); err != nil {
		return len(items), fmt.Errorf("update subscription time: %w", err)
	}

	return len(items), nil
}

// SyncAll syncs all enabled subscriptions.
func (s *SubscriptionService) SyncAll() (map[string]int, error) {
	subs, err := s.SubStore.GetAll()
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)
	for _, sub := range subs {
		if !sub.Enabled {
			continue
		}
		count, err := s.Sync(sub.ID)
		if err != nil {
			results[sub.ID] = -1
			continue
		}
		results[sub.ID] = count
	}
	return results, nil
}

// StartAutoUpdate starts a background goroutine that periodically checks
// subscriptions and syncs those whose auto-update interval has elapsed.
func (s *SubscriptionService) StartAutoUpdate() {
	s.stopCh = make(chan struct{})
	go s.autoUpdateLoop()
}

// StopAutoUpdate stops the background auto-update goroutine.
func (s *SubscriptionService) StopAutoUpdate() {
	s.stopOnce.Do(func() {
		if s.stopCh != nil {
			close(s.stopCh)
		}
	})
}

func (s *SubscriptionService) autoUpdateLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkAndSync()
		}
	}
}

func (s *SubscriptionService) checkAndSync() {
	subs, err := s.SubStore.GetAll()
	if err != nil {
		return
	}

	now := time.Now().Unix()
	for _, sub := range subs {
		if !sub.Enabled || sub.AutoUpdateInterval <= 0 {
			continue
		}
		intervalSec := int64(sub.AutoUpdateInterval) * 60
		if now-sub.UpdateTime < intervalSec {
			continue
		}
		count, err := s.Sync(sub.ID)
		if s.OnAutoSync != nil {
			s.OnAutoSync(sub.ID, count, err)
		}
	}
}
