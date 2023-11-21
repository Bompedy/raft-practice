package main

import (
	"net"
	"sync"
)

const (
	LEADER = iota
	CANDIDATE
	FOLLOWER
)

type Node struct {
	host    string
	state   int
	clients []Connection
	leader  net.Conn
	log     Log
}

func node(state int, parent string, nodes ...string) (*Node, error) {
	var mutex sync.Mutex
	node := &Node{
		state:   state,
		host:    parent,
		clients: make([]Connection, 0),
		log: Log{
			file:     [][]byte{},
			mutex:    &mutex,
			majority: (len(nodes)+1)/2 + 1,
		},
	}

	listener, err := net.Listen("tcp", parent)
	if err != nil {
		return node, err
	}

	var waiter sync.WaitGroup
	for _, address := range nodes {
		waiter.Add(1)
		address := address
		go func() {
			defer waiter.Done()
			client, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			buffer := make([]byte, 65535)
			err = Read(client, buffer[:1])
			if err != nil {
				return
			}
			if buffer[0] == 1 {
				node.leader = client
			}
			node.clients = append(node.clients, Connection{client, buffer})
		}()
	}

	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				continue
			}
			buffer := make([]byte, 1)
			if state == LEADER {
				buffer[0] = 1
			} else {
				buffer[0] = 0
			}
			err = Write(client, buffer)
			if err != nil {
				return
			}
			go node.replicator(client)
		}
	}()
	waiter.Wait()
	return node, nil
}
