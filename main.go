package main

import (
	"Aryan-Kuntal/tcp_chat/chat"
	"log"
	"net"
)

func main() {
	server := chat.NewServer()
	go server.Run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("Error occured in server", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error occured while accepting", err.Error())
		}

		go server.NewClient(conn)
	}
}
