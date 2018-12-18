package api

import "github.com/yeqown/gateway/config/presistence"

var (
	global presistence.Store
)

// Global ...
func Global() presistence.Store {
	return global
}

// SetGlobal ..
func SetGlobal(store presistence.Store) {
	global = store
}
