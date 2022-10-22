package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "proto/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userName   = flag.String("username", "anonymous", "The name others will see you by")
	serverPort = flag.String("sPort", ":50051", "The port of the server")
)

func main() {
	flag.Parse()
	conn := ConnectToServer(*serverPort)
	defer conn.Close()
	c := pb.NewTemplateClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// get a stream to the server
	stream, err := c.SendChatMessage(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		var line string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line = scanner.Text()

			switch line {
			case "exit":
				stream.CloseSend()
				conn.Close()
				cancel()
				log.Fatal("Goodbye ", *userName)
			default:
				stream.Send(&pb.OutgoingChatMessage{UserName: *userName, Process: 1, Actions: 1, Message: line})
			}
		}
	}
}

func ConnectToServer(port string) *grpc.ClientConn {
	conn, connectionErr := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if connectionErr != nil {
		log.Fatalf("did not connect: %v", connectionErr)
	}
	return conn
}
