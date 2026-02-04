package store

import "github.com/RayUI/RayUI/internal/model"

// DNSStore manages DNSItem persistence.
type DNSStore struct {
	store *Store[model.DNSItem]
}

func NewDNSStore() *DNSStore {
	return &DNSStore{
		store: NewStore[model.DNSItem]("dns.json", model.DefaultDNSItem()),
	}
}

func (s *DNSStore) Load() (model.DNSItem, error) {
	return s.store.Load()
}

func (s *DNSStore) Save(item model.DNSItem) error {
	return s.store.Save(item)
}
