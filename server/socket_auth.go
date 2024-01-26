package server

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func SocketAuth(conn net.Conn) error {
	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return errors.New("read connection header error" + err.Error())
	}
	ver, nMethod := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}
	n, err = io.ReadFull(conn, buf[:nMethod])
	if n != nMethod {
		return errors.New("reading method" + err.Error())
	}

	n, err = conn.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("failure response" + err.Error())
	}
	fmt.Println("auth access ...")
	return nil
}
