package main

import (
	"fmt"
)

func main() {
	go func() {
		node, err := node(LEADER, "localhost:8080", "localhost:8081", "localhost:8082")
		if err != nil {
			fmt.Printf("Error creating node on %s: %v", node.host, err)
			return
		}
		node.append("leader")
	}()
	go func() {
		node, err := node(FOLLOWER, "localhost:8081", "localhost:8080", "localhost:8082")
		if err != nil {
			fmt.Printf("Error creating node on %s: %v", node.host, err)
			return
		}
		node.append("node 1")
	}()
	go func() {
		node, err := node(FOLLOWER, "localhost:8082", "localhost:8081", "localhost:8080")
		if err != nil {
			fmt.Printf("Error creating node on %s: %v", node.host, err)
			return
		}
		node.append("node 2")
	}()
	select {}
}
