// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zocket/distributed-kv-store/hashing"
	"zocket/distributed-kv-store/store"
)

var newHash *hashing.HashRing
var kvStore *store.Store

func main() {
	newHash = hashing.NewHashRing(3)
	newHash.AddNode("node1")
	newHash.AddNode("node2")
	newHash.AddNode("node3")

	kvStore = store.NewStore()

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	value, exists := kvStore.Get(key)

	if exists {
		resp := map[string]string{"value": value}
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Key not found", http.StatusNotFound)
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	node := newHash.GetNode(key)

	if node == "" {
		http.Error(w, "No node found for key", http.StatusInternalServerError)
		return
	}

	kvStore.Post(key, value)
	fmt.Fprintf(w, "Stored key=%s, value=%s on node %s\n", key, value, node)
}
