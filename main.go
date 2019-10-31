package main

import (
	"Chating/chat"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	chat.RegisterChatServer(s, NewChatService("chat.db", 3*time.Hour, 10, 3*time.Minute))
	s.Serve(lis)
}
