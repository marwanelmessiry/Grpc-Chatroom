package main

import (
	"log"
	"net"
	"sync"

	pb "github.com/marwanelmessiry/ChatRoomGrpc/proto"
	"google.golang.org/grpc"
)

func main() {
	db, err := InitDB("messages.db")
	if err != nil {
		log.Fatalf("error in initializing db: %v", err)
	}
	log.Println("db is initialized")
	defer db.Close()

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("error in listening to port 8000: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatAppServer(s, &server{
		UnimplementedChatAppServer: pb.UnimplementedChatAppServer{},
		db:                         db,
		clients:                    make(map[string]chan *pb.Message),
		mu:                         sync.Mutex{},
	})
	log.Println("server is running on port 8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error in running server: %v", err)
	}
}
