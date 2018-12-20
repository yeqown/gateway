package presistence

import (
	"time"
)

const (
	// DefaultExpire 0 second expire duration
	DefaultExpire = time.Duration(0)
	// ForeverExpire forever expire duration
	ForeverExpire = time.Duration(-1)
)

// Store interface to presist data into database or file etc.
type Store interface {

	// Set func to Set item with params
	Set(key string, value []byte, expire time.Duration) error

	// Get func to Get item with params
	Get(key string) ([]byte, error)

	// Replace func to Replace item with params
	Replace(key string, newVal []byte, expire time.Duration) error

	// Exists func to Exists item with params
	Exists(key string) bool

	// Delete func to Delete item with params
	Delete(key string) error

	// Flush func to flush all data
	Flush() error
}
