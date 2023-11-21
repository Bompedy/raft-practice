package main

import (
	"net"
	"sync"
)

type Log struct {
	file            [][]byte
	mutex           *sync.Mutex
	majority, index int
}

func (node *Node) broadcast(message []byte) {
	println("Leader Got: ", string(message))
	//defer node.log.mutex.Lock()
	//node.log.mutex.Lock()
	for _, connection := range node.clients {
		connection := connection
		go func() {
			err := WriteBytes(connection.conn, connection.buffer, message)
			if err != nil {
				println("Error writing to clients")
			}
		}()
	}
}

func (node *Node) replicator(server net.Conn) {
	buffer := make([]byte, 65535)
	for {
		if node.state == LEADER {
			message, err := ReadBytes(server, buffer)
			if err != nil {
				continue
			}
			node.broadcast(message)
		} else {
			message, err := ReadBytes(server, buffer)
			if err != nil {
				continue
			}
			// add to log here
			println("Client Got: ", string(message))
		}
	}
}

func (node *Node) append(value string) {
	if node.state == LEADER {
		node.broadcast([]byte(value))
	} else {
		err := WriteBytes(node.leader, make([]byte, uint32(len(value))+4), []byte(value))
		if err != nil {
			return
		}
	}
}
