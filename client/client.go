package main

import (
	"context"
	"flag"
	"log"
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
	c := pb.NewTemplateClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// get a stream to the server
	stream, err := c.SendChatMessage(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	// send some messages to the server
	stream.Send(&pb.OutgoingChatMessage{UserName: *userName, Process: 1, Actions: 1, Message: "hello world"})
	stream.Send(&pb.OutgoingChatMessage{UserName: *userName, Process: 1, Actions: 2, Message: "the sequel"})
	stream.Send(&pb.OutgoingChatMessage{UserName: *userName, Process: 1, Actions: 3, Message: "the seqseqsequel"})

	// close the stream
	stream.CloseSend()
	conn.Close()
	cancel()
}

func ConnectToServer(port string) *grpc.ClientConn {
	conn, connectionErr := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if connectionErr != nil {
		log.Fatalf("did not connect: %v", connectionErr)
	}
	return conn
}
