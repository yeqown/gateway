package presistence

import (
	"testing"
	"time"
)

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	if store == nil {
		t.Errorf("could not new InMemoryStore")
	}
}

func TestInMemoryStore_Set_Get(t *testing.T) {
	store := NewInMemoryStore()
	key := "key"
	value := "value"
	_ = store.Set(key, []byte(value), time.Second*3)
	// exists
	if b := store.Exists(key); !b {
		t.Errorf("could not found key: %s", key)
	}

	byts, err := store.Get(key)
	if err != nil {
		t.Errorf("could not Get from store: %v", err)
	}
	if string(byts) != value {
		t.Errorf("InMemoryStore.Get from store: want %v, get %v", value, string(byts))
	}

	// replace ...
	newVal := "new value"
	if err := store.Replace(key, []byte(newVal), time.Second*3); err != nil {
		t.Errorf("could not Replace to store: %v", err)
	}
	byts, err = store.Get(key)
	if err != nil {
		t.Errorf("could not Get from store: %v", err)
	}
	if string(byts) != newVal {
		t.Errorf("InMemoryStore.Get from store: want %v, get %v", newVal, string(byts))
	}
}

func TestInMemoryStore_Flush(t *testing.T) {
	store := NewInMemoryStore()
	value := "value value 1"
	store.Set("key1", []byte(value), time.Second*2)
	store.Set("key2", []byte(value), time.Second*2)
	store.Set("key3", []byte(value), time.Second*2)

	if !store.Exists("key1") || !store.Exists("key2") || !store.Exists("key3") {
		t.Error("keys not exists")
	}

	// call flush
	store.Flush()

	if store.Exists("key1") || store.Exists("key2") || store.Exists("key3") {
		t.Error("some key exists")
	}
}
