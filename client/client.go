package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "proto/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverPort = flag.String("sPort", ":50051", "The port of the server")
	clientPort = flag.String("cPort", ":10000", "The port of the client")
)

func main() {
	flag.Parse()
	for {
		reader := bufio.NewReader(os.Stdin)
		_, readerErr := reader.ReadString('\n')
		if readerErr != nil {
			log.Fatal(readerErr)
		}
		conn, connectionErr := grpc.Dial(*serverPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if connectionErr != nil {
			log.Fatalf("did not connect: %v", connectionErr)
		}
		defer conn.Close()
		c := pb.NewTemplateClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, requestErr := c.GetWeather(ctx, &pb.WeatherRequest{Clientport: *clientPort, Location: ""})
		fmt.Println(r.GetWeather())
		if requestErr != nil {
			log.Fatal(requestErr)
		}
	}
}
