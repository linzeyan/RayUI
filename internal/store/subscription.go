package store

import (
	"fmt"

	"github.com/RayUI/RayUI/internal/model"
)

// SubStore manages SubItem persistence.
type SubStore struct {
	store *Store[[]model.SubItem]
}

func NewSubStore() *SubStore {
	return &SubStore{
		store: NewStore[[]model.SubItem]("subscriptions.json", []model.SubItem{}),
	}
}

func (s *SubStore) GetAll() ([]model.SubItem, error) {
	return s.store.Load()
}

func (s *SubStore) GetByID(id string) (*model.SubItem, error) {
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

func (s *SubStore) Add(item model.SubItem) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	items = append(items, item)
	return s.store.Save(items)
}

func (s *SubStore) Update(item model.SubItem) error {
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
	return fmt.Errorf("subscription %s not found", item.ID)
}

func (s *SubStore) Delete(id string) error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}
	var remaining []model.SubItem
	for _, item := range items {
		if item.ID != id {
			remaining = append(remaining, item)
		}
	}
	return s.store.Save(remaining)
}
