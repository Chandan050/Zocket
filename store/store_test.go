package store

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

// Mock server to simulate other nodes for replication
func startMockNode(t *testing.T, wg *sync.WaitGroup) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	wg.Done()
	return server
}

func TestStore_PostGet(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	// Start mock nodes
	node1 := startMockNode(t, &wg)
	node2 := startMockNode(t, &wg)
	node3 := startMockNode(t, &wg)
	wg.Wait()

	// Override Nodes with mock server URLs
	Nodes = []string{node1.URL, node2.URL, node3.URL}

	store := NewStore()

	// Test Post and Get
	err := store.Post("key1", "value1")
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}

	val, ok := store.Get("key1")
	if !ok || val != "value1" {
		t.Fatalf("Get failed: expected value1, got %v", val)
	}

	// Test replication failure scenario by shutting down one node
	node1.Close()

	err = store.Post("key2", "value2")
	if err != nil {
		t.Fatalf("Post with one node down failed: %v", err)
	}

	val, ok = store.Get("key2")
	if !ok || val != "value2" {
		t.Fatalf("Get after replication failed: expected value2, got %v", val)
	}

	// Cleanup
	node2.Close()
	node3.Close()
}

func TestStore_ConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	store := NewStore()

	// Override Nodes with empty slice to avoid replication during concurrency test
	Nodes = []string{}

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			err := store.Post(key, value)
			if err != nil {
				t.Errorf("Post failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			_, _ = store.Get(key)
		}(i)
	}

	wg.Wait()
}
