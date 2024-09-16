package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/marwanelmessiry/ChatRoomGrpc/proto"
)

type server struct {
	pb.UnimplementedChatAppServer
	db      *sql.DB
	clients map[string]chan *pb.Message
	mu      sync.Mutex
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	timestamp, err := ptypes.Timestamp(req.Message.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert timestamp: %v", err)
	}

	_, err = s.db.Exec("INSERT INTO messages (sender, recipient, content, timestamp) VALUES (?, ?, ?, ?)",
		req.Message.Sender, req.Message.Recipient, req.Message.Content, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to insert message into database: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Message.Recipient == "" {
		for _, ch := range s.clients {
			ch <- req.Message
		}
	} else {
		if ch, ok := s.clients[req.Message.Recipient]; ok {
			ch <- req.Message
		}
	}

	return &pb.SendMessageResponse{}, nil
}

func (s *server) ReceiveMessages(req *pb.ReceiveMessagesRequest, stream pb.ChatApp_ReceiveMessagesServer) error {
	s.mu.Lock()
	if s.clients == nil {
		s.clients = make(map[string]chan *pb.Message)
	}
	ch := make(chan *pb.Message)
	s.clients[req.User] = ch
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, req.User)
		s.mu.Unlock()
		close(ch)
	}()

	for msg := range ch {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
