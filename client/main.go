package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/marwanelmessiry/ChatRoomGrpc/proto"
	"google.golang.org/grpc"
)

func main() {
	// Define command-line flags for sender and recipient
	sender := flag.String("sender", "", "Sender of the message")
	recipient := flag.String("recipient", "", "Recipient of the message (optional)")
	flag.Parse()

	if *sender == "" {
		log.Fatalf("sender must be specified")
	}

	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatAppClient(conn)

	go receiveMessages(client, *sender)

	reader := bufio.NewReader(os.Stdin)
	for {
		//fmt.Print("Enter message: ")
		content, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("could not read input: %v", err)
		}

		content = strings.TrimSpace(content)

		sendMessage(client, *sender, *recipient, content)
	}
}

func sendMessage(client pb.ChatAppClient, sender, recipient, content string) {
	timestamp, _ := ptypes.TimestampProto(time.Now())
	message := &pb.Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Timestamp: timestamp,
	}

	req := &pb.SendMessageRequest{
		Message: message,
	}

	_, err := client.SendMessage(context.Background(), req)
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	log.Printf(" %s: %s", sender, content)
}

func receiveMessages(client pb.ChatAppClient, user string) {
	req := &pb.ReceiveMessagesRequest{
		User: user,
	}

	stream, err := client.ReceiveMessages(context.Background(), req)
	if err != nil {
		log.Fatalf("could not receive messages: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
		}

		log.Printf(" %s: %s", msg.Sender, msg.Content)
	}
}
