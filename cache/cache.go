package cache

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"

	"gitee.com/ikongjix/go_common/redis_db/master_db"
)

// GetCachedStruct Use this method to get cached struct
func GetCachedStruct(cacheKey string, destStruct interface{}) (err error) {
	res := master_db.MasterRedis.Get(cacheKey)
	bs, err := res.Bytes()
	if err != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewReader(bs))
	return decoder.Decode(destStruct)
}

// CacheStruct Use this method to cache structs
func CacheStruct(cacheKey string, destStruct interface{}, expire time.Duration) (ok bool, err error) {
	var structBuff bytes.Buffer
	encoder := gob.NewEncoder(&structBuff)
	encoder.Encode(destStruct)

	res := master_db.MasterRedis.Set(cacheKey, structBuff.String(), 0)
	v, err := res.Result()
	if err != nil {
		return
	}

	ok = strings.ToLower(v) == "ok"

	return
}

func MGetCachedStructs(keys []string, destStructs []interface{}) ([]int, error) {
	keySize := len(keys)
	var err error
	var unCachedIndex []int

	// for _, key := range keys {
	// 	err = conn.Send("GET", key)
	// 	if err != nil {
	// 		return unCachedIndex, mGetCachedStructsSendError
	// 	}
	// }

	// if err = conn.Flush(); err != nil {
	// 	return unCachedIndex, mGetCachedStructsFlushError
	// }

	// for i := 0; i < keySize; i++ {
	// 	destStruct := destStructs[i]
	// 	v, err := redis.Bytes(conn.Receive())
	// 	if err == nil {
	// 		decoder := gob.NewDecoder(bytes.NewReader(v))
	// 		derr := decoder.Decode(destStruct)
	// 		if derr != nil {
	// 			destStructs[i] = nil
	// 			unCachedIndex = append(unCachedIndex, i)
	// 		}
	// 	} else {
	// 		destStructs[i] = nil
	// 		unCachedIndex = append(unCachedIndex, i)
	// 	}

	// }
	unCachedIndex = append(unCachedIndex, 0, 1)
	return unCachedIndex, nil
}
