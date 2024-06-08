package core

import (
	"log"
	"time"
)

func expireSample() float32 {
	limit := 20
	expireCount := 0

	for key, obj := range store {
		if obj.ExpiredAt != -1 {
			limit--

			// delete if key is expired
			if obj.ExpiredAt <= time.Now().UnixMilli() {
				delete(store, key)
				expireCount++
			}
		}

		if limit == 20 {
			break
		}
	}

	return float32(expireCount) / float32(20)
}

func DeleteExpireKeys() {
	for {
		fraction := expireSample()

		if fraction < 0.25 {
			break
		}
	}
	log.Println("expired keys deleted, total keys ", len(store))
}
