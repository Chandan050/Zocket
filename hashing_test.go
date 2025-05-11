package main_test

import (
	"fmt"
	"zocket/distributed-kv-store/hashing"
)

func testMain() {
	hr := hashing.NewHashRing(3)
	hr.AddNode("node1")
	hr.AddNode("node2")
	hr.AddNode("node3")
	fmt.Println(hr.GetNode("testkey"))
}
