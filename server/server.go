package main

import (
	"context"
	"flag"
	"fmt"
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

func (s *Server) SendChatMessage(ctx context.Context) error {
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	server := &Server{}

	pb.RegisterTemplateServer(grpcServer, server)

	grpcServer.Serve(lis)
}
