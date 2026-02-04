package store

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/RayUI/RayUI/internal/util"
)

// Store[T] is a generic JSON file-backed store.
type Store[T any] struct {
	filename   string
	defaultVal T
	mu         sync.RWMutex
}

// NewStore creates a Store that persists to ~/.RayUI/{filename}.
func NewStore[T any](filename string, defaultVal T) *Store[T] {
	return &Store[T]{
		filename:   filename,
		defaultVal: defaultVal,
	}
}

// GetPath returns the full file path for this store.
func (s *Store[T]) GetPath() string {
	return filepath.Join(util.AppDataDir(), s.filename)
}

// Load reads the JSON file into T. If the file does not exist, it writes
// the default value and returns it.
func (s *Store[T]) Load() (T, error) {
	s.mu.RLock()
	path := s.GetPath()
	s.mu.RUnlock()

	var val T
	if err := util.ReadJSON(path, &val); err != nil {
		if os.IsNotExist(err) {
			val = s.defaultVal
			if saveErr := s.Save(val); saveErr != nil {
				return val, saveErr
			}
			return val, nil
		}
		return val, err
	}
	return val, nil
}

// Save writes data to the JSON file atomically.
func (s *Store[T]) Save(data T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return util.AtomicWriteJSON(s.GetPath(), data)
}
