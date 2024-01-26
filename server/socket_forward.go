package server

import (
	"fmt"
	"io"
	"net"
)

func SocketForward(client, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		_, err := io.Copy(src, dest)
		if err != nil {
		}
		fmt.Println("forward success")
	}
	go forward(client, target)
	go forward(target, client)
}
