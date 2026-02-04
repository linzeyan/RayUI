package store

import (
	"fmt"

	"github.com/RayUI/RayUI/internal/model"
)

// ProfileStore manages ProfileItem persistence.
type ProfileStore struct {
	store *Store[[]model.ProfileItem]
}

// NewProfileStore creates a ProfileStore backed by profiles.json.
func NewProfileStore() *ProfileStore {
	return &ProfileStore{
		store: NewStore[[]model.ProfileItem]("profiles.json", []model.ProfileItem{}),
	}
}

func (s *ProfileStore) GetAll() ([]model.ProfileItem, error) {
	return s.store.Load()
}

func (s *ProfileStore) GetByID(id string) (*model.ProfileItem, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].ID == id {
			return &items[i], nil
		}
	}
	return nil, nil
}

func (s *ProfileStore) GetBySubID(subID string) ([]model.ProfileItem, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	var result []model.ProfileItem
	for _, item := range items {
		if item.SubID == subID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *ProfileStore) Add(item model.ProfileItem) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	items = append(items, item)
	return s.store.Save(items)
}

func (s *ProfileStore) Update(item model.ProfileItem) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	for i := range items {
		if items[i].ID == item.ID {
			items[i] = item
			return s.store.Save(items)
		}
	}
	return fmt.Errorf("profile %s not found", item.ID)
}

func (s *ProfileStore) Delete(ids []string) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}
	var remaining []model.ProfileItem
	for _, item := range items {
		if _, del := idSet[item.ID]; !del {
			remaining = append(remaining, item)
		}
	}
	return s.store.Save(remaining)
}

// ReplaceBySubID replaces all profiles belonging to a subscription.
func (s *ProfileStore) ReplaceBySubID(subID string, newItems []model.ProfileItem) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	var remaining []model.ProfileItem
	for _, item := range items {
		if item.SubID != subID {
			remaining = append(remaining, item)
		}
	}
	remaining = append(remaining, newItems...)
	return s.store.Save(remaining)
}
