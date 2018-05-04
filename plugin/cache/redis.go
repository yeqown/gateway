// Package cache ... do connect to redis with RedisConfig ref to common or other where?
// // declare interfaces to use cahce in common
package cache

// type RedisConfig struct {
// 	Host string `json:"host"`
// 	Port int    `json:"port"`
// }

// func ConnRedis(rc *RedisConfig) {

// }

// type KeyVal interface {
// 	Key() string
// 	Expire() uint
// 	Marshal() string
// 	Unmarshal(string) error
// }

// func IsKeyExisted(kv KeyVal) bool {
// 	return false
// }

// func SaveKV(kv KeyVal) {
// 	key := kv.Key()
// 	val := kv.Marshal()
// 	expire := kv.Expire()
// 	// do save
// }

// func GetKV(kv KeyVal) error {
// 	key := kv.Key()
// 	// get string from redis
// 	s := ""

// 	return kv.Unmarshal(s)
// }
