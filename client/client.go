package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"

	pb "proto/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userName   = flag.String("username", "anonymous", "The name others will see you by")
	serverPort = flag.String("sPort", ":50051", "The port of the server")
	actions    = 0
	process    = -1
)

func main() {
	flag.Parse()
	conn := ConnectToServer(*serverPort)
	defer conn.Close()
	c := pb.NewTemplateClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	initReply, _ := c.InitialiseConnection(ctx, &pb.Dummy{})
	process = int(initReply.Process)
	defer cancel()
	// get a stream to the server
	stream, err := c.SendChatMessage(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	var line string
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			if scanner.Scan() {
				line = scanner.Text()
				actions++
				switch line {
				case "exit":
					stream.CloseSend()
					conn.Close()
					cancel()
					log.Fatal("Goodbye ", *userName)
				default:
					stream.Send(&pb.OutgoingChatMessage{UserName: *userName, Process: int64(process), Actions: int64(actions), Message: line})
				}
			}
		}
	}()

	for {
		msg, _ := stream.Recv()
		log.Printf("%s (%d, %d): %s", msg.UserName, msg.Process, msg.Actions, msg.Message)
	}

}

func ConnectToServer(port string) *grpc.ClientConn {
	conn, connectionErr := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if connectionErr != nil {
		log.Fatalf("did not connect: %v", connectionErr)
	}
	actions++
	return conn
}
