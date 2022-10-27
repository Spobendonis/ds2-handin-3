package main

import (
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

func (s *Server) SendChatMessage(msgStream pb.Template_SendChatMessageServer) error {

	clientChannel := make(chan messageStruct, 2)
	clientChannels = append(clientChannels, clientChannel)

	initialMessage, initialError := msgStream.Recv()

	if initialError != nil {
		return initialError
	}

	processes++

	username := initialMessage.UserName
	process := processes
	msgStream.Send(&pb.IncomingMessage{UserName: username, Message: "", Process: int64(process), Actions: initialMessage.Actions})

	log.Print("Connected")
	handleMessageReceived(" has connected", username, int64(process), initialMessage.Actions)

	go func() {
		for {
			toSend := <-clientChannel
			log.Print("Sending to client")
			toSendGRPC := &pb.IncomingMessage{UserName: toSend.username, Message: toSend.message, Process: toSend.process, Actions: toSend.actions}
			msgStream.Send(toSendGRPC)
		}
	}()

	for {
		// get the next message from the stream
		msg, err := msgStream.Recv()

		// the stream is closed so we can exit the loop
		if err == io.EOF {
			handleMessageReceived(" has left", username, int64(process), initialMessage.Actions)
			break
		}
		// some other error
		if err != nil {
			handleMessageReceived(" has disconnected", username, int64(process), initialMessage.Actions)
			return err
		}
		// log the message
		log.Printf("%s (%d, %d): %s", msg.UserName, msg.Process, msg.Actions, msg.Message)
		handleMessageReceived(": "+msg.Message, msg.UserName, msg.Process, msg.Actions)
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
