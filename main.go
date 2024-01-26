package main

import (
	"fmt"
	"log"
	"net"
	"socket5_feature/server"
)

func main() {
	listner, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return
	}
	for {
		client, err := listner.Accept()
		if err != nil {
			log.Println("access fatal")
			continue
		}
		fmt.Println("start process")
		go process(client)
	}
}
func process(client net.Conn) {
	err := server.SocketAuth(client)
	if err != nil {
		client.Close()
		return
	}
	target, err := server.SocketConnect(client)
	if err != nil {
		client.Close()
		return
	}
	server.SocketForward(client, target)

}
