package main

import (
	"encoding/binary"
	"net"
)

func Read(connection net.Conn, buffer []byte) error {
	for start := 0; start != len(buffer); {
		amount, reason := connection.Read(buffer[start:])
		if reason != nil {
			return reason
		}
		start += amount
	}
	return nil
}
func Write(connection net.Conn, buffer []byte) error {
	for start := 0; start != len(buffer); {
		amount, reason := connection.Write(buffer[start:])
		if reason != nil {
			return reason
		}
		start += amount
	}
	return nil
}

func ReadBytes(connection net.Conn, buffer []byte) ([]byte, error) {
	err := Read(connection, buffer[:4])
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(buffer[:4])
	err = Read(connection, buffer[4:4+length])
	if err != nil {
		return nil, err
	}
	return buffer[4 : 4+length], nil
}

func WriteBytes(connection net.Conn, buffer []byte, value []byte) error {
	var length = uint32(len(value))
	binary.LittleEndian.PutUint32(buffer[:4], length)
	copy(buffer[4:], value)
	err := Write(connection, buffer[:length+4])
	if err != nil {
		return err
	}
	return nil
}

type Connection struct {
	conn   net.Conn
	buffer []byte
}
