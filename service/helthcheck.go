package service

import (
	"fmt"
	"net/http"
	"time"
	"zocket/distributed-kv-store/store"
)

func CheckHealth() {
	for {
		for _, node := range store.Nodes {
			_, err := http.Get(fmt.Sprintf("%s/health", node))
			if err != nil {
				fmt.Printf("Node %s failed! Rebalancing...\n", node)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
