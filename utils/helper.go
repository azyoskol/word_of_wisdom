package utils

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func FatalApplication(msg string, err error) {
	log.Fatalf("%s > %s\n", msg, err)
}

func ReadMessage(conn net.Conn) (msg []byte, err error) {
	var length uint64
	if err = binary.Read(conn, binary.BigEndian, &length); err != nil {
		return
	}

	msg = make([]byte, length)
	_, err = io.ReadFull(conn, msg)
	return
}

func WriteMessage(conn net.Conn, msg []byte) (err error) {
	if err = binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return
	}
	_, err = conn.Write(msg)
	return
}
