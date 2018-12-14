// Package cache ... do connect to redis with RedisConfig ref to common or other where?
// declare interfaces to use cahce in common
package cache

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_responseCache_Encode_Decode(t *testing.T) {
	cache := responseCache{
		Header: http.Header{
			"Content-Type": []string{"appliction/json", "xhtml"},
			"X-Real-Ip":    []string{"127.0.0.1"},
		},
		Status: 200,
		Data:   []byte("this is body bytes"),
	}

	byts, err := encodeCache(&cache)
	if err != nil {
		t.Errorf("could encode cache: %v", err)
	}

	// logf got bytes
	// t.Logf("got encode string: %s", string(byts))

	if got, err := decodeToCache(byts); err != nil {
		t.Errorf("could encode cache: %v", err)
	} else if !reflect.DeepEqual(got, cache) {
		t.Errorf("could not decode in same way: want %v, got %v", cache, got)
	} else {
		t.Logf("status: %d, data: %s", got.Status, string(got.Data))
	}
}

func Test_cachedWriter(t *testing.T) {
	var _ http.ResponseWriter = cachedWriter{}
}
