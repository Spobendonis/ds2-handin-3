package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "proto/proto"

	"google.golang.org/grpc"
)

var (
	port      = flag.Int("port", 50051, "The server port")
	processes = 0
)

type messageStruct struct {
	message  string
	username string
	process  int64
	actions  int64
}

var clientChannels []chan messageStruct

func handleMessageReceived(message string, username string, process int64, actions int64) {

	log.Print("Broadcasting")

	for _, channel := range clientChannels {
		channel <- messageStruct{message, username, process, actions}
	}

}

type Server struct {
	pb.UnimplementedTemplateServer
}

func (s Server) InitialiseConnection(ctx context.Context, in *pb.Dummy) (*pb.ProcessMessage, error) {
	processes++
	return (&pb.ProcessMessage{Process: int64(processes)}), nil
}

func (s *Server) SendChatMessage(msgStream pb.Template_SendChatMessageServer) error {

	clientChannel := make(chan messageStruct, 2)
	clientChannels = append(clientChannels, clientChannel)

	log.Print("Connected")
	handleMessageReceived("New user has connected", "", int64(processes), 0)

	go func() {
		for {
			toSend := <-clientChannel
			log.Print("Sending to client")
			toSendGRPC := &pb.IncomingChatMessage{UserName: toSend.username, Message: toSend.message, Process: toSend.process, Actions: toSend.actions}
			msgStream.Send(toSendGRPC)
		}
	}()

	for {
		// get the next message from the stream
		msg, err := msgStream.Recv()

		// the stream is closed so we can exit the loop
		if err == io.EOF {
			break
		}
		// some other error
		if err != nil {
			return err
		}
		// log the message
		log.Printf("%s (%d, %d): %s", msg.UserName, msg.Process, msg.Actions, msg.Message)
		handleMessageReceived(msg.Message, msg.UserName, msg.Process, msg.Actions)
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTemplateServer(grpcServer, &Server{})

	grpcServer.Serve(lis)
}
