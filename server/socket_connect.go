package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

/*
	ver : version
	cmd : connection method
	rev : not use, reserved filed
	atyp: addr type (ipv4 = 1, ipv6 = 2, domain name 3
	dst.addr: destination addr
	dst.port: destination port (2 byte)
*/

func SocketConnect(conn net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)
	n, err := io.ReadFull(conn, buf[:4])
	if n != 4 {
		return nil, errors.New("read header error" + err.Error())
	}
	ver, cmd, _, atyp := int(buf[0]), int(buf[1]), int(buf[2]), int(buf[3])
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid version or connection type")
	}
	addr := ""

	switch atyp {
	case 1:
		n, err := io.ReadFull(conn, buf[:4])
		if n != 4 {
			return nil, errors.New("ipv4 connection error" + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", int(buf[0]), int(buf[1]), int(buf[2]), int(buf[3]))

	case 2:
		n, err := io.ReadFull(conn, buf[:16])
		if n != 16 {
			return nil, errors.New("ipv6 connection error" + err.Error())
		}
		for i := 0; i < 16; i += 2 {
			addr += strconv.FormatUint(uint64(buf[i])<<8+uint64(buf[i+1]), 16) + ":"
		}
		addr = addr[:len(addr)-1]
	case 3:
		// not implement domain name resolve
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	default:
		return nil, errors.New("invalid atyp")

	}

	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 {
		return nil, errors.New("resolve port error" + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	dest, err := net.Dial("tcp", addr+":"+strconv.Itoa(int(port)))
	if err != nil {
		return nil, errors.New("dial error " + err.Error())
	}
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write response error" + err.Error())
	}
	fmt.Println("connect success")
	return dest, err

}
