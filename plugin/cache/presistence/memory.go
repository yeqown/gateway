package presistence

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// NewInMemoryStore ...
func NewInMemoryStore() *InMemoryStore {
	defaultExpiration := 5 * time.Minute
	cleanupInterval := time.Duration(0) // never clean up
	return &InMemoryStore{Cache: cache.New(defaultExpiration, cleanupInterval)}
}

// InMemoryStore to save cache into memory can be read and write
type InMemoryStore struct {
	Cache *cache.Cache
}

// Set func implement presistence.Store interface
func (s *InMemoryStore) Set(key string, value []byte, expire time.Duration) error {
	s.Cache.Set(key, value, expire)
	return nil
}

// Get func implement presistence.Store interface
func (s *InMemoryStore) Get(key string) ([]byte, error) {
	v, ok := s.Cache.Get(key)
	if !ok {
		return nil, fmt.Errorf("could not found key: %s", key)
	}
	byts, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("could not assert to bytes by key: %s", key)
	}

	return byts, nil
}

// Replace func implement presistence.Store interface
func (s *InMemoryStore) Replace(
	key string, newVal []byte, expire time.Duration) error {
	if err := s.Cache.Replace(key, newVal, expire); err != nil {
		return fmt.Errorf("could not replace cache: %v", err)
	}
	return nil
}

// Exists func implement presistence.Store interface
func (s *InMemoryStore) Exists(key string) bool {
	_, ok := s.Cache.Get(key)
	return ok
}

// Delete func implement presistence.Store interface
func (s *InMemoryStore) Delete(key string) error {
	s.Cache.Delete(key)
	return nil
}

// Flush func implement presistence.Store interface
func (s *InMemoryStore) Flush() error {
	s.Cache.Flush()
	return nil
}
