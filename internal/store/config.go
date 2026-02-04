package store

import "github.com/RayUI/RayUI/internal/model"

// ConfigStore manages the global Config persistence.
type ConfigStore struct {
	store *Store[model.Config]
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		store: NewStore[model.Config]("config.json", model.DefaultConfig()),
	}
}

func (s *ConfigStore) Load() (model.Config, error) {
	return s.store.Load()
}

func (s *ConfigStore) Save(cfg model.Config) error {
	return s.store.Save(cfg)
}
