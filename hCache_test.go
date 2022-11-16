package HCache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup_Get(t *testing.T) {
	counts := make(map[string]int, len(db))
	hee := NewGroup("scores", GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("search key", key)
			if v, ok := db[key]; ok {
				if _, ok := counts[key]; !ok {
					counts[key] = 0
				}
				counts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)
	for k, v := range db {
		if val, err := hee.Get(k); err != nil || val.String() != v {
			t.Fatal("fail to get value")
		}
		if _, err := hee.Get(k); err != nil || counts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if v, err := hee.Get("empty"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", v)
	}
}
