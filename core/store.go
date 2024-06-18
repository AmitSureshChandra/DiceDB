package core

import (
	"dicedb/config"
	"time"
)

type Obj struct {
	Value     interface{}
	ExpiredAt int64
}

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64) *Obj {

	var expiredAt int64 = -1

	if durationMs > 0 {
		expiredAt = time.Now().UnixMilli() + +durationMs
	}

	return &Obj{
		Value:     value,
		ExpiredAt: expiredAt,
	}
}

func Put(k string, obj *Obj) {
	if len(store) >= config.KeysLimit {
		evict()
	}

	store[k] = obj
}

func Get(k string) *Obj {
	if val, ok := store[k]; ok {
		return val
	}
	return nil
}

func Del(key string) bool {
	if _, ok := store[key]; ok {
		delete(store, key)
		return true
	}
	return false
}
