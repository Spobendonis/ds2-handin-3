package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "proto/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userName      = flag.String("username", "anonymous", "The name others will see you by")
	serverAddress = flag.String("sAddress", ":50051", "The port of the server")
	actions       = 0
	process       = -1
)

func main() {
	flag.Parse()
	stream, conn, cancel, err := GetServerInfo(*serverAddress)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	defer cancel()
	var line string
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			if scanner.Scan() {
				line = scanner.Text()
				actions++
				switch line {
				case "join": //join chat
					fmt.Println("Choose an address to join:")
					if scanner.Scan() {
						a := scanner.Text()
						stream, conn, cancel, err = GetServerInfo(a)
						if err != nil {
							log.Println(err)
						}
					}
				case "leave": //leave chat
					actions = 0
					stream.CloseSend()
					conn.Close()
					cancel()
					log.Print("Leaving chat")
				case "exit": //exit application
					stream.CloseSend()
					conn.Close()
					cancel()
					log.Fatal("Goodbye ", *userName)
				default:
					stream.Send(&pb.OutgoingMessage{UserName: *userName, Process: int64(process), Actions: int64(actions), Message: line})
				}
			}
		}
	}()

	for {
		msg, e := stream.Recv()
		if e == nil {
			if actions < int(msg.Actions) {
				actions = int(msg.Actions)
			}
			log.Printf("%s (%d, %d)%s", msg.UserName, msg.Process, msg.Actions, msg.Message)
		}
	}
}

func GetServerInfo(addr string) (pb.Template_SendChatMessageClient, *grpc.ClientConn, context.CancelFunc, error) {
	conn := ConnectToServer(addr)
	c := pb.NewTemplateClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	// get a stream to the server
	stream, err := c.SendChatMessage(ctx)
	if err != nil {
		log.Println(err)
	}
	stream.Send(&pb.OutgoingMessage{UserName: *userName, Process: 0, Actions: 0, Message: ""})
	msg, e := stream.Recv()
	if e != nil {
		log.Println(e)
	}

	process = int(msg.Process)

	return stream, conn, cancel, err
}

func ConnectToServer(port string) *grpc.ClientConn {
	conn, connectionErr := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if connectionErr != nil {
		log.Fatalf("did not connect: %v", connectionErr)
	}
	actions++
	return conn
}
