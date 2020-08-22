package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/satyamsoni2211/go_proto_example/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	address := flag.String("address", "[::]:8080", "port to connect to client")
	flag.Parse()
	log.Printf("Connecting client at %v", *address)
	con, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to server %v", err)
	}
	client := pb.NewFibonacciServiceClient(con)
	req := &pb.FibonacciRequest{
		Number: 10,
	}
	stream, err := client.CalculateFibonacci(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not receive response from %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream")
			return
		}
		if err != io.EOF && err != nil {
			log.Fatalf("Unable to receive stream response from server %v", err)
		}
		fmt.Println(res.GetResult())
	}
}
