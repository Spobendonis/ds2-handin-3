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
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	pb.UnimplementedTemplateServer
}

func (s *Server) SendChatMessage(msgStream pb.Template_SendChatMessageServer) error {
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
