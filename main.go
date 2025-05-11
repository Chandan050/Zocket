// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"zocket/distributed-kv-store/hashing"
	"zocket/distributed-kv-store/service"
	"zocket/distributed-kv-store/store"
)

var newHash *hashing.HashRing
var kvStore *store.Store
var currentNode string

func main() {
	currentNode = os.Getenv("NODE_NAME")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if currentNode == "" {
		fmt.Println("NODE_NAME environment variable not set")
		return
	}

	newHash = hashing.NewHashRing(3)
	newHash.AddNode("node1")
	newHash.AddNode("node2")
	newHash.AddNode("node3")

	kvStore = store.NewStore()

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	fmt.Printf("Server running on port 8080 as %s...\n", currentNode)
	go service.CheckHealth()
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	if err := http.ListenAndServe(fmt.Sprintf(":"+port), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func forwardGetRequest(node, key string) (string, int, error) {
	url := fmt.Sprintf("http://%s/get?key=%s", node, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	return string(body), resp.StatusCode, nil
}

func forwardSetRequest(node, key, value string) (string, int, error) {
	url := fmt.Sprintf("http://%s/set?key=%s&value=%s", node, key, value)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	return string(body), resp.StatusCode, nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	node := newHash.GetNode(key)
	if node == "" {
		http.Error(w, "No node found for key", http.StatusInternalServerError)
		return
	}

	if node != currentNode {
		body, status, err := forwardGetRequest(node, key)
		if err != nil {
			http.Error(w, "Failed to forward get request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}
	fmt.Printf("Fetching key=%s from node %s\n", key, node)

	value, exists := kvStore.Get(key)
	if exists {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{"value": value}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
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

	if node != currentNode {
		body, status, err := forwardSetRequest(node, key, value)
		if err != nil {
			http.Error(w, "Failed to forward set request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	err := kvStore.Post(key, value)
	if err != nil {
		http.Error(w, "Failed to store key-value: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Stored key=%s, value=%s on node %s\n", key, value, node)
}
