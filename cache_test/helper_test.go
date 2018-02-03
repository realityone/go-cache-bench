package cache_test

import (
	"fmt"

	"github.com/google/uuid"

	"math/rand"
)

var keys []string

func loader(key interface{}) (interface{}, error) {
	return vValue(key), nil
}

func vValue(key interface{}) interface{} {
	return fmt.Sprintf("valueFor%s", key)
}

func randomKey() (init func(int), gen func() string) {
	init = func(_ int) {}
	gen = func() string {
		u, _ := uuid.NewRandom()
		return u.String()
	}
	return
}

func scheduledKey() (init func(int), gen func() string) {
	init = func(size int) {
		keys = make([]string, size)
		for i := 0; i < size; i++ {
			u, _ := uuid.NewRandom()
			keys[i] = u.String()
		}
	}
	gen = func() string {
		return keys[rand.Intn(len(keys))]
	}
	return
}
