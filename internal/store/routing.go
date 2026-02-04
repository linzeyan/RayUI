package store

import (
	"fmt"

	"github.com/RayUI/RayUI/internal/model"
)

// RoutingStore manages RoutingItem persistence.
// On first load (file missing), it seeds with DefaultRoutingItems.
type RoutingStore struct {
	store *Store[[]model.RoutingItem]
}

func NewRoutingStore() *RoutingStore {
	return &RoutingStore{
		store: NewStore[[]model.RoutingItem]("routing.json", model.DefaultRoutingItems()),
	}
}

func (s *RoutingStore) GetAll() ([]model.RoutingItem, error) {
	return s.store.Load()
}

func (s *RoutingStore) GetByID(id string) (*model.RoutingItem, error) {
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

func (s *RoutingStore) Add(item model.RoutingItem) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	items = append(items, item)
	return s.store.Save(items)
}

func (s *RoutingStore) Update(item model.RoutingItem) error {
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
	return fmt.Errorf("routing item %s not found", item.ID)
}

func (s *RoutingStore) Delete(id string) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	var remaining []model.RoutingItem
	for _, item := range items {
		if item.ID != id {
			if item.Locked {
				remaining = append(remaining, item)
				continue
			}
			remaining = append(remaining, item)
		} else if item.Locked {
			return fmt.Errorf("cannot delete built-in routing item %q", item.Remarks)
		}
	}
	return s.store.Save(remaining)
}
