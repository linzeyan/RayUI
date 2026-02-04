package store

import (
	"time"

	"github.com/RayUI/RayUI/internal/model"
)

// StatsStore manages per-server traffic statistics.
type StatsStore struct {
	store *Store[[]model.ServerStatItem]
}

func NewStatsStore() *StatsStore {
	return &StatsStore{
		store: NewStore[[]model.ServerStatItem]("stats.json", []model.ServerStatItem{}),
	}
}

func (s *StatsStore) GetAll() ([]model.ServerStatItem, error) {
	return s.store.Load()
}

func (s *StatsStore) GetByProfileID(profileID string) (*model.ServerStatItem, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].ProfileID == profileID {
			return &items[i], nil
		}
	}
	return nil, nil
}

// UpdateStats adds the given delta bytes to a profile's statistics.
// Creates the record if it doesn't exist.
func (s *StatsStore) UpdateStats(profileID string, deltaUp, deltaDown int64) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	found := false
	for i := range items {
		if items[i].ProfileID == profileID {
			found = true
			// Reset daily counters on date change.
			if items[i].DateNow != today {
				items[i].TodayUp = 0
				items[i].TodayDown = 0
				items[i].DateNow = today
			}
			items[i].TotalUp += deltaUp
			items[i].TotalDown += deltaDown
			items[i].TodayUp += deltaUp
			items[i].TodayDown += deltaDown
			items[i].LastUpdate = now.Unix()
			break
		}
	}

	if !found {
		items = append(items, model.ServerStatItem{
			ProfileID:  profileID,
			TotalUp:    deltaUp,
			TotalDown:  deltaDown,
			TodayUp:    deltaUp,
			TodayDown:  deltaDown,
			DateNow:    today,
			LastUpdate: now.Unix(),
		})
	}

	return s.store.Save(items)
}

// DeleteByProfileID removes stats for the given profile.
func (s *StatsStore) DeleteByProfileID(profileID string) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	var remaining []model.ServerStatItem
	for _, item := range items {
		if item.ProfileID != profileID {
			remaining = append(remaining, item)
		}
	}
	return s.store.Save(remaining)
}

// Clear removes all statistics.
func (s *StatsStore) Clear() error {
	return s.store.Save([]model.ServerStatItem{})
}

// ResetDaily resets todayUp/todayDown for all entries whose date has changed.
func (s *StatsStore) ResetDaily() error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	changed := false
	for i := range items {
		if items[i].DateNow != today {
			items[i].TodayUp = 0
			items[i].TodayDown = 0
			items[i].DateNow = today
			changed = true
		}
	}

	if changed {
		return s.store.Save(items)
	}
	return nil
}
