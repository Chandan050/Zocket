package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"zocket/distributed-kv-store/models"
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

var Nodes = []string{}

// replicate sends the key-value pair to other nodes and waits for quorum acknowledgments
func replicate(key, value string) error {
	if len(Nodes) == 0 {
		// No other nodes to replicate to, consider success
		return nil
	}

	post := models.Request{Key: key, Value: value}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	if err != nil {
		return err
	}

	quorum := len(Nodes)/2 + 1
	successCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, node := range Nodes {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Post(fmt.Sprintf("%s/set", n), "application/json", &buf)
			if err == nil && resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(node)
	}

	wg.Wait()

	if successCount >= quorum {
		return nil
	}
	return errors.New("failed to replicate to quorum of nodes")
}

// Post stores the data into the database with strong consistency
func (s *Store) Post(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	err := replicate(key, value)
	if err != nil {
		// rollback or handle failure accordingly
		delete(s.data, key)
		return err
	}
	return nil
}

// Get retrieves the value by key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}
