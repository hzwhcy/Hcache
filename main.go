package main

import (
	"HCache/hcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	hcache.NewGroup("scores", hcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)

	addr := "localhost:9999"
	peers := hcache.NewHttpPool(addr)
	log.Println("hcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
